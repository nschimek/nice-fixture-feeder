package request

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nschimek/nice-fixture-feeder/core"
)

const (
	leaguesEndpoint = "leagues"
)

type LeagueRequest struct {
	config *core.Config
	requester *Requester
	client *http.Client
}

func NewLeagueRequest(config *core.Config) *LeagueRequest {
	return &LeagueRequest{
		config: config,
		requester: NewRequester(config),
		client: http.DefaultClient,
	}
}

func (r *LeagueRequest) Request(id, season int) {
	p := url.Values{}
	p.Add("id", strconv.Itoa(id))
	p.Add("season", strconv.Itoa(season))

	fmt.Println(string(r.requester.Get(leaguesEndpoint, p)))
}