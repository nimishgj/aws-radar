package collector

import "testing"

func TestCollector_msk_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewMSKCollector(), "msk")
}

func TestCollector_msk_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewMSKCollector(), true)
}
