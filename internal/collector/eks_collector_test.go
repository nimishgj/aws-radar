package collector

import "testing"

func TestCollector_eks_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewEKSCollector(), "eks")
}

func TestCollector_eks_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewEKSCollector(), true)
}
