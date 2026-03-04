package collector

import (
	"testing"
)

func TestCollector_cognito_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewCognitoCollector(), "cognito")
}

func TestCollector_cognito_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewCognitoCollector(), true)
}
