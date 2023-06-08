package request

import (
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

//go:generate mockery --name LeagueRequest
type LeagueRequest interface {
	Request(ids ...string)
	Persist()
	PostPersist()
	GetData() []model.League
}

type leagueRequest struct {
	config *core.Config
	requester Requester[model.League]
	repo repository.Repository[model.League]
	imageService service.ImageService
	RequestedData []model.League
}

func NewLeagueRequest(config *core.Config, repo repository.Repository[model.League], is service.ImageService) LeagueRequest {
	return &leagueRequest{
		config: config,
		requester: NewRequester[model.League](config),
		imageService: is,
		repo: repo,
	}
}

func (r *leagueRequest) Request(ids ...string) {
	for id := range core.IdArrayToMap(ids) {
		if leagues, err := r.request(id); err == nil {
			r.RequestedData = append(r.RequestedData, leagues...)
		} else {
			core.Log.Errorf("Could not get league %s: %v", id, err)
		}
	}
}

func (r *leagueRequest) request(id string) ([]model.League, error) {
	p := url.Values{}
	p.Add("id", id)
	p.Add("season", strconv.Itoa(r.config.Season))

	resp, err := r.requester.Get(leaguesEndpoint, p)

	if err != nil {
		return nil, err
	}

	return resp.Response, nil
}

func (r *leagueRequest) Persist() {
	rs := r.repo.Upsert(r.RequestedData)
	rs.LogErrors()
	rs.LogSuccesses()
}

func (r *leagueRequest) PostPersist() {
	for _, league := range r.RequestedData {
		r.imageService.TransferURL(league.League.Logo, r.config.AWS.BucketName, leagueKeyFormat)
		r.imageService.TransferURL(league.Country.Flag, r.config.AWS.BucketName, countryKeyFormat)
	}
}

func (r *leagueRequest) GetData() []model.League {
	return r.RequestedData
}