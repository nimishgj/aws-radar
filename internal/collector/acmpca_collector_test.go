package collector

import (
	"testing"
)

func TestCollector_acmpca_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewACMPCACollector(), "acmpca")
}

func TestCollector_acmpca_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewACMPCACollector(), true)
}
