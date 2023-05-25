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

type LeagueRequest struct {
	config *core.Config
	requester Requester[model.League]
	repo repository.Repository[model.League]
	imageService service.ImageService
	RequestedData []model.League
}

func NewLeagueRequest(config *core.Config, repo *repository.LeagueRepository, is service.ImageService) *LeagueRequest {
	return &LeagueRequest{
		config: config,
		requester: NewRequester[model.League](config),
		imageService: is,
		repo: repo,
	}
}

func (r *LeagueRequest) Request(season, id int) {
	p := url.Values{}
	p.Add("id", strconv.Itoa(id))
	p.Add("season", strconv.Itoa(season))

	resp, err := r.requester.Get(leaguesEndpoint, p)
	
	if err != nil {
		core.Log.Errorf("Could not get league %d: %v", id, err)
	} else {
		r.RequestedData = append(r.RequestedData, resp.Response...)
	}
}

func (r *LeagueRequest) Persist() {
	rs := r.repo.Upsert(r.RequestedData)
	rs.LogErrors()
	rs.LogSuccesses()
}

func (r *LeagueRequest) PostPersist() {
	for _, league := range r.RequestedData {
		r.imageService.TransferURL(league.League.Logo, r.config.AWS.BucketName, leagueKeyFormat)
		r.imageService.TransferURL(league.Country.Flag, r.config.AWS.BucketName, countryKeyFormat)
	}
}