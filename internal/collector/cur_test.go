package collector

import (
	"archive/zip"
	"bytes"
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"testing"
	"time"

	"github.com/nimishgj/aws-radar/internal/config"
)

func TestGlobalCollector_cost_cur_Name(t *testing.T) {
	c := NewCURCollector(config.CostCURConfig{Bucket: "b", ReportName: "r"})
	assertGlobalCollectorName(t, c, "cost_cur")
}

func TestGlobalCollector_cost_cur_ErrorContract(t *testing.T) {
	c := NewCURCollector(config.CostCURConfig{Bucket: "b", ReportName: "r"})
	assertGlobalCollectorErrorContract(t, c, true)
}

func TestCURCollectorDefaults(t *testing.T) {
	c := NewCURCollector(config.CostCURConfig{})
	if c.cfg.Frequency != "daily" {
		t.Fatalf("expected default frequency 'daily', got %q", c.cfg.Frequency)
	}
	if c.cfg.Region != "us-east-1" {
		t.Fatalf("expected default region 'us-east-1', got %q", c.cfg.Region)
	}
	if c.cfg.MaxResources != 100 {
		t.Fatalf("expected default max_resources 100, got %d", c.cfg.MaxResources)
	}
}

func TestCURCollectorInvalidFrequency(t *testing.T) {
	c := NewCURCollector(config.CostCURConfig{Frequency: "weird"})
	if c.cfg.Frequency != "daily" {
		t.Fatalf("expected invalid frequency to fall back to 'daily', got %q", c.cfg.Frequency)
	}
}

func TestCURCollectorShouldRefresh(t *testing.T) {
	c := NewCURCollector(config.CostCURConfig{Frequency: "daily"})
	now := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)

	if !c.shouldRefresh(now) {
		t.Fatalf("expected refresh when no previous fetch")
	}

	c.lastFetch = now.Add(-2 * time.Hour)
	if c.shouldRefresh(now) {
		t.Fatalf("did not expect refresh on same UTC date")
	}

	c.lastFetch = now.Add(-26 * time.Hour)
	if !c.shouldRefresh(now) {
		t.Fatalf("expected refresh on next UTC date")
	}
}

func TestCURCollectorShouldRefreshHourly(t *testing.T) {
	c := NewCURCollector(config.CostCURConfig{Frequency: "hourly"})
	now := time.Date(2026, 3, 1, 12, 30, 0, 0, time.UTC)

	c.lastFetch = now.Add(-20 * time.Minute)
	if c.shouldRefresh(now) {
		t.Fatalf("expected no refresh within same hour")
	}

	c.lastFetch = now.Add(-90 * time.Minute)
	if !c.shouldRefresh(now) {
		t.Fatalf("expected refresh on next hour")
	}
}

func TestParseManifest(t *testing.T) {
	data := []byte(`{
		"reportKeys": [
			"prefix/report/20260301-20260401/report-00001.snappy.parquet",
			"prefix/report/20260301-20260401/report-00002.snappy.parquet"
		]
	}`)

	m, err := ParseManifest(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.ReportKeys) != 2 {
		t.Fatalf("expected 2 report keys, got %d", len(m.ReportKeys))
	}
	if m.ReportKeys[0] != "prefix/report/20260301-20260401/report-00001.snappy.parquet" {
		t.Fatalf("unexpected report key: %s", m.ReportKeys[0])
	}
}

func TestParseManifestCUR2Variants(t *testing.T) {
	cases := []struct {
		name string
		json string
		want []string
	}{
		{"dataFileS3Keys strings", `{"dataFileS3Keys":["pfx/exp/data/BILLING_PERIOD=2026-08/exp-00001.snappy.parquet"]}`,
			[]string{"pfx/exp/data/BILLING_PERIOD=2026-08/exp-00001.snappy.parquet"}},
		{"dataFiles objects with key", `{"dataFiles":[{"key":"a.parquet"},{"key":"b.parquet"}]}`,
			[]string{"a.parquet", "b.parquet"}},
		{"files objects with s3Key", `{"files":[{"s3Key":"c.csv.gz"}]}`,
			[]string{"c.csv.gz"}},
		{"s3 uri normalized", `{"dataFileS3Keys":["s3://my-bucket/pfx/exp-00001.snappy.parquet"]}`,
			[]string{"pfx/exp-00001.snappy.parquet"}},
		{"legacy reportKeys", `{"reportKeys":["r.csv.gz"]}`,
			[]string{"r.csv.gz"}},
		{"dataFiles as s3 uri strings", `{"executionId":"e-1","dataFiles":["s3://bkt/pfx/exp/data/BILLING_PERIOD=2026-06/part.parquet"]}`,
			[]string{"pfx/exp/data/BILLING_PERIOD=2026-06/part.parquet"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := ParseManifest([]byte(tc.json))
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			got := m.keys()
			if len(got) != len(tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("expected %v, got %v", tc.want, got)
				}
			}
		})
	}
}

