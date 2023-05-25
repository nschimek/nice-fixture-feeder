package core

import (
	"errors"
	"testing"

	goHttp "net/http"

	awsHttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/smithy-go/transport/http"

	"github.com/stretchr/testify/assert"
)

func TestS3ErrorHandlerNotFound(t *testing.T) {
	test := &awsHttp.ResponseError{
		ResponseError: &http.ResponseError{
			Response: &http.Response{
				Response: &goHttp.Response{
					StatusCode: goHttp.StatusNotFound,
				},
			},
		},
	}

	r, err := s3NotFoundErrorHandler(test)
	assert.False(t, r)
	assert.Nil(t, err)
}

func TestS3ErrorHandlerOtherError(t *testing.T) {
	r, err := s3NotFoundErrorHandler(errors.New("test"))
	assert.False(t, r)
	assert.Error(t, err)
}

func TestS3ErrorHandlerNoError(t *testing.T) {
	r, err := s3NotFoundErrorHandler(nil)
	assert.True(t, r)
	assert.Nil(t, err)
}