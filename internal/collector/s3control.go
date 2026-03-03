package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type S3ControlCollector struct{}

func NewS3ControlCollector() *S3ControlCollector { return &S3ControlCollector{} }

func (c *S3ControlCollector) Name() string { return "s3control" }

func (c *S3ControlCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := s3control.NewFromConfig(cfg)

	apPaginator := s3control.NewListAccessPointsPaginator(client, &s3control.ListAccessPointsInput{AccountId: aws.String(account)})
	var accessPointCount float64
	for apPaginator.HasMorePages() {
		page, err := apPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		accessPointCount += float64(len(page.AccessPointList))
	}
	metrics.S3AccessPoints.WithLabelValues(account, accountName, region).Set(accessPointCount)

	slPaginator := s3control.NewListStorageLensConfigurationsPaginator(client, &s3control.ListStorageLensConfigurationsInput{AccountId: aws.String(account)})
	var storageLensCount float64
	for slPaginator.HasMorePages() {
		page, err := slPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		storageLensCount += float64(len(page.StorageLensConfigurationList))
	}
	metrics.S3StorageLensConfigurations.WithLabelValues(account, accountName, region).Set(storageLensCount)

	return nil
}