func TestManifestKeyCandidates(t *testing.T) {
	now := time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC)

	legacy := NewCURCollector(config.CostCURConfig{
		Prefix: "cur-reports", ReportName: "my-report",
	})
	got := legacy.manifestKeyCandidates(now)
	want := "cur-reports/my-report/20260801-20260901/my-report-Manifest.json"
	if len(got) != 1 || got[0] != want {
		t.Fatalf("expected [%s], got %v", want, got)
	}

	cur2 := NewCURCollector(config.CostCURConfig{
		ReportType: "cur2", Prefix: "exports", ReportName: "my-export",
	})
	got = cur2.manifestKeyCandidates(now)
	want0 := "exports/my-export/metadata/BILLING_PERIOD=2026-08/my-export-Manifest.json"
	want1 := "exports/my-export/metadata/BILLING_PERIOD=2026-08/Manifest.json"
	if len(got) != 2 || got[0] != want0 || got[1] != want1 {
		t.Fatalf("expected [%s %s], got %v", want0, want1, got)
	}
}

func TestParseManifestEmpty(t *testing.T) {
	data := []byte(`{"reportKeys": []}`)
	m, err := ParseManifest(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.ReportKeys) != 0 {
		t.Fatalf("expected 0 report keys, got %d", len(m.ReportKeys))
	}
}

func TestParseManifestInvalid(t *testing.T) {
	_, err := ParseManifest([]byte(`not json`))
	if err == nil {
		t.Fatalf("expected error for invalid JSON")
	}
}

func TestCSVColumnMap(t *testing.T) {
	header := []string{
		"lineItem/UsageAccountId",
		"lineItem/ProductCode",
		"lineItem/UnblendedCost",
		"lineItem/ResourceId",
		"lineItem/UsageType",
		"lineItem/LineItemType",
		"resourceTags/user:Environment",
		"resourceTags/user:Team",
		"someOtherColumn",
	}

	m := CSVColumnMap(header)
	if m.accountID != 0 {
		t.Fatalf("expected accountID at index 0, got %d", m.accountID)
	}
	if m.service != 1 {
		t.Fatalf("expected service at index 1, got %d", m.service)
	}
	if m.cost != 2 {
		t.Fatalf("expected cost at index 2, got %d", m.cost)
	}
	if m.resource != 3 {
		t.Fatalf("expected resource at index 3, got %d", m.resource)
	}
	if m.usageType != 4 {
		t.Fatalf("expected usageType at index 4, got %d", m.usageType)
	}
	if m.lineItemType != 5 {
		t.Fatalf("expected lineItemType at index 5, got %d", m.lineItemType)
	}
	if len(m.tagCols) != 2 {
		t.Fatalf("expected 2 tag columns, got %d", len(m.tagCols))
	}
	if m.tagCols[6] != "Environment" {
		t.Fatalf("expected tag 'Environment' at index 6, got %q", m.tagCols[6])
	}
	if m.tagCols[7] != "Team" {
		t.Fatalf("expected tag 'Team' at index 7, got %q", m.tagCols[7])
	}
}

