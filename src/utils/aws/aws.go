package awsutil

import (
	"context"
	"os"
	"y-ogura/echo_lambda_sample/src/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type AWSService interface {
}

type awsService struct {
	Config aws.Config
}

func NewAWSService() *awsService {
	var cfg aws.Config
	var err error
	region := os.Getenv("REGION")
	ctx := context.Background()
	if os.Getenv("LOCALSTACK_ENDPOINT") != "" {
		// LocalStack用
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				os.Getenv("LOCALSTACK_KEY"), os.Getenv("LOCALSTACK_SECRET"), "")),
		)
		if err != nil {
			logger.Panic("Failed to load default config for LocalStack", err)
		}
	} else {
		// Dev,Stg,Prod用
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
		)
		if err != nil {
			logger.Panic("Failed to load default config for AWS Connection", err)
		}
	}
	return &awsService{
		Config: cfg,
	}
}
