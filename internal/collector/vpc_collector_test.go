package collector

import "testing"

func TestCollector_vpc_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewVPCCollector(), "vpc")
}

func TestCollector_vpc_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewVPCCollector(), false)
}
