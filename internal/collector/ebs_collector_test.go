package collector

import "testing"

func TestCollector_ebs_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewEBSCollector(), "ebs")
}

func TestCollector_ebs_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewEBSCollector(), true)
}
