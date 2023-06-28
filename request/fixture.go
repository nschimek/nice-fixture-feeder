package request

import (
	"net/url"
	"strconv"
	"time"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
)

const fixturesEndpoint = "fixtures"

type FixtureRequest interface {
	RequestAll()
	Request(startDate, endDate time.Time)
	Persist()
	GetMap() map[int]model.Fixture
	GetIds() []int
}

type fixtureRequest struct {
	config *core.Config
	requester Requester[model.Fixture]
	repo repository.UpsertRepository[model.Fixture]
	requestedData []model.Fixture
	fixtureMap map[int]model.Fixture
	fixtureIds []int
}

func NewFixtureRequest(config *core.Config, repo repository.FixtureRepository) FixtureRequest {
	return &fixtureRequest{
		config: config,
		requester: NewRequester[model.Fixture](config),
		fixtureMap: make(map[int]model.Fixture),
		repo: repo,
	}
}

func (r *fixtureRequest) RequestAll() {
	r.Request(time.Time{}, time.Time{})
}

func (r *fixtureRequest) Request(startDate, endDate time.Time) {
	for leagueId := range core.IdArrayToMap(r.config.Leagues) {
		if fixtures, err := r.request(startDate, endDate, leagueId); err == nil {
			r.requestedData = append(r.requestedData, fixtures...)
		} else {
			core.Log.Errorf("Could not get fixtures for league ID %d: %v", leagueId, err)
		}
	}
}

func (r *fixtureRequest) request(startDate, endDate time.Time, leagueId int) ([]model.Fixture, error) {
	p := url.Values{}
	p.Add("league", strconv.Itoa(leagueId))
	p.Add("season", strconv.Itoa(r.config.Season))

	if !startDate.IsZero() && !endDate.IsZero() {
		p.Add("from", startDate.Format(core.YYYY_MM_DD))
		p.Add("to", endDate.Format(core.YYYY_MM_DD))
	} 

	resp, err := r.requester.Get(fixturesEndpoint, p)

	if err != nil {
		return nil, err
	}

	return resp.Response, nil
}

func (r *fixtureRequest) Persist() {
	var err error
	r.requestedData, err = r.repo.Upsert(r.requestedData)
	if err == nil {
		r.postPersist()
	}
}

func (r *fixtureRequest) postPersist() {
	for _, fixture := range r.requestedData {
		r.fixtureIds = append(r.fixtureIds, fixture.Fixture.Id)
		r.fixtureMap[fixture.Fixture.Id] = fixture
	}
}

func (r *fixtureRequest) GetIds() []int {
	return r.fixtureIds
}

func (r *fixtureRequest) GetMap() map[int]model.Fixture {
	return r.fixtureMap
}