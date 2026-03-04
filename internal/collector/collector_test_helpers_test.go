package collector

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// mockRoute maps a request matcher to a response body.
type mockRoute struct {
	// matcher returns true if this route handles the request.
	matcher func(r *http.Request) bool
	body    string
}

// newMockAWSServer creates an httptest.Server that routes requests to canned
// responses. Routes are evaluated in order; the first match wins.
func newMockAWSServer(routes []mockRoute) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read body for form-encoded services (Query protocol)
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

		for _, route := range routes {
			if route.matcher(r) {
				w.Header().Set("Content-Type", "text/xml")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(route.body))
				return
			}
		}
		// Return empty 200 for unmatched requests
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
}

// mockAWSConfig creates an aws.Config that points to the given test server.
func mockAWSConfig(serverURL, region string) aws.Config {
	resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, reg string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: serverURL, SigningRegion: region}, nil
		},
	)
	return aws.Config{
		Region:                      region,
		Credentials:                 aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("test", "test", "")),
		EndpointResolverWithOptions: resolver,
	}
}

// queryAction returns a matcher for AWS Query-protocol requests (ELB, RDS, ElastiCache).
func queryAction(action string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		bodyBytes, _ := io.ReadAll(r.Body)
		body := string(bodyBytes)
		r.Body = io.NopCloser(strings.NewReader(body))
		return strings.Contains(body, "Action="+action)
	}
}

// jsonTarget returns a matcher for AWS JSON-protocol requests (ECS).
func jsonTarget(target string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		return strings.Contains(r.Header.Get("X-Amz-Target"), target)
	}
}

// restPath returns a matcher for AWS REST-JSON requests (SES v2).
func restPath(pathPrefix string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		return strings.HasPrefix(r.URL.Path, pathPrefix)
	}
}

// gaugeValue is a test helper to read the current float64 of a gauge with labels.
func gaugeValue(g *prometheus.GaugeVec, labels ...string) float64 {
	return testutil.ToFloat64(g.WithLabelValues(labels...))
}

type alwaysFailHTTPClient struct{}

func (c *alwaysFailHTTPClient) Do(*http.Request) (*http.Response, error) {
	return nil, errors.New("forced transport failure")
}

func testAWSConfig(region string) aws.Config {
	return aws.Config{
		Region:      region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("test", "test", "")),
		HTTPClient:  &alwaysFailHTTPClient{},
	}
}

func assertRegionalCollectorName(t *testing.T, c Collector, expected string) {
	t.Helper()
	if got := c.Name(); got != expected {
		t.Fatalf("collector name mismatch: expected %q got %q", expected, got)
	}
}

func assertGlobalCollectorName(t *testing.T, c GlobalCollector, expected string) {
	t.Helper()
	if got := c.Name(); got != expected {
		t.Fatalf("global collector name mismatch: expected %q got %q", expected, got)
	}
}

func assertRegionalCollectorErrorContract(t *testing.T, c Collector, wantErr bool) {
	t.Helper()
	cfg := testAWSConfig("us-east-1")
	err := c.Collect(context.Background(), cfg, "us-east-1", "123456789012", "test-account")
	if wantErr && err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !wantErr && err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func assertGlobalCollectorErrorContract(t *testing.T, c GlobalCollector, wantErr bool) {
	t.Helper()
	cfg := testAWSConfig("us-east-1")
	err := c.Collect(context.Background(), cfg, "123456789012", "test-account")
	if wantErr && err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !wantErr && err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
