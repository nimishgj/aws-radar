package collector

import "testing"

func TestCollector_ec2_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewEC2Collector(), "ec2")
}

func TestCollector_ec2_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewEC2Collector(), true)
}
