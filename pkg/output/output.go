package output

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tamirsinai/onboarding-golang/models"
	awsLocal "github.com/tamirsinai/onboarding-golang/pkg/awslocal"
)

var outputFileName string = "output.json"
var OutputFilePath string = "/tmp/" + outputFileName

func WriteOutputFile(scan *models.Scan) error {
	jsonData, err := json.Marshal(&scan)
	if err != nil {
		return err
	}
	if err := os.WriteFile(OutputFilePath, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

func Send(ctx context.Context) {
	file, err := os.Open(OutputFilePath)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer file.Close()

	_, err = awsLocal.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &awsLocal.BucketName,
		Key:    &outputFileName,
		Body:   file,
	})
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}
	fmt.Println("File uploaded successfully.")
}
