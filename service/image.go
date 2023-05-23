package service

import (
	"errors"
	"io"
	"net/http"
	"path"

	"github.com/nschimek/nice-fixture-feeder/core"
)

type ImageService struct {
	s3 *core.AwsS3
}

func NewImageService(s3 *core.AwsS3) *ImageService {
	return &ImageService{
		s3: s3,
	}
}

func (is *ImageService) TransferURL(url, bucket string) bool {
	core.Log.Infof("Transferring image %s to S3 bucket %s...", url, bucket)
	fileName := path.Base(url)

	if e, err := is.s3.Exists(bucket, fileName); !e {
		bytes, err := is.download(url)
		if err != nil {
			core.Log.Errorf("Issue while downloading image: %v", err)
			return false
		}

		err = is.s3.Upload(bytes, bucket, fileName)
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

func (is *ImageService) download(url string) ([]byte, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	} else if (res.StatusCode != 200) {
		return nil, errors.New("received non-200 response code")
	}

	body, err := io.ReadAll(res.Body)

	return body, err
}