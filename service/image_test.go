package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/stretchr/testify/suite"
)

const (
	testData = "dummy test data"
)

type ImageServiceTestSuite struct {
	suite.Suite
	mockHttp *httptest.Server
	mockS3 *core.MockS3Client
}

func (suite *ImageServiceTestSuite) SetupTest() {
	suite.mockHttp = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testData)
	}))
	defer suite.mockHttp.Close()
	suite.mockS3 = &core.MockS3Client{}
}

func (suite *ImageServiceTestSuite) TestSuccessful() {
	url, bucket, keyFormat := "http://wwww.test.com/test.txt", "test-bucket", "/test/%s"

	suite.mockS3.On("Exists", bucket, "/test/test-bucket/test.txt").Return(false, nil)
	suite.mockS3.On("Upload", []byte(testData), bucket, "/test/test-bucket/test.txt").Return(nil)

	is := &imageService{s3: suite.mockS3}

	r := is.TransferURL(url, bucket, keyFormat)

	suite.True(r)
}
