package collector

import "testing"

func TestCollector_codebuild_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewCodeBuildCollector(), "codebuild")
}

func TestCollector_codebuild_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewCodeBuildCollector(), true)
}