func TestCSVColumnMapSnakeCase(t *testing.T) {
	header := []string{
		"line_item_usage_account_id",
		"line_item_product_code",
		"line_item_unblended_cost",
		"line_item_resource_id",
		"line_item_usage_type",
		"line_item_line_item_type",
		"resource_tags_user_Environment",
		"some_other_column",
	}

	m := CSVColumnMap(header)
	if m.accountID != 0 {
		t.Fatalf("expected accountID at index 0, got %d", m.accountID)
	}
	if m.service != 1 {
		t.Fatalf("expected service at index 1, got %d", m.service)
	}
	if m.cost != 2 {
		t.Fatalf("expected cost at index 2, got %d", m.cost)
	}
	if m.resource != 3 {
		t.Fatalf("expected resource at index 3, got %d", m.resource)
	}
	if m.usageType != 4 {
		t.Fatalf("expected usageType at index 4, got %d", m.usageType)
	}
	if m.lineItemType != 5 {
		t.Fatalf("expected lineItemType at index 5, got %d", m.lineItemType)
	}
	if m.tagCols[6] != "Environment" {
		t.Fatalf("expected tag 'Environment' at index 6, got %q", m.tagCols[6])
	}
}

func TestAggregateCSVRow(t *testing.T) {
	header := []string{
		"lineItem/ProductCode",
		"lineItem/UnblendedCost",
		"lineItem/ResourceId",
		"lineItem/UsageType",
		"lineItem/LineItemType",
		"resourceTags/user:Environment",
	}
	colMap := CSVColumnMap(header)
	agg := NewCURAggregation("20260301")

	// Row 1: EC2 instance (Usage)
	AggregateCSVRow([]string{"AmazonEC2", "10.50", "i-abc123", "BoxUsage:t3.micro", "Usage", "prod"}, colMap, agg)
	// Row 2: Same service, different resource (Usage)
	AggregateCSVRow([]string{"AmazonEC2", "5.25", "i-def456", "BoxUsage:t3.small", "Usage", "dev"}, colMap, agg)
	// Row 3: Different service (Tax)
	AggregateCSVRow([]string{"AmazonS3", "1.00", "my-bucket", "TimedStorage-ByteHrs", "Tax", "prod"}, colMap, agg)
	// Row 4: Zero cost (should be skipped)
	AggregateCSVRow([]string{"AmazonEC2", "0", "i-skip", "BoxUsage:t3.micro", "Usage", "prod"}, colMap, agg)
	// Row 5: Credit (should be skipped)
	AggregateCSVRow([]string{"AmazonEC2", "-3.00", "i-abc123", "BoxUsage:t3.micro", "Credit", "prod"}, colMap, agg)
	// Row 6: SavingsPlanNegation (should be skipped)
	AggregateCSVRow([]string{"AmazonEC2", "-10.50", "i-abc123", "BoxUsage:t3.micro", "SavingsPlanNegation", "prod"}, colMap, agg)

	// Check totals.
	if !floatClose(agg.total, 16.75) {
		t.Fatalf("expected total 16.75, got %f", agg.total)
	}

	// Check by service.
	if !floatClose(agg.byService["AmazonEC2"], 15.75) {
		t.Fatalf("expected EC2 cost 15.75, got %f", agg.byService["AmazonEC2"])
	}
	if !floatClose(agg.byService["AmazonS3"], 1.00) {
		t.Fatalf("expected S3 cost 1.00, got %f", agg.byService["AmazonS3"])
	}

	// Check by resource.
	if len(agg.byResource) != 3 {
		t.Fatalf("expected 3 resources, got %d", len(agg.byResource))
	}
	rc := agg.byResource["AmazonEC2|i-abc123"]
	if !floatClose(rc.cost, 10.50) {
		t.Fatalf("expected resource cost 10.50, got %f", rc.cost)
	}

	// Check by usage type.
	if !floatClose(agg.byUsageType["AmazonEC2|BoxUsage:t3.micro"], 10.50) {
		t.Fatalf("expected usage type cost 10.50, got %f", agg.byUsageType["AmazonEC2|BoxUsage:t3.micro"])
	}

	// Check by tag.
	if !floatClose(agg.byTag["Environment|prod"], 11.50) {
		t.Fatalf("expected tag prod cost 11.50, got %f", agg.byTag["Environment|prod"])
	}
	if !floatClose(agg.byTag["Environment|dev"], 5.25) {
		t.Fatalf("expected tag dev cost 5.25, got %f", agg.byTag["Environment|dev"])
	}
}

