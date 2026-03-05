package collector

import (
	"math"
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

	c.cfg.Format = "Parquet"
	if f := c.detectFormat("report.csv.gz"); f != "parquet" {
		t.Fatalf("expected parquet override, got %s", f)
	}
}

func floatClose(a, b float64) bool {
	return math.Abs(a-b) < 0.001
}
