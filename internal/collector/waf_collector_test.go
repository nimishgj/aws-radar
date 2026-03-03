package collector

import "testing"

func TestCollector_waf_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewWAFCollector(), "waf")
}

func TestCollector_waf_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewWAFCollector(), true)
}
