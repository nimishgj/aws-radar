package collector

import "testing"

func TestCollector_dms_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewDMSCollector(), "dms")
}

func TestCollector_dms_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewDMSCollector(), true)
}
