package core

// Can be injected into tests to provide mocked config values
var MockConfig = &Config{
	Season: 2022,
	Debug: true,
	API: configAPI{
		Host: "api.sample-host.com",
		UrlFormat: "http://%s/v1",
		Key: "not-a-real-key",
	},
	Database: configDatabase{
		User: "test",
		Password: "fake",
		Location: "localhost",
		Name: "nice-fixture",
	},
	AWS: configAWS{
		Region: "us-test-1",
		AccessKeyId: "FAKE-ACCESS-KEY",
		SecretAccessKey: "fakeSecretKey",
		BucketName: "test123",
	},
}