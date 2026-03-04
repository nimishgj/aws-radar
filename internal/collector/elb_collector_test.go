package collector

import (
	"context"
	"testing"

	"github.com/nimishgj/aws-radar/internal/metrics"
)

func TestCollector_elb_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewELBCollector(), "elb")
}

func TestCollector_elb_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewELBCollector(), false)
}

func TestCollector_elb_Metrics(t *testing.T) {
	// ELB Classic uses Query protocol (DescribeLoadBalancers)
	// ELBv2 uses Query protocol (DescribeLoadBalancers, DescribeListeners, DescribeRules, DescribeTargetGroups)
	server := newMockAWSServer([]mockRoute{
		{
			// Classic ELB - DescribeLoadBalancers (elasticloadbalancing uses same Action name)
			// We match on the classic API by checking the body doesn't contain the v2 marker.
			// Since both use Action=DescribeLoadBalancers, we rely on request order:
			// classic is called first. We'll use two separate handlers below.
			matcher: queryAction("DescribeLoadBalancers"),
			body: `<DescribeLoadBalancersResponse>
  <DescribeLoadBalancersResult>
    <LoadBalancerDescriptions>
      <member>
        <LoadBalancerName>classic-lb-1</LoadBalancerName>
        <Scheme>internet-facing</Scheme>
      </member>
      <member>
        <LoadBalancerName>classic-lb-2</LoadBalancerName>
        <Scheme>internal</Scheme>
      </member>
    </LoadBalancerDescriptions>
  </DescribeLoadBalancersResult>
</DescribeLoadBalancersResponse>`,
		},
		{
			matcher: queryAction("DescribeListeners"),
			body: `<DescribeListenersResponse>
  <DescribeListenersResult>
    <Listeners>
      <member>
        <ListenerArn>arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-alb/1234/5678</ListenerArn>
        <Protocol>HTTPS</Protocol>
      </member>
      <member>
        <ListenerArn>arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-alb/1234/9012</ListenerArn>
        <Protocol>HTTP</Protocol>
      </member>
    </Listeners>
  </DescribeListenersResult>
</DescribeListenersResponse>`,
		},
		{
			matcher: queryAction("DescribeRules"),
			body: `<DescribeRulesResponse>
  <DescribeRulesResult>
    <Rules>
      <member>
        <RuleArn>arn:aws:elasticloadbalancing:us-east-1:123456789012:listener-rule/app/my-alb/1234/5678/rule1</RuleArn>
        <IsDefault>true</IsDefault>
      </member>
      <member>
        <RuleArn>arn:aws:elasticloadbalancing:us-east-1:123456789012:listener-rule/app/my-alb/1234/5678/rule2</RuleArn>
        <IsDefault>false</IsDefault>
      </member>
    </Rules>
  </DescribeRulesResult>
</DescribeRulesResponse>`,
		},
		{
			matcher: queryAction("DescribeTargetGroups"),
			body: `<DescribeTargetGroupsResponse>
  <DescribeTargetGroupsResult>
    <TargetGroups>
      <member>
        <TargetGroupArn>arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/tg1/1234</TargetGroupArn>
        <TargetType>instance</TargetType>
        <LoadBalancerArns>
          <member>arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-alb/1234</member>
        </LoadBalancerArns>
      </member>
      <member>
        <TargetGroupArn>arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/tg2/5678</TargetGroupArn>
        <TargetType>ip</TargetType>
        <LoadBalancerArns>
          <member>arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-alb/1234</member>
        </LoadBalancerArns>
      </member>
    </TargetGroups>
  </DescribeTargetGroupsResult>
</DescribeTargetGroupsResponse>`,
		},
	})
	defer server.Close()

	metrics.ResetAll()

	// Note: Both classic and v2 share DescribeLoadBalancers action name but hit different
	// service endpoints. With a single mock server the classic call will get the same XML.
	// The classic collector parses <LoadBalancerDescriptions> and the v2 collector parses
	// <LoadBalancers>. Since v2 won't find <LoadBalancers>, it'll get 0 results — which
	// means we can only fully test classic from a single mock, OR we test them independently.
	// Let's test them independently for accuracy.

	cfg := mockAWSConfig(server.URL, "us-east-1")
	err := NewELBCollector().Collect(context.Background(), cfg, "us-east-1", "123456789012", "test-account")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Classic LBs: 1 internet-facing, 1 internal
	if v := gaugeValue(metrics.ELBClassic, "123456789012", "test-account", "us-east-1", "internet-facing"); v != 1 {
		t.Errorf("ELBClassic internet-facing: expected 1, got %v", v)
	}
	if v := gaugeValue(metrics.ELBClassic, "123456789012", "test-account", "us-east-1", "internal"); v != 1 {
		t.Errorf("ELBClassic internal: expected 1, got %v", v)
	}
}

