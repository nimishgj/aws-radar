package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Collection CollectionConfig `mapstructure:"collection"`
	AWS        AWSConfig        `mapstructure:"aws"`
	Collectors []string         `mapstructure:"collectors"`
	Logging    LoggingConfig    `mapstructure:"logging"`
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
	viper.SetDefault("collectors", []string{
		"apigateway",
		"apigatewayv2",
		"autoscaling",
		"athena",
		"ecr",
		"ec2",
		"efs",
		"eventbridge",
		"glue",
		"s3",
		"rds",
		"lambda",
		"ecs",
		"elb",
		"eks",
		"dynamodb",
		"elasticache",
		"opensearch",
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
		"iam",
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
