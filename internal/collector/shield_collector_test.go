package collector

import "testing"

func TestGlobalCollector_shield_Name(t *testing.T) {
	assertGlobalCollectorName(t, NewShieldCollector(), "shield")
}

func TestGlobalCollector_shield_ErrorContract(t *testing.T) {
	assertGlobalCollectorErrorContract(t, NewShieldCollector(), true)
}
