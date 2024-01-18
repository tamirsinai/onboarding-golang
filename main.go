package main

import (
	awsLocal "github.com/tamirsinai/onboarding-golang/pkg/awslocal"
	"github.com/tamirsinai/onboarding-golang/pkg/env"
	"github.com/tamirsinai/onboarding-golang/pkg/input"
	"github.com/tamirsinai/onboarding-golang/pkg/logger"
)

func main() {
	env.Load()
	logger.Init()
	awsLocal.LoadLocalConfig()
	awsLocal.LoadSqsConfig()
	awsLocal.LoadS3Config()
	input.Receive()
}
