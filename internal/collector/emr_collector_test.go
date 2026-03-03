package collector

import "testing"

func TestCollector_emr_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewEMRCollector(), "emr")
}

func TestCollector_emr_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewEMRCollector(), true)
}
