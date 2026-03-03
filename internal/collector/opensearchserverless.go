package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type OpenSearchServerlessCollector struct{}

func NewOpenSearchServerlessCollector() *OpenSearchServerlessCollector {
	return &OpenSearchServerlessCollector{}
}

func (c *OpenSearchServerlessCollector) Name() string { return "opensearchserverless" }

func (c *OpenSearchServerlessCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := opensearchserverless.NewFromConfig(cfg)
	paginator := opensearchserverless.NewListCollectionsPaginator(client, &opensearchserverless.ListCollectionsInput{})
	counts := make(map[string]float64)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, col := range page.CollectionSummaries {
			status := string(col.Status)
			if status == "" {
				status = "UNKNOWN"
			}
			typeName := aws.ToString(col.CollectionGroupName)
			if typeName == "" {
				typeName = "UNKNOWN"
			}
			counts[status+"|"+typeName]++
		}
	}
	for key, count := range counts {
		parts := splitKey(key, 2)
		metrics.OpenSearchServerlessCollections.WithLabelValues(account, accountName, region, parts[0], parts[1]).Set(count)
	}
	return nil
}
