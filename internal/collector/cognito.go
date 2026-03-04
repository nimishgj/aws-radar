package collector

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognitoidentity "github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	cognitoidp "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type CognitoCollector struct{}

func NewCognitoCollector() *CognitoCollector {
	return &CognitoCollector{}
}

func (c *CognitoCollector) Name() string {
	return "cognito"
}

func (c *CognitoCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	// Collect User Pools
	idpClient := cognitoidp.NewFromConfig(cfg)
	statusCounts := make(map[string]float64)

	var nextToken *string
	for {
		output, err := idpClient.ListUserPools(ctx, &cognitoidp.ListUserPoolsInput{
			MaxResults: aws.Int32(60),
			NextToken:  nextToken,
		})
		if err != nil {
			return err
		}

		for _, pool := range output.UserPools {
			status := string(pool.Status)
			if status == "" {
				status = "unknown"
			}
			statusCounts[status]++

			// Get estimated user count
			poolName := aws.ToString(pool.Name)
			if poolName == "" {
				poolName = "unknown"
			}
			descOutput, descErr := idpClient.DescribeUserPool(ctx, &cognitoidp.DescribeUserPoolInput{
				UserPoolId: pool.Id,
			})
			if descErr != nil {
				log.Warn().Err(descErr).Str("region", region).Str("pool", poolName).Msg("Failed to describe Cognito user pool")
				continue
			}
			if descOutput.UserPool != nil {
				metrics.CognitoUserPoolUsers.WithLabelValues(account, accountName, region, poolName).Set(float64(descOutput.UserPool.EstimatedNumberOfUsers))
			}
		}

		if output.NextToken == nil || *output.NextToken == "" {
			break
		}
		nextToken = output.NextToken
	}

	for status, count := range statusCounts {
		metrics.CognitoUserPools.WithLabelValues(account, accountName, region, status).Set(count)
	}

	// Collect Identity Pools
	identityClient := cognitoidentity.NewFromConfig(cfg)
	var identityCount float64
	var identityNextToken *string
	for {
		output, err := identityClient.ListIdentityPools(ctx, &cognitoidentity.ListIdentityPoolsInput{
			MaxResults: aws.Int32(60),
			NextToken:  identityNextToken,
		})
		if err != nil {
			log.Warn().Err(err).Str("region", region).Msg("Failed to list Cognito identity pools")
			break
		}
		identityCount += float64(len(output.IdentityPools))
		if output.NextToken == nil || *output.NextToken == "" {
			break
		}
		identityNextToken = output.NextToken
	}

	metrics.CognitoIdentityPools.WithLabelValues(account, accountName, region).Set(identityCount)

	log.Debug().
		Str("region", region).
		Int("user_pool_statuses", len(statusCounts)).
		Str("identity_pools", strconv.FormatFloat(identityCount, 'f', 0, 64)).
		Msg("Cognito collection completed")

	return nil
}
