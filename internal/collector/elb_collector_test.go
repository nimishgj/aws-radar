package collector

import "testing"

func TestCollector_elb_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewELBCollector(), "elb")
}

func TestCollector_elb_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewELBCollector(), false)
}
