package request

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
)

const (
	leaguesEndpoint = "leagues"
)

type LeagueRequest struct {
	config *core.Config
	requester *Requester
	client *http.Client
	repo *repository.LeagueRepository
	RequestedData []model.League
}

func NewLeagueRequest(config *core.Config, repo *repository.LeagueRepository) *LeagueRequest {
	return &LeagueRequest{
		config: config,
		requester: NewRequester(config),
		client: http.DefaultClient,
		repo: repo,
	}
}

func (r *LeagueRequest) Request(season, id int) {
	p := url.Values{}
	p.Add("id", strconv.Itoa(id))
	p.Add("season", strconv.Itoa(season))

	var response Response[model.League]
	json.Unmarshal(r.requester.Get(leaguesEndpoint, p), &response)

	r.RequestedData = append(r.RequestedData, response.Response...)
}

func (r *LeagueRequest) Persist() {
	r.repo.UpsertLeagues(r.RequestedData)
}