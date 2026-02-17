package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestResetAllResetsGauges(t *testing.T) {
	EC2Instances.WithLabelValues("acct", "us-east-1", "t3.micro", "running", "us-east-1a").Set(5)
	S3Buckets.WithLabelValues("acct", "us-east-1").Set(3)
	CollectionUp.WithLabelValues("acct", "ec2", "us-east-1").Set(1)

	ResetAll()

	if got := testutil.ToFloat64(EC2Instances.WithLabelValues("acct", "us-east-1", "t3.micro", "running", "us-east-1a")); got != 0 {
		t.Fatalf("expected EC2Instances to reset to 0, got %v", got)
	}
	if got := testutil.ToFloat64(S3Buckets.WithLabelValues("acct", "us-east-1")); got != 0 {
		t.Fatalf("expected S3Buckets to reset to 0, got %v", got)
	}
	if got := testutil.ToFloat64(CollectionUp.WithLabelValues("acct", "ec2", "us-east-1")); got != 0 {
		t.Fatalf("expected CollectionUp to reset to 0, got %v", got)
	}
}
