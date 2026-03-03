package collector

import "testing"

func TestCollector_cloudtrail_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewCloudTrailCollector(), "cloudtrail")
}

func TestCollector_cloudtrail_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewCloudTrailCollector(), true)
}
