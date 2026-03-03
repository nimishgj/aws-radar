package collector

import "testing"

func TestCollector_efs_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewEFSCollector(), "efs")
}

func TestCollector_efs_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewEFSCollector(), true)
}
