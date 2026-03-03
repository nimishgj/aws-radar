package collector

import "testing"

func TestCollector_dynamodb_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewDynamoDBCollector(), "dynamodb")
}

func TestCollector_dynamodb_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewDynamoDBCollector(), true)
}