func TestTopNResources(t *testing.T) {
	resources := map[string]resourceCost{
		"svc|r1": {service: "svc", resourceID: "r1", cost: 10},
		"svc|r2": {service: "svc", resourceID: "r2", cost: 50},
		"svc|r3": {service: "svc", resourceID: "r3", cost: 30},
		"svc|r4": {service: "svc", resourceID: "r4", cost: 5},
		"svc|r5": {service: "svc", resourceID: "r5", cost: 20},
	}

	top3 := topNResources(resources, 3)
	if len(top3) != 3 {
		t.Fatalf("expected 3 resources, got %d", len(top3))
	}
	if top3[0].resourceID != "r2" {
		t.Fatalf("expected top resource r2, got %s", top3[0].resourceID)
	}
	if top3[1].resourceID != "r3" {
		t.Fatalf("expected second resource r3, got %s", top3[1].resourceID)
	}

	// N larger than total.
	all := topNResources(resources, 100)
	if len(all) != 5 {
		t.Fatalf("expected 5 resources, got %d", len(all))
	}
}

func TestBillingPeriod(t *testing.T) {
	now := time.Date(2026, 3, 15, 10, 0, 0, 0, time.UTC)
	start, end := billingPeriod(now)
	if start != "20260301" {
		t.Fatalf("expected start 20260301, got %s", start)
	}
	if end != "20260401" {
		t.Fatalf("expected end 20260401, got %s", end)
	}
}

func TestDetectFormat(t *testing.T) {
	c := NewCURCollector(config.CostCURConfig{})

	if f := c.detectFormat("report.snappy.parquet"); f != "parquet" {
		t.Fatalf("expected parquet, got %s", f)
	}
	if f := c.detectFormat("report.csv.gz"); f != "csv" {
		t.Fatalf("expected csv, got %s", f)
	}
	if f := c.detectFormat("report.csv.zip"); f != "zip" {
		t.Fatalf("expected zip, got %s", f)
	}

	c.cfg.Format = "Parquet"
	if f := c.detectFormat("report.csv.gz"); f != "parquet" {
		t.Fatalf("expected parquet override, got %s", f)
	}
}

func TestParseZIPSample(t *testing.T) {
	f, err := os.Open("testdata/cur-sample.csv.zip")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer func() { _ = f.Close() }()

	stat, err := f.Stat()
	if err != nil {
		t.Fatalf("stat fixture: %v", err)
	}

	agg := &curAggregation{
		byService:   make(map[string]float64),
		byResource:  make(map[string]resourceCost),
		byUsageType: make(map[string]float64),
		byTag:       make(map[string]float64),
	}

	if err := parseZIP(f, stat.Size(), agg); err != nil {
		t.Fatalf("parseZIP: %v", err)
	}

	if !floatClose(agg.total, 0.0035) {
		t.Fatalf("expected total 0.0035, got %f", agg.total)
	}
	if !floatClose(agg.byService["AmazonEC2"], 0.0035) {
		t.Fatalf("expected AmazonEC2 cost 0.0035, got %f", agg.byService["AmazonEC2"])
	}
	if _, ok := agg.byResource["AmazonEC2|i-0f0bf49951ab73706"]; !ok {
		t.Fatalf("expected resource entry for AmazonEC2|i-0f0bf49951ab73706, got %v", agg.byResource)
	}
}

// newMockS3Server serves objects by exact URL path (path-style addressing),
// returning S3-style 404s for unknown keys, and records requested paths.
func newMockS3Server(objects map[string]string, requested *[]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*requested = append(*requested, r.URL.Path)
		if body, ok := objects[r.URL.Path]; ok {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(body))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`<?xml version="1.0"?><Error><Code>NoSuchKey</Code></Error>`))
	}))
}

func mockS3Config(serverURL string) aws.Config {
	resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, reg string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: serverURL, SigningRegion: "us-east-1", HostnameImmutable: true}, nil
		},
	)
	return aws.Config{
		Region:                      "us-east-1",
		Credentials:                 aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("test", "test", "")),
		EndpointResolverWithOptions: resolver,
		Retryer:                     func() aws.Retryer { return aws.NopRetryer{} },
	}
}

