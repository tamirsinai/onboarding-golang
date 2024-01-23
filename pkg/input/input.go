package input

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/tamirsinai/onboarding-golang/models"
	awsLocal "github.com/tamirsinai/onboarding-golang/pkg/awslocal"
	"github.com/tamirsinai/onboarding-golang/pkg/logger"
	"github.com/tamirsinai/onboarding-golang/pkg/output"
	"github.com/tamirsinai/onboarding-golang/pkg/repo"
	"github.com/tamirsinai/onboarding-golang/pkg/scan"
	"go.uber.org/zap"
)

const inputFileName string = "input.json"

func ReadFile() (*models.Input, error) {
	jsonData, err := os.ReadFile(inputFileName)
	if err != nil {
		return nil, err
	}

	var input models.Input
	if err := json.Unmarshal(jsonData, &input); err != nil {
		return nil, err
	}

	return &input, nil
}

func Receive(ctx context.Context, event events.SQSEvent) {
	if err := repo.DeleteClonedProjectsDir(); err != nil {
		logger.Logger.Error("Error with delete repo:", zap.Error(err))
		return
	}
	for _, record := range event.Records {
		var bodyMessage models.Input
		if err := json.Unmarshal([]byte(record.Body), &bodyMessage); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}

		if err := repo.CloneRepositoryToScan(bodyMessage.CloneUrl); err != nil {
			logger.Logger.Error("Error clone repo:", zap.Error(err))
			return
		}

		scan, err := scan.ScanRepoFiles(repo.ClonedProjectsDir, bodyMessage.Size)
		if err != nil {
			logger.Logger.Error("Error scanning repo files:", zap.Error(err))
			return
		}

		if err := output.WriteOutputFile(scan); err != nil {
			logger.Logger.Error("Error write output file:", zap.Error(err))
			return
		}

		path, err := filepath.Abs(output.OutputFilePath)
		if err != nil {
			logger.Logger.Error("Error getting absolute path:", zap.Error(err))
			return
		}
		logger.Logger.Info(path)

		if err := repo.DeleteClonedProjectsDir(); err != nil {
			logger.Logger.Error("Error with delete repo:", zap.Error(err))
			return
		}

		output.Send(ctx)
		delete(ctx, awsLocal.SqsClient, record.ReceiptHandle)
	}
}

func delete(ctx context.Context, client *sqs.Client, receiptHandle string) {
	input := &sqs.DeleteMessageInput{
		ReceiptHandle: &receiptHandle,
		QueueUrl:      &awsLocal.QueueURL,
	}

	_, err := client.DeleteMessage(ctx, input)
	if err != nil {
		fmt.Println("Error deleting message:", err)
		return
	}

	fmt.Println("Message deleted.")
}
