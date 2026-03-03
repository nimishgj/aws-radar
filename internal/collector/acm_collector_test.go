package collector

import "testing"

func TestCollector_acm_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewACMCollector(), "acm")
}

func TestCollector_acm_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewACMCollector(), true)
}
