package collector

import "testing"

func TestCollector_rds_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewRDSCollector(), "rds")
}

func TestCollector_rds_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewRDSCollector(), true)
}
