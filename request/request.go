package request

import (
	"net/http"

	"github.com/nschimek/nice-fixture-feeder/core"
)

const (
	headerKey = "X-RapidAPI-Key"
	headerHost = "X-RapidAPI-Host"
	baseUrl = "https://api-football-v1.p.rapidapi.com/v3"
)

func NewHttpRequest(config *core.Config, endpoint string) *http.Request {
	req, _ := http.NewRequest("GET", baseUrl + "/" + endpoint, nil)
	req.Header.Add(headerKey, config.Api.Key)
	req.Header.Add(headerHost, config.Api.Host)
	return req
}