package request

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/service"
)

const (
	leaguesEndpoint = "leagues"
	leagueKeyFormat = "images/logos/leagues/%s"
	countryKeyFormat = "images/flags/%s"
)

type LeagueRequest struct {
	config *core.Config
	requester *Requester
	client *http.Client
	repo *repository.LeagueRepository
	imageService *service.ImageService
	RequestedData []model.League
}

func NewLeagueRequest(config *core.Config, repo *repository.LeagueRepository, is *service.ImageService) *LeagueRequest {
	return &LeagueRequest{
		config: config,
		requester: NewRequester(config),
		client: http.DefaultClient,
		imageService: is,
		repo: repo,
	}
}

func (r *LeagueRequest) Request(season, id int) {
	p := url.Values{}
	p.Add("id", strconv.Itoa(id))
	p.Add("season", strconv.Itoa(season))

	var response Response[model.League]
	data, err := r.requester.Get(leaguesEndpoint, p)
	
	if err != nil {
		core.Log.Error("Could not get league %d: %v", id, err)
	}

	err = json.Unmarshal(data, &response)

	if err == nil {
		r.RequestedData = append(r.RequestedData, response.Response...)
	} else {
		core.Log.Errorf("Could not unmarshal league %d JSON: %v", id, err)
	}
}

func (r *LeagueRequest) Persist() {
	rs := r.repo.UpsertLeagues(r.RequestedData)
	rs.LogErrors()
	rs.LogSuccesses()
}

func (r *LeagueRequest) PostPersist() {
	for _, league := range r.RequestedData {
		r.imageService.TransferURL(league.League.Logo, r.config.AWS.BucketName, leagueKeyFormat)
		r.imageService.TransferURL(league.Country.Flag, r.config.AWS.BucketName, countryKeyFormat)
	}
}