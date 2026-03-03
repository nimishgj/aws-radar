package collector

import "testing"

func TestCollector_opensearchserverless_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewOpenSearchServerlessCollector(), "opensearchserverless")
}

func TestCollector_opensearchserverless_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewOpenSearchServerlessCollector(), true)
}
