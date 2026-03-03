package collector

import "testing"

func TestCollector_apigateway_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewAPIGatewayCollector(), "apigateway")
}

func TestCollector_apigateway_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewAPIGatewayCollector(), true)
}
