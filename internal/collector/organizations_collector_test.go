package collector

import "testing"

func TestGlobalCollector_organizations_Name(t *testing.T) {
	assertGlobalCollectorName(t, NewOrganizationsCollector(), "organizations")
}

func TestGlobalCollector_organizations_ErrorContract(t *testing.T) {
	assertGlobalCollectorErrorContract(t, NewOrganizationsCollector(), true)
}
