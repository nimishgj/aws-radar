package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type ServiceCatalogCollector struct{}

func NewServiceCatalogCollector() *ServiceCatalogCollector {
	return &ServiceCatalogCollector{}
}

func (c *ServiceCatalogCollector) Name() string {
	return "servicecatalog"
}

func (c *ServiceCatalogCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := servicecatalog.NewFromConfig(cfg)

	// List Portfolios
	var portfolioCount float64
	paginator := servicecatalog.NewListPortfoliosPaginator(client, &servicecatalog.ListPortfoliosInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		portfolioCount += float64(len(page.PortfolioDetails))
	}
	metrics.ServiceCatalogPortfolios.WithLabelValues(account, accountName, region).Set(portfolioCount)

	// Search Products
	var productCount float64
	var productToken *string
	for {
		output, err := client.SearchProductsAsAdmin(ctx, &servicecatalog.SearchProductsAsAdminInput{
			PageToken: productToken,
		})
		if err != nil {
			log.Warn().Err(err).Str("region", region).Msg("Failed to search Service Catalog products")
			break
		}
		productCount += float64(len(output.ProductViewDetails))
		if output.NextPageToken == nil || *output.NextPageToken == "" {
			break
		}
		productToken = output.NextPageToken
	}
	metrics.ServiceCatalogProducts.WithLabelValues(account, accountName, region).Set(productCount)

	// Search Provisioned Products
	provisionedCounts := make(map[string]float64)
	var ppToken *string
	for {
		output, err := client.SearchProvisionedProducts(ctx, &servicecatalog.SearchProvisionedProductsInput{
			PageToken: ppToken,
		})
		if err != nil {
			log.Warn().Err(err).Str("region", region).Msg("Failed to search Service Catalog provisioned products")
			break
		}
		for _, pp := range output.ProvisionedProducts {
			status := string(pp.Status)
			if status == "" {
				status = "unknown"
			}
			provisionedCounts[status]++
		}
		if output.NextPageToken == nil || *output.NextPageToken == "" {
			break
		}
		ppToken = output.NextPageToken
	}

	for status, count := range provisionedCounts {
		metrics.ServiceCatalogProvisionedProducts.WithLabelValues(account, accountName, region, status).Set(count)
	}

	log.Debug().
		Str("region", region).
		Float64("portfolios", portfolioCount).
		Float64("products", productCount).
		Int("provisioned_statuses", len(provisionedCounts)).
		Msg("Service Catalog collection completed")

	return nil
}
