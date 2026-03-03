package collector

import "testing"

func TestCollector_lambda_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewLambdaCollector(), "lambda")
}

func TestCollector_lambda_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewLambdaCollector(), true)
}
