package collector

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/shield"
	shieldTypes "github.com/aws/aws-sdk-go-v2/service/shield/types"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type ShieldCollector struct{}

func NewShieldCollector() *ShieldCollector {
	return &ShieldCollector{}
}

func (c *ShieldCollector) Name() string {
	return "shield"
}

func (c *ShieldCollector) Collect(ctx context.Context, cfg aws.Config, account, accountName string) error {
	globalCfg := cfg.Copy()
	if globalCfg.Region == "" {
		globalCfg.Region = "us-east-1"
	}

	client := shield.NewFromConfig(globalCfg)
	_, err := client.DescribeSubscription(ctx, &shield.DescribeSubscriptionInput{})
	if err != nil {
		var notFound *shieldTypes.ResourceNotFoundException
		if errors.As(err, &notFound) {
			metrics.ShieldSubscriptions.WithLabelValues(account, accountName).Set(0)
			return nil
		}
		return err
	}

	metrics.ShieldSubscriptions.WithLabelValues(account, accountName).Set(1)
	return nil
}
