package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/sirupsen/logrus"
)

//go:generate mockery --name Image --filename image_mock.go
type Image interface {
	TransferURL(url, bucket, keyFormat string) bool
}

type image struct {
	s3 core.S3Client
}

func NewImage(s3 core.S3Client) *image {
	return &image{
		s3: s3,
	}
}

func (is *image) TransferURL(url, bucket, keyFormat string) bool {
	time.Sleep(500 * time.Millisecond) // half second rate limit 
	finalKeyName := fmt.Sprintf(keyFormat, path.Base(url))
	core.Log.WithFields(logrus.Fields{
		"url": url,
		"bucket": bucket,
		"key": finalKeyName,
	}).Infof("Transferring image to S3...")

	if e, err := is.s3.Exists(bucket, finalKeyName); !e {
		bytes, err := is.download(url)
		if err != nil {
			core.Log.Errorf("Issue while downloading image: %v", err)
			return false
		}

		err = is.s3.Upload(bytes, bucket, finalKeyName)
		if err != nil {
			core.Log.Errorf("Issue while uploading image: %v", err)
			return false
		}
		
		return true
	} else if err != nil {
		core.Log.Errorf("Issue while determining if image exists: %v", err)
		return false
	} else {
		core.Log.Infof("Image already exists, skipping download!")
		return false
	}
}

func (is *image) download(url string) ([]byte, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	} else if (res.StatusCode != http.StatusOK) {
		return nil, errors.New("received non-200 response code")
	}
	
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	return body, err
}