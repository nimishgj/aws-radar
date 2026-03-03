package collector

import "testing"

func TestCollector_memorydb_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewMemoryDBCollector(), "memorydb")
}

func TestCollector_memorydb_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewMemoryDBCollector(), true)
}
