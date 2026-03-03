package collector

import "testing"

func TestGlobalCollector_cloudfront_Name(t *testing.T) {
	assertGlobalCollectorName(t, NewCloudFrontCollector(), "cloudfront")
}

func TestGlobalCollector_cloudfront_ErrorContract(t *testing.T) {
	assertGlobalCollectorErrorContract(t, NewCloudFrontCollector(), true)
}