func TestCollector_elbv2_Metrics(t *testing.T) {
	// Dedicated test for ELBv2 with proper response format
	server := newMockAWSServer([]mockRoute{
		{
			matcher: queryAction("DescribeLoadBalancers"),
			body: `<DescribeLoadBalancersResponse>
  <DescribeLoadBalancersResult>
    <LoadBalancers>
      <member>
        <LoadBalancerArn>arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-alb/1234</LoadBalancerArn>
        <LoadBalancerName>my-alb</LoadBalancerName>
        <Type>application</Type>
        <Scheme>internet-facing</Scheme>
        <IpAddressType>ipv4</IpAddressType>
        <State><Code>active</Code></State>
        <AvailabilityZones>
          <member><ZoneName>us-east-1a</ZoneName><SubnetId>subnet-aaa</SubnetId></member>
          <member><ZoneName>us-east-1b</ZoneName><SubnetId>subnet-bbb</SubnetId></member>
        </AvailabilityZones>
      </member>
      <member>
        <LoadBalancerArn>arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/net/my-nlb/5678</LoadBalancerArn>
        <LoadBalancerName>my-nlb</LoadBalancerName>
        <Type>network</Type>
        <Scheme>internal</Scheme>
        <IpAddressType>ipv4</IpAddressType>
        <State><Code>active</Code></State>
        <AvailabilityZones>
          <member><ZoneName>us-east-1a</ZoneName><SubnetId>subnet-aaa</SubnetId></member>
        </AvailabilityZones>
      </member>
    </LoadBalancers>
  </DescribeLoadBalancersResult>
</DescribeLoadBalancersResponse>`,
		},
		{
			matcher: queryAction("DescribeListeners"),
			body: `<DescribeListenersResponse>
  <DescribeListenersResult>
    <Listeners>
      <member>
        <ListenerArn>arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-alb/1234/listener1</ListenerArn>
        <Protocol>HTTPS</Protocol>
      </member>
    </Listeners>
  </DescribeListenersResult>
</DescribeListenersResponse>`,
		},
		{
			matcher: queryAction("DescribeRules"),
			body: `<DescribeRulesResponse>
  <DescribeRulesResult>
    <Rules>
      <member><RuleArn>rule1</RuleArn></member>
      <member><RuleArn>rule2</RuleArn></member>
      <member><RuleArn>rule3</RuleArn></member>
    </Rules>
  </DescribeRulesResult>
</DescribeRulesResponse>`,
		},
		{
			matcher: queryAction("DescribeTargetGroups"),
			body: `<DescribeTargetGroupsResponse>
  <DescribeTargetGroupsResult>
    <TargetGroups>
      <member>
        <TargetType>instance</TargetType>
        <LoadBalancerArns>
          <member>arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-alb/1234</member>
        </LoadBalancerArns>
      </member>
      <member>
        <TargetType>ip</TargetType>
        <LoadBalancerArns>
          <member>arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/net/my-nlb/5678</member>
        </LoadBalancerArns>
      </member>
    </TargetGroups>
  </DescribeTargetGroupsResult>
</DescribeTargetGroupsResponse>`,
		},
	})
	defer server.Close()

	metrics.ResetAll()
	cfg := mockAWSConfig(server.URL, "us-east-1")

	collector := NewELBCollector()
	err := collector.collectV2(context.Background(), cfg, "us-east-1", "123456789012", "test-account")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ELBV2 counts: 1 application/internet-facing, 1 network/internal
	if v := gaugeValue(metrics.ELBV2, "123456789012", "test-account", "us-east-1", "application", "internet-facing"); v != 1 {
		t.Errorf("ELBV2 application/internet-facing: expected 1, got %v", v)
	}
	if v := gaugeValue(metrics.ELBV2, "123456789012", "test-account", "us-east-1", "network", "internal"); v != 1 {
		t.Errorf("ELBV2 network/internal: expected 1, got %v", v)
	}

	// Detailed counts
	if v := gaugeValue(metrics.ELBV2Detailed, "123456789012", "test-account", "us-east-1", "application", "internet-facing", "ipv4", "active"); v != 1 {
		t.Errorf("ELBV2Detailed application: expected 1, got %v", v)
	}

	// Listeners: both LBs get 1 HTTPS listener each from the mock (same response for both)
	// ALB gets HTTPS listener
	if v := gaugeValue(metrics.ELBV2Listeners, "123456789012", "test-account", "us-east-1", "application", "internet-facing", "HTTPS"); v < 1 {
		t.Errorf("ELBV2Listeners application/HTTPS: expected >= 1, got %v", v)
	}

	// Rules per ALB (only counted for application type)
	if v := gaugeValue(metrics.ELBV2RulesPerALB, "123456789012", "test-account", "us-east-1", "my-alb"); v < 1 {
		t.Errorf("ELBV2RulesPerALB my-alb: expected >= 1, got %v", v)
	}

	// AZ count per LB: my-alb has 2 AZs
	if v := gaugeValue(metrics.ELBV2AvailabilityZonesPerLB, "123456789012", "test-account", "us-east-1", "my-alb", "application", "internet-facing"); v != 2 {
		t.Errorf("ELBV2AvailabilityZonesPerLB my-alb: expected 2, got %v", v)
	}

	// Subnets per LB: my-alb has 2 subnets
	if v := gaugeValue(metrics.ELBV2SubnetsPerLB, "123456789012", "test-account", "us-east-1", "my-alb", "application", "internet-facing"); v != 2 {
		t.Errorf("ELBV2SubnetsPerLB my-alb: expected 2, got %v", v)
	}

	// my-nlb has 1 AZ
	if v := gaugeValue(metrics.ELBV2AvailabilityZonesPerLB, "123456789012", "test-account", "us-east-1", "my-nlb", "network", "internal"); v != 1 {
		t.Errorf("ELBV2AvailabilityZonesPerLB my-nlb: expected 1, got %v", v)
	}

	// Target groups: 1 instance (application), 1 ip (network)
	if v := gaugeValue(metrics.ELBV2TargetGroups, "123456789012", "test-account", "us-east-1", "application", "instance"); v != 1 {
		t.Errorf("ELBV2TargetGroups application/instance: expected 1, got %v", v)
	}
	if v := gaugeValue(metrics.ELBV2TargetGroups, "123456789012", "test-account", "us-east-1", "network", "ip"); v != 1 {
		t.Errorf("ELBV2TargetGroups network/ip: expected 1, got %v", v)
	}
}
