package collector

import "testing"

func TestCollector_apigatewayv2_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewAPIGatewayV2Collector(), "apigatewayv2")
}

func TestCollector_apigatewayv2_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewAPIGatewayV2Collector(), true)
}
