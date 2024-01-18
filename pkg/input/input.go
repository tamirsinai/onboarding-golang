package input

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

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

func Receive() {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:              &awsLocal.QueueURL,
		MaxNumberOfMessages:   1,
		WaitTimeSeconds:       20,
		MessageAttributeNames: []string{"All"},
	}

	for {
		resp, err := awsLocal.SqsClient.ReceiveMessage(context.TODO(), input)
		if err != nil {
			fmt.Println("Error receiving messages:", err)
			return
		}

		for _, message := range resp.Messages {
			var bodyMessage models.Input
			if err := json.Unmarshal([]byte(*message.Body), &bodyMessage); err != nil {
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

			path, err := filepath.Abs(output.OutputFileName)
			if err != nil {
				logger.Logger.Error("Error getting absolute path:", zap.Error(err))
				return
			}
			logger.Logger.Info(path)

			if err := repo.DeleteClonedProjectsDir(); err != nil {
				logger.Logger.Error("Error with delete repo:", zap.Error(err))
				return
			}

			output.Send()
			delete(awsLocal.SqsClient, awsLocal.QueueURL, *message.ReceiptHandle)
		}

		time.Sleep(3 * time.Second)
	}
}

func delete(client *sqs.Client, queueURL, receiptHandle string) {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: &receiptHandle,
	}

	_, err := client.DeleteMessage(context.TODO(), input)
	if err != nil {
		fmt.Println("Error deleting message:", err)
		return
	}

	fmt.Println("Message deleted.")
}
