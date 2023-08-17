package mocks

import "github.com/nschimek/nice-fixture-feeder/core"

// Can be injected into tests to provide mocked config values
var Config = core.Config{
	Season: 2022,
	Leagues: []int{39, 140},
	Debug: true,
	API: core.ConfigAPI{
		Host: "api.sample-host.com",
		UrlFormat: "http://%s/v1",
		Key: "not-a-real-key",
	},
	Database: core.ConfigDatabase{
		User: "lol",
		Password: "fake",
		Location: "localhost",
		Name: "nice-fixture",
	},
	AWS: core.ConfigAWS{
		Region: "us-test-1",
		AccessKeyId: "FAKE-ACCESS-KEY",
		SecretAccessKey: "fakeSecretKey",
		BucketName: "test123",
	},
}