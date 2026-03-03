package collector

import "testing"

func TestCollector_kinesis_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewKinesisCollector(), "kinesis")
}

func TestCollector_kinesis_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewKinesisCollector(), true)
}
