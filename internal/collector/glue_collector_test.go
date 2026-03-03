package collector

import "testing"

func TestCollector_glue_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewGlueCollector(), "glue")
}

func TestCollector_glue_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewGlueCollector(), true)
}
