package collector

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type fakeCollector struct {
	name         string
	called       bool
	region       string
	account      string
	accountName  string
	returnsError bool
}

func (f *fakeCollector) Name() string { return f.name }

func (f *fakeCollector) Collect(_ context.Context, _ aws.Config, region, account, accountName string) error {
	f.called = true
	f.region = region
	f.account = account
	f.accountName = accountName
	if f.returnsError {
		return errors.New("boom")
	}
	return nil
}

type fakeGlobalCollector struct {
	name         string
	called       bool
	account      string
	accountName  string
	returnsError bool
}

func (f *fakeGlobalCollector) Name() string { return f.name }

func (f *fakeGlobalCollector) Collect(_ context.Context, _ aws.Config, account, accountName string) error {
	f.called = true
	f.account = account
	f.accountName = accountName
	if f.returnsError {
		return errors.New("boom")
	}
	return nil
}

func TestNormalizeCollectors(t *testing.T) {
	got := normalizeCollectors([]string{" EC2 ", "s3", "S3", "", "  "})
	if len(got) != 2 {
		t.Fatalf("expected 2 collectors, got %d", len(got))
	}
	if _, ok := got["ec2"]; !ok {
		t.Fatalf("expected ec2 to be normalized")
	}
	if _, ok := got["s3"]; !ok {
		t.Fatalf("expected s3 to be normalized")
	}
}

func TestFilterCollectors(t *testing.T) {
	all := []Collector{
		&fakeCollector{name: "ec2"},
		&fakeCollector{name: "s3"},
	}
	enabled := map[string]struct{}{"s3": {}}
	filtered := filterCollectors(all, enabled)
	if len(filtered) != 1 || filtered[0].Name() != "s3" {
		t.Fatalf("expected only s3 collector, got %v", filtered)
	}
}

func TestRunCollectorRecordsSuccess(t *testing.T) {
	metrics.ResetAll()

	o := &Orchestrator{timeout: 2 * time.Second}
	c := &fakeCollector{name: "fake"}
	account := "acct"
	accountName := "my-account"
	region := "us-east-1"

	beforeErrors := testutil.ToFloat64(metrics.CollectionErrors.WithLabelValues(account, accountName, c.Name(), region))
	o.runCollector(context.Background(), c, aws.Config{}, region, account, accountName)

	if !c.called || c.account != account || c.accountName != accountName || c.region != region {
		t.Fatalf("collector not invoked with expected values")
	}
	if got := testutil.ToFloat64(metrics.CollectionUp.WithLabelValues(account, accountName, c.Name(), region)); got != 1 {
		t.Fatalf("expected up=1, got %v", got)
	}
	afterErrors := testutil.ToFloat64(metrics.CollectionErrors.WithLabelValues(account, accountName, c.Name(), region))
	if afterErrors != beforeErrors {
		t.Fatalf("expected no errors increment on success")
	}
}

func TestRunCollectorRecordsFailure(t *testing.T) {
	metrics.ResetAll()

	o := &Orchestrator{timeout: 2 * time.Second}
	c := &fakeCollector{name: "fake", returnsError: true}
	account := "acct"
	accountName := "my-account"
	region := "us-east-1"

	beforeErrors := testutil.ToFloat64(metrics.CollectionErrors.WithLabelValues(account, accountName, c.Name(), region))
	o.runCollector(context.Background(), c, aws.Config{}, region, account, accountName)

	if got := testutil.ToFloat64(metrics.CollectionUp.WithLabelValues(account, accountName, c.Name(), region)); got != 0 {
		t.Fatalf("expected up=0, got %v", got)
	}
	afterErrors := testutil.ToFloat64(metrics.CollectionErrors.WithLabelValues(account, accountName, c.Name(), region))
	if afterErrors != beforeErrors+1 {
		t.Fatalf("expected errors to increment by 1")
	}
}

func TestRunGlobalCollectorRecordsSuccess(t *testing.T) {
	metrics.ResetAll()

	o := &Orchestrator{timeout: 2 * time.Second}
	c := &fakeGlobalCollector{name: "global"}
	account := "acct"
	accountName := "my-account"

	beforeErrors := testutil.ToFloat64(metrics.CollectionErrors.WithLabelValues(account, accountName, c.Name(), "global"))
	o.runGlobalCollector(context.Background(), c, aws.Config{}, account, accountName)

	if !c.called || c.account != account || c.accountName != accountName {
		t.Fatalf("global collector not invoked with expected values")
	}
	if got := testutil.ToFloat64(metrics.CollectionUp.WithLabelValues(account, accountName, c.Name(), "global")); got != 1 {
		t.Fatalf("expected up=1, got %v", got)
	}
	afterErrors := testutil.ToFloat64(metrics.CollectionErrors.WithLabelValues(account, accountName, c.Name(), "global"))
	if afterErrors != beforeErrors {
		t.Fatalf("expected no errors increment on success")
	}
}
