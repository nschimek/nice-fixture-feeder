package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
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

	var response Response[model.League]
	json.Unmarshal(r.requester.Get(leaguesEndpoint, p), &response)

	fmt.Printf("%+v\n", response)
}