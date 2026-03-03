package collector

import "testing"

func TestCollector_bedrock_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewBedrockCollector(), "bedrock")
}

func TestCollector_bedrock_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewBedrockCollector(), true)
}
