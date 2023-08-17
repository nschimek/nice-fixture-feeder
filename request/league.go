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

//go:generate mockery --name League --filename league_mock.go
type League interface {
	Request()
	Persist()
}

type league struct {
	config *core.Config
	requester Requester[model.League]
	repo repository.UpsertRepository[model.League]
	imageService service.Image
	requestedData []model.League
}

func NewLeague(config *core.Config, repo repository.UpsertRepository[model.League], is service.Image) League {
	return &league{
		config: config,
		requester: NewRequester[model.League](config),
		imageService: is,
		repo: repo,
	}
}

func (r *league) Request() {
	core.Log.WithField("leagues", r.config.Leagues).Info("Requesting leagues...")
	for id := range core.IdArrayToMap(r.config.Leagues) {
		if leagues, err := r.request(id); err == nil {
			r.requestedData = append(r.requestedData, leagues...)
		} else {
			core.Log.Errorf("Could not get league %d: %v", id, err)
		}
	}
}

func (r *league) request(id int) ([]model.League, error) {
	p := url.Values{}
	p.Add("id", strconv.Itoa(id))
	p.Add("season", strconv.Itoa(r.config.Season))

	resp, err := r.requester.Get(leaguesEndpoint, p)

	if err != nil {
		return nil, err
	}

	return resp.Response, nil
}

func (r *league) Persist() {
	var err error
	r.requestedData, err = r.repo.Upsert(r.requestedData)
	if err == nil {
		r.postPersist()
	}
}

func (r *league) postPersist() {
	for _, league := range r.requestedData {
		r.imageService.TransferURL(league.League.Logo, r.config.AWS.BucketName, leagueKeyFormat)
		r.imageService.TransferURL(league.Country.Flag, r.config.AWS.BucketName, countryKeyFormat)
	}
}
