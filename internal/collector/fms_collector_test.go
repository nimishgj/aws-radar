package collector

import (
	"testing"
)

func TestCollector_fms_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewFMSCollector(), "fms")
}

func TestCollector_fms_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewFMSCollector(), true)
}
