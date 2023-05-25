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

var S3 S3Client

func SetupS3(config *Config) {
	S3 = &s3Client{
		client: s3.NewFromConfig(getAwsConfig(config)),
	}
}

// S3Client interface represents a DAL with the S3 Client library
type S3Client interface {
	Exists(bucket, key string) (bool, error)
	Upload(data []byte, bucket, key string) error
}

type s3Client struct {
	client *s3.Client
}

func (s *s3Client) Exists(bucket, key string) (bool, error) {
	_, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(key),
	})
	return s3NotFoundErrorHandler(err)
}

func (s *s3Client) Upload(data []byte, bucket, key string) error {
	uploader := manager.NewUploader(s.client)
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(key),
		Body: strings.NewReader(string(data)),
	})
	return err
}

// checks if we have an actual response error or if the key simply is not found
func s3NotFoundErrorHandler(err error) (bool, error) {
	if err != nil {
		var respErr *awsHttp.ResponseError
		if errors.As(err, &respErr) && respErr.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func getAwsConfig(cfg *Config) aws.Config {
	Log.WithField("region", cfg.AWS.Region).Info("Setting up AWS Config...")
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

