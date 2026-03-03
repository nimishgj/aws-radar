package collector

import "testing"

func TestCollector_documentdb_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewDocumentDBCollector(), "documentdb")
}

func TestCollector_documentdb_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewDocumentDBCollector(), true)
}
