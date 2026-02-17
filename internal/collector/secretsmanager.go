package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type SecretsManagerCollector struct{}

func NewSecretsManagerCollector() *SecretsManagerCollector {
	return &SecretsManagerCollector{}
}

func (c *SecretsManagerCollector) Name() string {
	return "secretsmanager"
}

func (c *SecretsManagerCollector) Collect(ctx context.Context, cfg aws.Config, region, account string) error {
	client := secretsmanager.NewFromConfig(cfg)

	var count float64
	paginator := secretsmanager.NewListSecretsPaginator(client, &secretsmanager.ListSecretsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.SecretList))
	}

	metrics.SecretsManagerSecrets.WithLabelValues(account, region).Set(count)

	log.Debug().
		Str("region", region).
		Float64("secret_count", count).
		Msg("Secrets Manager collection completed")

	return nil
}
