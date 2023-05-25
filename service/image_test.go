package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/stretchr/testify/suite"
)

const (
	testFileName = "test.txt"
	testBucket = "test-bucket"
	testKeyFormat = "/test/%s"
	testData = "dummy test data"
	testFinalKeyName = "/test/test.txt"
)

type ImageServiceTestSuite struct {
	suite.Suite
	mockS3 *core.MockS3Client
	is *imageService
}

func TestLeagueRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ImageServiceTestSuite))
}

func (suite *ImageServiceTestSuite) SetupTest() {
	suite.mockS3 = &core.MockS3Client{}
	suite.is = &imageService{s3: suite.mockS3}
}

func (suite *ImageServiceTestSuite) TestSuccessful() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testData)
	}))
	defer svr.Close()

	suite.mockS3.On("Exists", testBucket, testFinalKeyName).Return(false, nil)
	suite.mockS3.On("Upload", []byte(testData), testBucket, testFinalKeyName).Return(nil)

	r := suite.is.TransferURL(svr.URL + "/" + testFileName, testBucket, testKeyFormat)

	suite.True(r)
}

func (suite *ImageServiceTestSuite) TestExistsError() {
	suite.mockS3.On("Exists", testBucket, testFinalKeyName).Return(false, errors.New("test"))

	// url does not matter for this test, as we never get there
	r := suite.is.TransferURL(testFileName, testBucket, testKeyFormat)

	suite.False(r)
}

func (suite *ImageServiceTestSuite) TestDownloadError() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Test Error", http.StatusBadRequest)
	}))
	defer svr.Close()

	suite.mockS3.On("Exists", testBucket, testFinalKeyName).Return(false, nil)

	r := suite.is.TransferURL(svr.URL + "/" + testFileName, testBucket, testKeyFormat)

	suite.False(r)
}

func (suite *ImageServiceTestSuite) TestUploadError() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testData)
	}))
	defer svr.Close()

	suite.mockS3.On("Exists", testBucket, testFinalKeyName).Return(false, nil)
	suite.mockS3.On("Upload", []byte(testData), testBucket, testFinalKeyName).Return(errors.New("test"))

	r := suite.is.TransferURL(svr.URL + "/" + testFileName, testBucket, testKeyFormat)
	suite.False(r)
}

func (suite *ImageServiceTestSuite) TestImageExists() {
	suite.mockS3.On("Exists", testBucket, testFinalKeyName).Return(true, nil)

	// url does not matter for this test, as we never get there
	r := suite.is.TransferURL(testFileName, testBucket, testKeyFormat)

	suite.False(r)
}
