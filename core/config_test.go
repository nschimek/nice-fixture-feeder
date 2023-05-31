package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// this is more of an integration test, but it does the job
func TestLoadConfigFile(t *testing.T) {
	SetupViper()
	SetupConfigFile("../config/sample.yaml")

	assert.Equal(t, 2022, Cfg.Season)
	assert.Equal(t, true, Cfg.Debug)
	assert.Equal(t, "fake-sample.api.com", Cfg.API.Host)
	assert.Equal(t, "https://%s/v1", Cfg.API.UrlFormat)
	assert.Equal(t, "fake-api-key", Cfg.API.Key)
	assert.Equal(t, "fake-user", Cfg.Database.User)
	assert.Equal(t, "fake-password", Cfg.Database.Password)
	assert.Equal(t, "localhost", Cfg.Database.Location)
	assert.Equal(t, 3306, Cfg.Database.Port)
	assert.Equal(t, "us-east-1", Cfg.AWS.Region)
	assert.Equal(t, "fake-access-key-id", Cfg.AWS.AccessKeyId)
	assert.Equal(t, "fake-secret-access-key", Cfg.AWS.SecretAccessKey)
	assert.Equal(t, "fake-bucket-name", Cfg.AWS.BucketName)
}

func TestWithoutConfigFile(t *testing.T) {
	SetupViper()

	os.Setenv("NF_USE_CONFIG_FILE", "false")
	os.Setenv("NF_SEASON", "2022")
	os.Setenv("NF_API_HOST", "fake-sample.api.com")
	os.Setenv("NF_API_URL_FORMAT", "https://%s/v1")
	os.Setenv("NF_API_KEY", "fake-api-key")
	os.Setenv("NF_DATABASE_USER", "fake-user")
	os.Setenv("NF_DATABASE_PASSWORD", "fake-password")
	os.Setenv("NF_DATABASE_LOCATION", "localhost")
	os.Setenv("NF_DATABASE_PORT", "3306")
	os.Setenv("NF_AWS_REGION", "us-east-1")
	os.Setenv("NF_AWS_ACCESS_KEY_ID", "fake-access-key-id")
	os.Setenv("NF_AWS_SECRET_ACCESS_KEY", "fake-secret-access-key")
	os.Setenv("NF_AWS_BUCKET_NAME", "fake-bucket-name")

	SetupConfigFile("")
	
	assert.Equal(t, false, Cfg.Debug)
	assert.Equal(t, "fake-sample.api.com", Cfg.API.Host)
	assert.Equal(t, "https://%s/v1", Cfg.API.UrlFormat)
	assert.Equal(t, "fake-api-key", Cfg.API.Key)
	assert.Equal(t, "fake-user", Cfg.Database.User)
	assert.Equal(t, "fake-password", Cfg.Database.Password)
	assert.Equal(t, "localhost", Cfg.Database.Location)
	assert.Equal(t, "us-east-1", Cfg.AWS.Region)
	assert.Equal(t, "fake-access-key-id", Cfg.AWS.AccessKeyId)
	assert.Equal(t, "fake-secret-access-key", Cfg.AWS.SecretAccessKey)
	assert.Equal(t, "fake-bucket-name", Cfg.AWS.BucketName)
}