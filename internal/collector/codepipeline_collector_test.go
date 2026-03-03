package collector

import "testing"

func TestCollector_codepipeline_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewCodePipelineCollector(), "codepipeline")
}

func TestCollector_codepipeline_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewCodePipelineCollector(), true)
}
