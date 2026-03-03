package collector

import "testing"

func TestCollector_macie_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewMacieCollector(), "macie")
}

func TestCollector_macie_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewMacieCollector(), true)
}
