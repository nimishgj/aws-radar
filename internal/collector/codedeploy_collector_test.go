package collector

import "testing"

func TestCollector_codedeploy_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewCodeDeployCollector(), "codedeploy")
}

func TestCollector_codedeploy_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewCodeDeployCollector(), true)
}
