package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsLocal "github.com/tamirsinai/onboarding-golang/pkg/awslocal"
	"github.com/tamirsinai/onboarding-golang/pkg/env"
	"github.com/tamirsinai/onboarding-golang/pkg/input"
	"github.com/tamirsinai/onboarding-golang/pkg/logger"
)

func main() {
	lambda.Start(startLambda)
}

func startLambda(ctx context.Context, event events.SQSEvent) {
	env.Load()
	logger.Init()
	awsLocal.LoadLocalConfig(ctx)
	awsLocal.LoadSqsConfig()
	awsLocal.LoadS3Config()
	input.Receive(ctx, event)
}
