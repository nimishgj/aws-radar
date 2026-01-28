package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type VPCCollector struct{}

func NewVPCCollector() *VPCCollector {
	return &VPCCollector{}
}

func (c *VPCCollector) Name() string {
	return "vpc"
}

func (c *VPCCollector) Collect(ctx context.Context, cfg aws.Config, region string) error {
	client := ec2.NewFromConfig(cfg)

	// Collect VPCs
	if err := c.collectVPCs(ctx, client, region); err != nil {
		log.Warn().Err(err).Str("region", region).Msg("Failed to collect VPCs")
	}

	// Collect Subnets
	if err := c.collectSubnets(ctx, client, region); err != nil {
		log.Warn().Err(err).Str("region", region).Msg("Failed to collect Subnets")
	}

	// Collect Security Groups
	if err := c.collectSecurityGroups(ctx, client, region); err != nil {
		log.Warn().Err(err).Str("region", region).Msg("Failed to collect Security Groups")
	}

	// Collect NAT Gateways
	if err := c.collectNATGateways(ctx, client, region); err != nil {
		log.Warn().Err(err).Str("region", region).Msg("Failed to collect NAT Gateways")
	}

	// Collect Internet Gateways
	if err := c.collectInternetGateways(ctx, client, region); err != nil {
		log.Warn().Err(err).Str("region", region).Msg("Failed to collect Internet Gateways")
	}

	return nil
}

func (c *VPCCollector) collectVPCs(ctx context.Context, client *ec2.Client, region string) error {
	counts := make(map[string]float64)

	paginator := ec2.NewDescribeVpcsPaginator(client, &ec2.DescribeVpcsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, vpc := range page.Vpcs {
			state := string(vpc.State)
			counts[state]++
		}
	}

	for state, count := range counts {
		metrics.VPCs.WithLabelValues(region, state).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("vpc_states", len(counts)).
		Msg("VPC collection completed")

	return nil
}

func (c *VPCCollector) collectSubnets(ctx context.Context, client *ec2.Client, region string) error {
	counts := make(map[string]float64)

	paginator := ec2.NewDescribeSubnetsPaginator(client, &ec2.DescribeSubnetsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, subnet := range page.Subnets {
			az := aws.ToString(subnet.AvailabilityZone)
			counts[az]++
		}
	}

	for az, count := range counts {
		metrics.Subnets.WithLabelValues(region, az).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("subnet_azs", len(counts)).
		Msg("Subnet collection completed")

	return nil
}

func (c *VPCCollector) collectSecurityGroups(ctx context.Context, client *ec2.Client, region string) error {
	counts := make(map[string]float64)

	paginator := ec2.NewDescribeSecurityGroupsPaginator(client, &ec2.DescribeSecurityGroupsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, sg := range page.SecurityGroups {
			vpcId := aws.ToString(sg.VpcId)
			if vpcId == "" {
				vpcId = "ec2-classic"
			}
			counts[vpcId]++
		}
	}

	for vpcId, count := range counts {
		metrics.SecurityGroups.WithLabelValues(region, vpcId).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("security_group_vpcs", len(counts)).
		Msg("Security Group collection completed")

	return nil
}

func (c *VPCCollector) collectNATGateways(ctx context.Context, client *ec2.Client, region string) error {
	counts := make(map[string]float64)

	paginator := ec2.NewDescribeNatGatewaysPaginator(client, &ec2.DescribeNatGatewaysInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, nat := range page.NatGateways {
			state := string(nat.State)
			counts[state]++
		}
	}

	for state, count := range counts {
		metrics.NATGateways.WithLabelValues(region, state).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("nat_gateway_states", len(counts)).
		Msg("NAT Gateway collection completed")

	return nil
}

func (c *VPCCollector) collectInternetGateways(ctx context.Context, client *ec2.Client, region string) error {
	var count float64 = 0

	paginator := ec2.NewDescribeInternetGatewaysPaginator(client, &ec2.DescribeInternetGatewaysInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		count += float64(len(page.InternetGateways))
	}

	metrics.InternetGateways.WithLabelValues(region).Set(count)

	log.Debug().
		Str("region", region).
		Float64("internet_gateway_count", count).
		Msg("Internet Gateway collection completed")

	return nil
}
