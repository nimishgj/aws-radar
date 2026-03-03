package collector

import "testing"

func TestCollector_kinesisvideo_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewKinesisVideoCollector(), "kinesisvideo")
}

func TestCollector_kinesisvideo_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewKinesisVideoCollector(), true)
}
