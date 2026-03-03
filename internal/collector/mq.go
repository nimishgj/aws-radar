package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mq"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type MQCollector struct{}

func NewMQCollector() *MQCollector {
	return &MQCollector{}
}

func (c *MQCollector) Name() string {
	return "mq"
}

func (c *MQCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := mq.NewFromConfig(cfg)
	paginator := mq.NewListBrokersPaginator(client, &mq.ListBrokersInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.BrokerSummaries))
	}

	metrics.MQBrokers.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
