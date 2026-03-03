package collector

import "testing"

func TestCollector_workspaces_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewWorkSpacesCollector(), "workspaces")
}

func TestCollector_workspaces_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewWorkSpacesCollector(), true)
}
