package collector

import (
	"testing"
)

func TestCollector_servicecatalog_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewServiceCatalogCollector(), "servicecatalog")
}

func TestCollector_servicecatalog_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewServiceCatalogCollector(), true)
}
