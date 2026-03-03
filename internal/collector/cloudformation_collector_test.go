package collector

import "testing"

func TestCollector_cloudformation_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewCloudFormationCollector(), "cloudformation")
}

func TestCollector_cloudformation_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewCloudFormationCollector(), true)
}
