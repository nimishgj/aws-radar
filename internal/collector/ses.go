package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type SESCollector struct{}

func NewSESCollector() *SESCollector {
	return &SESCollector{}
}

func (c *SESCollector) Name() string {
	return "ses"
}

func (c *SESCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := sesv2.NewFromConfig(cfg)
	paginator := sesv2.NewListEmailIdentitiesPaginator(client, &sesv2.ListEmailIdentitiesInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.EmailIdentities))
	}

	metrics.SESIdentities.WithLabelValues(account, accountName, region).Set(count)

	var configSetCount float64
	configSetPaginator := sesv2.NewListConfigurationSetsPaginator(client, &sesv2.ListConfigurationSetsInput{})
	for configSetPaginator.HasMorePages() {
		page, err := configSetPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		configSetCount += float64(len(page.ConfigurationSets))
	}
	metrics.SESConfigSets.WithLabelValues(account, accountName, region).Set(configSetCount)

	var contactListCount float64
	contactListPaginator := sesv2.NewListContactListsPaginator(client, &sesv2.ListContactListsInput{})
	for contactListPaginator.HasMorePages() {
		page, err := contactListPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		contactListCount += float64(len(page.ContactLists))
	}
	metrics.SESContactLists.WithLabelValues(account, accountName, region).Set(contactListCount)
	return nil
}
