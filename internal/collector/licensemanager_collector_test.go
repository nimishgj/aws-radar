package collector

import (
	"testing"
)

func TestCollector_licensemanager_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewLicenseManagerCollector(), "licensemanager")
}

func TestCollector_licensemanager_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewLicenseManagerCollector(), true)
}
