package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type ElasticBeanstalkCollector struct{}

func NewElasticBeanstalkCollector() *ElasticBeanstalkCollector {
	return &ElasticBeanstalkCollector{}
}

func (c *ElasticBeanstalkCollector) Name() string {
	return "elasticbeanstalk"
}

func (c *ElasticBeanstalkCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := elasticbeanstalk.NewFromConfig(cfg)
	output, err := client.DescribeApplications(ctx, &elasticbeanstalk.DescribeApplicationsInput{})
	if err != nil {
		return err
	}

	count := float64(len(output.Applications))
	metrics.ElasticBeanstalkApplications.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
