package collector

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	smithy "github.com/aws/smithy-go"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type ECRCollector struct{}

func NewECRCollector() *ECRCollector {
	return &ECRCollector{}
}

func (c *ECRCollector) Name() string {
	return "ecr"
}

func (c *ECRCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := ecr.NewFromConfig(cfg)

	var count float64
	policyCounts := map[string]float64{
		"public":  0,
		"private": 0,
		"error":   0,
	}

	paginator := ecr.NewDescribeRepositoriesPaginator(client, &ecr.DescribeRepositoriesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Repositories))

		for _, repo := range page.Repositories {
			status := getRepositoryPolicyStatus(ctx, client, repo.RepositoryName)
			policyCounts[status]++
		}
	}

	metrics.ECRRepositories.WithLabelValues(account, accountName, region).Set(count)
	for status, n := range policyCounts {
		metrics.ECRRepositoriesByPolicy.WithLabelValues(account, accountName, region, status).Set(n)
	}

	log.Debug().
		Str("region", region).
		Float64("repository_count", count).
		Float64("public", policyCounts["public"]).
		Float64("private", policyCounts["private"]).
		Float64("error", policyCounts["error"]).
		Msg("ECR repository collection completed")

	return nil
}

// getRepositoryPolicyStatus returns "public", "private", or "error".
// Repos without any attached policy are "private".
func getRepositoryPolicyStatus(ctx context.Context, client *ecr.Client, repoName *string) string {
	out, err := client.GetRepositoryPolicy(ctx, &ecr.GetRepositoryPolicyInput{
		RepositoryName: repoName,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "RepositoryPolicyNotFoundException" {
			return "private"
		}
		log.Warn().
			Err(err).
			Str("repository", aws.ToString(repoName)).
			Msg("Failed to get ECR repository policy")
		return "error"
	}
	if isPublicECRPolicy(aws.ToString(out.PolicyText)) {
		return "public"
	}
	return "private"
}

// isPublicECRPolicy reports whether the policy JSON has any Allow statement
// granting access to the wildcard principal "*" (i.e. public access).
func isPublicECRPolicy(policyJSON string) bool {
	if policyJSON == "" {
		return false
	}
	var doc struct {
		Statement []struct {
			Effect    string      `json:"Effect"`
			Principal interface{} `json:"Principal"`
		} `json:"Statement"`
	}
	if err := json.Unmarshal([]byte(policyJSON), &doc); err != nil {
		return false
	}
	for _, stmt := range doc.Statement {
		if stmt.Effect != "Allow" {
			continue
		}
		if principalIsWildcard(stmt.Principal) {
			return true
		}
	}
	return false
}

// principalIsWildcard checks for the standard "public" indicators:
// Principal: "*"  OR  Principal: { "AWS": "*" }  OR  Principal: { "AWS": ["*", ...] }
func principalIsWildcard(p interface{}) bool {
	switch v := p.(type) {
	case string:
		return v == "*"
	case map[string]interface{}:
		for _, val := range v {
			switch vv := val.(type) {
			case string:
				if vv == "*" {
					return true
				}
			case []interface{}:
				for _, item := range vv {
					if s, ok := item.(string); ok && s == "*" {
						return true
					}
				}
			}
		}
	}
	return false
}
