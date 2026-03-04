package collector

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/nimishgj/aws-radar/internal/metrics"
)

func TestCollector_ses_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewSESCollector(), "ses")
}

func TestCollector_ses_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewSESCollector(), true)
}

func TestCollector_ses_Metrics(t *testing.T) {
	server := newMockAWSServer([]mockRoute{
		{
			// GET /v2/email/identities/{name} (GetEmailIdentity) — must match before ListEmailIdentities
			matcher: func(r *http.Request) bool {
				return r.Method == "GET" &&
					strings.HasPrefix(r.URL.Path, "/v2/email/identities/") &&
					r.URL.Path != "/v2/email/identities"
			},
			body: `{
				"IdentityType": "DOMAIN",
				"VerifiedForSendingStatus": true,
				"DkimAttributes": {"Status": "SUCCESS"},
				"MailFromAttributes": {"MailFromDomainStatus": "SUCCESS"}
			}`,
		},
		{
			// GET /v2/email/identities (ListEmailIdentities)
			matcher: func(r *http.Request) bool {
				return r.Method == "GET" && r.URL.Path == "/v2/email/identities"
			},
			body: `{
				"EmailIdentities": [
					{"IdentityName": "example.com", "IdentityType": "DOMAIN", "SendingEnabled": true, "VerificationStatus": "SUCCESS"},
					{"IdentityName": "test@example.com", "IdentityType": "EMAIL_ADDRESS", "SendingEnabled": false, "VerificationStatus": "PENDING"}
				]
			}`,
		},
		{
			// GET /v2/email/configuration-sets/{name}/event-destinations
			matcher: func(r *http.Request) bool {
				return r.Method == "GET" && strings.Contains(r.URL.Path, "/event-destinations")
			},
			body: `{
				"EventDestinations": [
					{"Name": "cw-dest", "Enabled": true, "CloudWatchDestination": {"DimensionConfigurations": []}},
					{"Name": "sns-dest", "Enabled": true, "SnsDestination": {"TopicArn": "arn:aws:sns:us-east-1:123:topic"}}
				]
			}`,
		},
		{
			// GET /v2/email/configuration-sets/{name} (GetConfigurationSet)
			matcher: func(r *http.Request) bool {
				return r.Method == "GET" &&
					strings.HasPrefix(r.URL.Path, "/v2/email/configuration-sets/") &&
					!strings.Contains(r.URL.Path, "/event-destinations")
			},
			body: `{
				"ConfigurationSetName": "my-config-set",
				"DeliveryOptions": {"SendingPoolName": "my-pool"}
			}`,
		},
		{
			// GET /v2/email/configuration-sets (ListConfigurationSets)
			matcher: func(r *http.Request) bool {
				return r.Method == "GET" && r.URL.Path == "/v2/email/configuration-sets"
			},
			body: `{"ConfigurationSets": ["my-config-set"]}`,
		},
		{
			// GET /v2/email/contact-lists
			matcher: restPath("/v2/email/contact-lists"),
			body:    `{"ContactLists": [{"ContactListName": "list1"}, {"ContactListName": "list2"}]}`,
		},
		{
			// GET /v2/email/suppression/addresses
			matcher: restPath("/v2/email/suppression/addresses"),
			body: `{
				"SuppressedDestinationSummaries": [
					{"EmailAddress": "bounce@example.com", "Reason": "BOUNCE"},
					{"EmailAddress": "complaint@example.com", "Reason": "COMPLAINT"},
					{"EmailAddress": "bounce2@example.com", "Reason": "BOUNCE"}
				]
			}`,
		},
		{
			// GET /v2/email/dedicated-ip-pools
			matcher: restPath("/v2/email/dedicated-ip-pools"),
			body:    `{"DedicatedIpPools": ["pool1", "pool2"]}`,
		},
		{
			// GET /v2/email/account
			matcher: restPath("/v2/email/account"),
			body: `{
				"SendingEnabled": true,
				"ProductionAccessEnabled": true,
				"DedicatedIpAutoWarmupEnabled": false,
				"SuppressionAttributes": {"SuppressedReasons": ["BOUNCE"]},
				"SendQuota": {"Max24HourSend": 50000, "MaxSendRate": 14, "SentLast24Hours": 1234}
			}`,
		},
	})
	defer server.Close()

	metrics.ResetAll()
	cfg := mockAWSConfig(server.URL, "us-east-1")
	err := NewSESCollector().Collect(context.Background(), cfg, "us-east-1", "123456789012", "test-account")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Total identities
	if v := gaugeValue(metrics.SESIdentities, "123456789012", "test-account", "us-east-1"); v != 2 {
		t.Errorf("SESIdentities: expected 2, got %v", v)
	}

	// Verification status
	if v := gaugeValue(metrics.SESIdentitiesByVerificationStatus, "123456789012", "test-account", "us-east-1", "SUCCESS"); v != 1 {
		t.Errorf("SESIdentitiesByVerificationStatus SUCCESS: expected 1, got %v", v)
	}
	if v := gaugeValue(metrics.SESIdentitiesByVerificationStatus, "123456789012", "test-account", "us-east-1", "PENDING"); v != 1 {
		t.Errorf("SESIdentitiesByVerificationStatus PENDING: expected 1, got %v", v)
	}

	// Auth status: both identities get same GetEmailIdentity response → 2x SUCCESS/SUCCESS/SUCCESS
	if v := gaugeValue(metrics.SESIdentityAuthStatus, "123456789012", "test-account", "us-east-1", "SUCCESS", "SUCCESS", "SUCCESS"); v != 2 {
		t.Errorf("SESIdentityAuthStatus SUCCESS/SUCCESS/SUCCESS: expected 2, got %v", v)
	}

	// Config sets
	if v := gaugeValue(metrics.SESConfigSets, "123456789012", "test-account", "us-east-1"); v != 1 {
		t.Errorf("SESConfigSets: expected 1, got %v", v)
	}

	// Event destinations
	if v := gaugeValue(metrics.SESConfigSetEventDestinations, "123456789012", "test-account", "us-east-1", "cloudwatch"); v != 1 {
		t.Errorf("SESConfigSetEventDestinations cloudwatch: expected 1, got %v", v)
	}
	if v := gaugeValue(metrics.SESConfigSetEventDestinations, "123456789012", "test-account", "us-east-1", "sns"); v != 1 {
		t.Errorf("SESConfigSetEventDestinations sns: expected 1, got %v", v)
	}

	// Config sets by sending pool
	if v := gaugeValue(metrics.SESConfigSetsBySendingPool, "123456789012", "test-account", "us-east-1", "my-pool"); v != 1 {
		t.Errorf("SESConfigSetsBySendingPool my-pool: expected 1, got %v", v)
	}

	// Contact lists
	if v := gaugeValue(metrics.SESContactLists, "123456789012", "test-account", "us-east-1"); v != 2 {
		t.Errorf("SESContactLists: expected 2, got %v", v)
	}

	// Suppressed destinations
	if v := gaugeValue(metrics.SESSuppressedDestinations, "123456789012", "test-account", "us-east-1", "BOUNCE"); v != 2 {
		t.Errorf("SESSuppressedDestinations BOUNCE: expected 2, got %v", v)
	}
	if v := gaugeValue(metrics.SESSuppressedDestinations, "123456789012", "test-account", "us-east-1", "COMPLAINT"); v != 1 {
		t.Errorf("SESSuppressedDestinations COMPLAINT: expected 1, got %v", v)
	}

	// Dedicated IP pools
	if v := gaugeValue(metrics.SESDedicatedIPPools, "123456789012", "test-account", "us-east-1"); v != 2 {
		t.Errorf("SESDedicatedIPPools: expected 2, got %v", v)
	}

	// Account settings
	if v := gaugeValue(metrics.SESAccountSettings, "123456789012", "test-account", "us-east-1", "sending_enabled"); v != 1 {
		t.Errorf("SESAccountSettings sending_enabled: expected 1, got %v", v)
	}
	if v := gaugeValue(metrics.SESAccountSettings, "123456789012", "test-account", "us-east-1", "production_access_enabled"); v != 1 {
		t.Errorf("SESAccountSettings production_access_enabled: expected 1, got %v", v)
	}
	if v := gaugeValue(metrics.SESAccountSettings, "123456789012", "test-account", "us-east-1", "dedicated_ip_auto_warmup_enabled"); v != 0 {
		t.Errorf("SESAccountSettings dedicated_ip_auto_warmup_enabled: expected 0, got %v", v)
	}
	if v := gaugeValue(metrics.SESAccountSettings, "123456789012", "test-account", "us-east-1", "suppression_enabled"); v != 1 {
		t.Errorf("SESAccountSettings suppression_enabled: expected 1, got %v", v)
	}

	// Sending quota
	if v := gaugeValue(metrics.SESSendingQuota, "123456789012", "test-account", "us-east-1", "max_24_hour_send"); v != 50000 {
		t.Errorf("SESSendingQuota max_24_hour_send: expected 50000, got %v", v)
	}
	if v := gaugeValue(metrics.SESSendingQuota, "123456789012", "test-account", "us-east-1", "max_send_rate"); v != 14 {
		t.Errorf("SESSendingQuota max_send_rate: expected 14, got %v", v)
	}
	if v := gaugeValue(metrics.SESSendingQuota, "123456789012", "test-account", "us-east-1", "sent_last_24_hours"); v != 1234 {
		t.Errorf("SESSendingQuota sent_last_24_hours: expected 1234, got %v", v)
	}
}
