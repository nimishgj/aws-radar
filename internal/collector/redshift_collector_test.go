package collector

import "testing"

func TestCollector_redshift_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewRedshiftCollector(), "redshift")
}

func TestCollector_redshift_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewRedshiftCollector(), true)
}
