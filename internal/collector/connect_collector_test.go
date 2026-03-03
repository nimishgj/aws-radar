package collector

import "testing"

func TestCollector_connect_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewConnectCollector(), "connect")
}

func TestCollector_connect_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewConnectCollector(), true)
}
