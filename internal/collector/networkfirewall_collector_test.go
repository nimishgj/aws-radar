package collector

import (
	"testing"
)

func TestCollector_networkfirewall_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewNetworkFirewallCollector(), "networkfirewall")
}

func TestCollector_networkfirewall_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewNetworkFirewallCollector(), true)
}
