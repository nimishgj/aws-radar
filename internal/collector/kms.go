package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type KMSCollector struct{}

func NewKMSCollector() *KMSCollector {
	return &KMSCollector{}
}

func (c *KMSCollector) Name() string {
	return "kms"
}

func (c *KMSCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := kms.NewFromConfig(cfg)
	paginator := kms.NewListKeysPaginator(client, &kms.ListKeysInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Keys))
	}

	metrics.KMSKeys.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
