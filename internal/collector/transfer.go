package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/transfer"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type TransferCollector struct{}

func NewTransferCollector() *TransferCollector {
	return &TransferCollector{}
}

func (c *TransferCollector) Name() string {
	return "transfer"
}

func (c *TransferCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := transfer.NewFromConfig(cfg)
	paginator := transfer.NewListServersPaginator(client, &transfer.ListServersInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Servers))
	}

	metrics.TransferServers.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
