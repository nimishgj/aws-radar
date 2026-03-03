package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type TimestreamCollector struct{}

func NewTimestreamCollector() *TimestreamCollector {
	return &TimestreamCollector{}
}

func (c *TimestreamCollector) Name() string {
	return "timestream"
}

func (c *TimestreamCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := timestreamwrite.NewFromConfig(cfg)

	dbPaginator := timestreamwrite.NewListDatabasesPaginator(client, &timestreamwrite.ListDatabasesInput{})
	var dbCount float64
	var tableCount float64

	for dbPaginator.HasMorePages() {
		page, err := dbPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		dbCount += float64(len(page.Databases))

		for _, db := range page.Databases {
			tablePaginator := timestreamwrite.NewListTablesPaginator(client, &timestreamwrite.ListTablesInput{
				DatabaseName: db.DatabaseName,
			})
			for tablePaginator.HasMorePages() {
				tablePage, err := tablePaginator.NextPage(ctx)
				if err != nil {
					return err
				}
				tableCount += float64(len(tablePage.Tables))
			}
		}
	}

	metrics.TimestreamDatabases.WithLabelValues(account, accountName, region).Set(dbCount)
	metrics.TimestreamTables.WithLabelValues(account, accountName, region).Set(tableCount)
	return nil
}