func TestRefreshManifestFallback(t *testing.T) {
	now := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	csvBody := "lineItem/ProductCode,lineItem/UnblendedCost,lineItem/LineItemType\nAmazonEC2,2.5,Usage\n"
	dataKey := "/test-bucket/exports/exp/data/BILLING_PERIOD=2026-09/exp-00001.csv"
	// Only the second manifest candidate (Manifest.json) exists.
	manifest := `{"dataFiles":["exports/exp/data/BILLING_PERIOD=2026-09/exp-00001.csv"]}`

	var requested []string
	server := newMockS3Server(map[string]string{
		"/test-bucket/exports/exp/metadata/BILLING_PERIOD=2026-09/Manifest.json": manifest,
		dataKey: csvBody,
	}, &requested)
	defer server.Close()

	c := NewCURCollector(config.CostCURConfig{
		Enabled: true, ReportType: "cur2", Bucket: "test-bucket",
		Prefix: "exports", ReportName: "exp", Region: "us-east-1",
	})

	if err := c.refresh(context.Background(), mockS3Config(server.URL), now); err != nil {
		t.Fatalf("refresh: %v", err)
	}

	firstCandidate := "/test-bucket/exports/exp/metadata/BILLING_PERIOD=2026-09/exp-Manifest.json"
	if len(requested) == 0 || requested[0] != firstCandidate {
		t.Fatalf("expected first request to be %s, got %v", firstCandidate, requested)
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.cached == nil {
		t.Fatal("expected cached aggregation after refresh")
	}
	if !floatClose(c.cached.total, 2.5) {
		t.Fatalf("expected total 2.5 from fallback manifest, got %f", c.cached.total)
	}
}

func TestRefreshAllManifestCandidatesFail(t *testing.T) {
	now := time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)
	var requested []string
	server := newMockS3Server(map[string]string{}, &requested)
	defer server.Close()

	c := NewCURCollector(config.CostCURConfig{
		Enabled: true, ReportType: "cur2", Bucket: "test-bucket",
		Prefix: "exports", ReportName: "exp", Region: "us-east-1",
	})

	err := c.refresh(context.Background(), mockS3Config(server.URL), now)
	if err == nil {
		t.Fatal("expected error when all manifest candidates fail")
	}
	if !strings.Contains(err.Error(), "exp-Manifest.json") || !strings.Contains(err.Error(), "tried") {
		t.Fatalf("expected error naming tried candidates, got: %v", err)
	}
	if len(requested) != 2 {
		t.Fatalf("expected both candidates requested, got %v", requested)
	}
}

func buildZIP(t *testing.T, entries map[string]string) *bytes.Reader {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range entries {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("create zip entry: %v", err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatalf("write zip entry: %v", err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	return bytes.NewReader(buf.Bytes())
}

func TestParseZIPEdgeCases(t *testing.T) {
	csvContent := "lineItem/ProductCode,lineItem/UnblendedCost\nAmazonEC2,1.5\n"

	t.Run("no csv entries", func(t *testing.T) {
		r := buildZIP(t, map[string]string{"readme.txt": "hello"})
		agg := NewCURAggregation("")
		if err := parseZIP(r, r.Size(), agg); err != nil {
			t.Fatalf("expected no error for zip without csv entries, got %v", err)
		}
		if agg.total != 0 {
			t.Fatalf("expected zero total, got %f", agg.total)
		}
	})

	t.Run("invalid csv propagates error", func(t *testing.T) {
		r := buildZIP(t, map[string]string{"bad.csv": "col1,col2\nonly_one_field\n"})
		agg := NewCURAggregation("")
		if err := parseZIP(r, r.Size(), agg); err == nil {
			t.Fatal("expected error for invalid csv content, got nil")
		}
	})

	t.Run("skips __MACOSX entries", func(t *testing.T) {
		r := buildZIP(t, map[string]string{
			"report.csv":            csvContent,
			"__MACOSX/._report.csv": "\x00\x01not a csv",
		})
		agg := NewCURAggregation("")
		if err := parseZIP(r, r.Size(), agg); err != nil {
			t.Fatalf("parseZIP: %v", err)
		}
		if !floatClose(agg.total, 1.5) {
			t.Fatalf("expected total 1.5, got %f", agg.total)
		}
	})

	t.Run("not a zip", func(t *testing.T) {
		r := bytes.NewReader([]byte("plain text, not a zip"))
		if err := parseZIP(r, r.Size(), NewCURAggregation("")); err == nil {
			t.Fatal("expected error for non-zip input, got nil")
		}
	})
}

func floatClose(a, b float64) bool {
	return math.Abs(a-b) < 0.001
}
