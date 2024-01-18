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

var OutputFileName string = "output.json"

func WriteOutputFile(scan *models.Scan) error {
	jsonData, err := json.Marshal(&scan)
	if err != nil {
		return err
	}
	if err := os.WriteFile(OutputFileName, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

func Send() {
	file, err := os.Open(OutputFileName)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer file.Close()

	_, err = awsLocal.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &awsLocal.BucketName,
		Key:    &OutputFileName,
		Body:   file,
	})
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}
	fmt.Println("File uploaded successfully.")
}
