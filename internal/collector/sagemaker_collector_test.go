package collector

import "testing"

func TestCollector_sagemaker_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewSageMakerCollector(), "sagemaker")
}

func TestCollector_sagemaker_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewSageMakerCollector(), true)
}
