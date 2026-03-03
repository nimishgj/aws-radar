package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type BackupCollector struct{}

func NewBackupCollector() *BackupCollector {
	return &BackupCollector{}
}

func (c *BackupCollector) Name() string {
	return "backup"
}

func (c *BackupCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := backup.NewFromConfig(cfg)
	paginator := backup.NewListBackupVaultsPaginator(client, &backup.ListBackupVaultsInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.BackupVaultList))
	}

	metrics.BackupVaults.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
