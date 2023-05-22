package core

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	awsHttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
)

var S3 *AwsS3

type AwsS3 struct {
	client *s3.Client
}

func SetupS3(config *Config) {
	S3 = &AwsS3{
		client: s3.NewFromConfig(getAwsConfig(config)),
	}
}

func getAwsConfig(cfg *Config) aws.Config {
	Log.WithField("region", cfg.AWS.Region).Info("Setting up AWS Connection...")
	creds := credentials.NewStaticCredentialsProvider(cfg.AWS.AccessKeyId, cfg.AWS.SecretAccessKey, "")
	
	awsConfig, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(cfg.AWS.Region),
	)

	if err != nil {
		Log.Fatal(err)
	}

	return awsConfig
} 

func (as *AwsS3) Exists(bucket, key string) (bool, error) {
	_, err := as.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(key),
	})
	if err != nil {
		var respErr *awsHttp.ResponseError
		if errors.As(err, &respErr) && respErr.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (as *AwsS3) Upload(data []byte, bucket, key string) error {
	uploader := manager.NewUploader(as.client)
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key: aws.String(key),
			Body: strings.NewReader(string(data)),
	})
	return err
}