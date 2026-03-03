package collector

import "testing"

func TestCollector_s3control_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewS3ControlCollector(), "s3control")
}

func TestCollector_s3control_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewS3ControlCollector(), true)
}
