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

type imageServiceTestSuite struct {
	suite.Suite
	mockS3 *core.MockS3Client
	imageService *imageService
}

func TestLeagueRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(imageServiceTestSuite))
}

func (s *imageServiceTestSuite) SetupTest() {
	s.mockS3 = &core.MockS3Client{}
	s.imageService = &imageService{s3: s.mockS3}
}

func (s *imageServiceTestSuite) TestSuccessful() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testData)
	}))
	defer svr.Close()

	s.mockS3.EXPECT().Exists(testBucket, testFinalKeyName).Return(false, nil)
	s.mockS3.EXPECT().Upload([]byte(testData), testBucket, testFinalKeyName).Return(nil)

	r := s.imageService.TransferURL(svr.URL + "/" + testFileName, testBucket, testKeyFormat)

	s.True(r)
}

func (s *imageServiceTestSuite) TestExistsError() {
	s.mockS3.EXPECT().Exists(testBucket, testFinalKeyName).Return(false, errors.New("test"))

	// url does not matter for this test, as we never get there
	r := s.imageService.TransferURL(testFileName, testBucket, testKeyFormat)

	s.False(r)
}

func (s *imageServiceTestSuite) TestDownloadResponseError() {
	s.mockS3.EXPECT().Exists(testBucket, testFinalKeyName).Return(false, nil)

	r, err := s.imageService.download("invalid URL")

	s.Nil(r)
	s.ErrorContains(err, "unsupported protocol scheme")
}

func (s *imageServiceTestSuite) TestDownloadNon200() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "test", http.StatusNoContent)
	}))
	defer svr.Close()

	r, err := s.imageService.download(svr.URL + "/" + testFileName)

	s.Nil(r)
	s.ErrorContains(err, "received non-200 response code")
}

func (s *imageServiceTestSuite) TestReadError() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
	}))
	defer svr.Close()

	r, err := s.imageService.download(svr.URL + "/" + testFileName)

	s.Equal(r, []byte{})
	s.ErrorContains(err, "unexpected EOF")
}

func (s *imageServiceTestSuite) TestUploadError() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testData)
	}))
	defer svr.Close()

	s.mockS3.EXPECT().Exists(testBucket, testFinalKeyName).Return(false, nil)
	s.mockS3.EXPECT().Upload([]byte(testData), testBucket, testFinalKeyName).Return(errors.New("test"))

	r := s.imageService.TransferURL(svr.URL + "/" + testFileName, testBucket, testKeyFormat)
	s.False(r)
}

func (s *imageServiceTestSuite) TestImageExists() {
	s.mockS3.EXPECT().Exists(testBucket, testFinalKeyName).Return(true, nil)

	// url does not matter for this test, as we never get there
	r := s.imageService.TransferURL(testFileName, testBucket, testKeyFormat)

	s.False(r)
}
