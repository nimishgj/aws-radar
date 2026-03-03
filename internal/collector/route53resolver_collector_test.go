package collector

import "testing"

func TestCollector_route53resolver_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewRoute53ResolverCollector(), "route53resolver")
}

func TestCollector_route53resolver_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewRoute53ResolverCollector(), true)
}
