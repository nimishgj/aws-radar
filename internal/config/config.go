package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Collection   CollectionConfig   `mapstructure:"collection"`
	AWS          AWSConfig          `mapstructure:"aws"`
	CostExplorer CostExplorerConfig `mapstructure:"cost_explorer"`
	CostCUR      CostCURConfig      `mapstructure:"cost_cur"`
	Collectors   []string           `mapstructure:"collectors"`
	Logging      LoggingConfig      `mapstructure:"logging"`
}

type ServerConfig struct {
	Port        int    `mapstructure:"port"`
	MetricsPath string `mapstructure:"metrics_path"`
	HealthPath  string `mapstructure:"health_path"`
}

type CollectionConfig struct {
	Interval time.Duration `mapstructure:"interval"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type AWSConfig struct {
	Regions []string `mapstructure:"regions"`
}

type CostExplorerConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Frequency string `mapstructure:"frequency"`
}

type CostCURConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	Bucket       string `mapstructure:"bucket"`
	Prefix       string `mapstructure:"prefix"`
	ReportName   string `mapstructure:"report_name"`
	Format       string `mapstructure:"format"`
	Frequency    string `mapstructure:"frequency"`
	Region       string `mapstructure:"region"`
	MaxResources int    `mapstructure:"max_resources"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app")

	// Check for custom config path
	if configPath := os.Getenv("AWS_RADAR_CONFIG"); configPath != "" {
		viper.SetConfigFile(configPath)
	}

	// Set defaults
	viper.SetDefault("server.port", 9090)
	viper.SetDefault("server.metrics_path", "/metrics")
	viper.SetDefault("server.health_path", "/health")
	viper.SetDefault("collection.interval", "60s")
	viper.SetDefault("collection.timeout", "30s")
	viper.SetDefault("aws.regions", []string{"us-east-1"})
	viper.SetDefault("cost_explorer.enabled", false)
	viper.SetDefault("cost_explorer.frequency", "daily")
	viper.SetDefault("cost_cur.enabled", false)
	viper.SetDefault("cost_cur.frequency", "daily")
	viper.SetDefault("cost_cur.region", "us-east-1")
	viper.SetDefault("cost_cur.max_resources", 100)
	viper.SetDefault("collectors", []string{
		"apigateway",
		"apigatewayv2",
		"autoscaling",
		"athena",
		"configservice",
		"mq",
		"ses",
		"cloudformation",
		"documentdb",
		"neptune",
		"memorydb",
		"timestream",
		"fsx",
		"backup",
		"kinesis",
		"firehose",
		"kinesisanalytics",
		"emr",
		"elasticbeanstalk",
		"kms",
		"cloudtrail",
		"kinesisvideo",
		"opensearchserverless",
		"s3control",
		"batch",
		"codebuild",
		"codepipeline",
		"codedeploy",
		"ecr",
		"ec2",
		"efs",
		"eventbridge",
		"glue",
		"apprunner",
		"transfer",
		"msk",
		"redshift",
		"s3",
		"rds",
		"lambda",
		"ecs",
		"elb",
		"eks",
		"dynamodb",
		"elasticache",
		"opensearch",
		"guardduty",
		"securityhub",
		"inspector2",
		"macie",
		"waf",
		"secretsmanager",
		"sfn",
		"ssm",
		"sqs",
		"sns",
		"ebs",
		"vpc",
		"acm",
		"cloudfront",
		"route53",
		"route53resolver",
		"iam",
		"shield",
		"ecrpublic",
		"directconnect",
		"organizations",
		"controltower",
		"bedrock",
		"sagemaker",
		"quicksight",
		"workspaces",
		"appstream",
		"connect",
		"amplify",
		"globalaccelerator",
		"datasync",
		"dms",
	})
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
