package collector

import "testing"

func TestCollector_opensearch_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewOpenSearchCollector(), "opensearch")
}

func TestCollector_opensearch_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewOpenSearchCollector(), true)
}
