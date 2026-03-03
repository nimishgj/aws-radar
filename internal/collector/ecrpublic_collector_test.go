package collector

import "testing"

func TestGlobalCollector_ecrpublic_Name(t *testing.T) {
	assertGlobalCollectorName(t, NewECRPublicCollector(), "ecrpublic")
}

func TestGlobalCollector_ecrpublic_ErrorContract(t *testing.T) {
	assertGlobalCollectorErrorContract(t, NewECRPublicCollector(), true)
}
