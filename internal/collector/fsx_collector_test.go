package collector

import "testing"

func TestCollector_fsx_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewFSxCollector(), "fsx")
}

func TestCollector_fsx_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewFSxCollector(), true)
}
