package collector

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

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
