package request

import (
	"fmt"
	"io"
	"net/http"

	"github.com/nschimek/nice-fixture-feeder/core"
)

const (
	leaguesEndpoint = "leagues"
)

type LeagueRequest struct {
	Config *core.Config
	request *http.Request
	client *http.Client
}

func NewLeagueRequest(config *core.Config) *LeagueRequest {
	return &LeagueRequest{
		Config: config,
		request: NewHttpRequest(config, leaguesEndpoint),
		client: http.DefaultClient,
	}
}

func (r *LeagueRequest) Request() {
	res, err := r.client.Do(r.request)
	core.IfErrorFatal(err)

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	core.IfErrorFatal(err)

	fmt.Println(string(body))
}