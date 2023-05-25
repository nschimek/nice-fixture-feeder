package request

import (
	"io"
	"net/http"
	"net/url"

	"github.com/nschimek/nice-fixture-feeder/core"
)

type Requester struct {
	config *core.Config
	client *http.Client
}


const (
	headerKey = "X-RapidAPI-Key"
	headerHost = "X-RapidAPI-Host"
	baseUrl = "https://api-football-v1.p.rapidapi.com/v3"
)

func NewRequester(config *core.Config) *Requester {
	return &Requester{
		config: config,
		client: http.DefaultClient,
	}
}

func (r *Requester) Get(endpoint string, params url.Values) ([]byte, error) {
	req, err := http.NewRequest("GET", baseUrl + "/" + endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add(headerKey, r.config.Api.Key)
	req.Header.Add(headerHost, r.config.Api.Host)
	
	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	return r.doRequest(req)
}

func (r *Requester) doRequest(req *http.Request) ([]byte, error) {
	core.Log.Infof("Requesting %s...", req.URL.String())

	res, err := r.client.Do(req)
	
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}