package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nschimek/nice-fixture-feeder/core"
)

//go:generate mockery --name Requester
type Requester[T any] interface {
	Get(endpoint string, params url.Values) (*Response[T], error)
}

type requester[T any] struct {
	BaseUrl string
	config *core.Config
}

const (
	headerKey = "X-RapidAPI-Key"
	headerHost = "X-RapidAPI-Host"
)

func NewRequester[T any](config *core.Config) *requester[T] {
	return &requester[T]{
		BaseUrl: fmt.Sprintf(config.API.UrlFormat, config.API.Host),
		config: config,
	}
}

func (r *requester[T]) Get(endpoint string, params url.Values) (*Response[T], error) {
	req, _ := http.NewRequest("GET", r.BaseUrl + "/" + endpoint, nil)

	req.Header.Add(headerKey, r.config.API.Key)
	req.Header.Add(headerHost, r.config.API.Host)
	
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	bytes, err := r.doRequest(req)

	if err != nil {
		return nil, err
	}

	return r.unmarshal(bytes)
}

func (r *requester[T]) doRequest(req *http.Request) ([]byte, error) {
	core.Log.Infof("Requesting %s...", req.URL.String())

	c := http.DefaultClient
	res, err := c.Do(req)
	
	if err != nil {
		return nil, err
	} else if (res.StatusCode != http.StatusOK) {
		return nil, errors.New("received non-200 response code")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *requester[T]) unmarshal(bytes []byte) (*Response[T], error) {
	var response Response[T]
	err := json.Unmarshal(bytes, &response)
	
	if len(response.Response) == 0 {
		core.Log.Warn("got 0 entities from API response")
	}

	if err == nil {
		return &response, nil
	} else {
		return nil, err
	}
}