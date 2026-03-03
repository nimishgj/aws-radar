package collector

import "testing"

func TestCollector_appstream_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewAppStreamCollector(), "appstream")
}

func TestCollector_appstream_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewAppStreamCollector(), true)
}
