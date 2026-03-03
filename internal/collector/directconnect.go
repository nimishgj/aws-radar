package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directconnect"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type DirectConnectCollector struct{}

func NewDirectConnectCollector() *DirectConnectCollector { return &DirectConnectCollector{} }

func (c *DirectConnectCollector) Name() string { return "directconnect" }

func (c *DirectConnectCollector) Collect(ctx context.Context, cfg aws.Config, account, accountName string) error {
	cfg.Region = "us-east-1"
	client := directconnect.NewFromConfig(cfg)

	output, err := client.DescribeConnections(ctx, &directconnect.DescribeConnectionsInput{})
	if err != nil {
		return err
	}

	counts := make(map[string]float64)
	for _, conn := range output.Connections {
		state := string(conn.ConnectionState)
		if state == "" {
			state = "unknown"
		}
		counts[state]++
	}
	for state, count := range counts {
		metrics.DirectConnectConnections.WithLabelValues(account, accountName, state).Set(count)
	}
	return nil
}
