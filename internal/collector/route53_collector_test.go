package collector

import "testing"

func TestGlobalCollector_route53_Name(t *testing.T) {
	assertGlobalCollectorName(t, NewRoute53Collector(), "route53")
}

func TestGlobalCollector_route53_ErrorContract(t *testing.T) {
	assertGlobalCollectorErrorContract(t, NewRoute53Collector(), true)
}
