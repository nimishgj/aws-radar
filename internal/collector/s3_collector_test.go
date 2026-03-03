package collector

import "testing"

func TestCollector_s3_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewS3Collector(), "s3")
}

func TestCollector_s3_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewS3Collector(), true)
}
