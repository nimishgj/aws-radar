package collector

import "testing"

func TestGlobalCollector_directconnect_Name(t *testing.T) {
	assertGlobalCollectorName(t, NewDirectConnectCollector(), "directconnect")
}

func TestGlobalCollector_directconnect_ErrorContract(t *testing.T) {
	assertGlobalCollectorErrorContract(t, NewDirectConnectCollector(), true)
}
