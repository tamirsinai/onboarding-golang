package awsLocal

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var localCfg aws.Config
var SqsClient *sqs.Client
var S3Client *s3.Client
var QueueURL = "https://sqs.us-east-1.amazonaws.com/586929748635/onboarding-repo-scan"
var BucketName = "onboarding-repo-scans-results"

func LoadLocalConfig(ctx context.Context) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return
	}
	localCfg = cfg
}

func LoadSqsConfig() {
	SqsClient = sqs.NewFromConfig(localCfg)
}

func LoadS3Config() {
	S3Client = s3.NewFromConfig(localCfg)
}
