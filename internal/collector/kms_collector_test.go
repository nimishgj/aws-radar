package collector

import "testing"

func TestCollector_kms_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewKMSCollector(), "kms")
}

func TestCollector_kms_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewKMSCollector(), true)
}
