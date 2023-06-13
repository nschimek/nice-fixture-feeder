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
	Request(startDate, endDate time.Time, leagueIds... string)
	Persist()
	GetData() []model.Fixture
	GetById(id int) *model.Fixture
	GetIds() []int
}

type fixtureRequest struct {
	config *core.Config
	requester Requester[model.Fixture]
	repo repository.Repository[model.Fixture]
	requestedData []model.Fixture
	fixtureMap map[int]model.Fixture
	fixtureIds []int
}

func NewFixtureRequest(config *core.Config, repo repository.Repository[model.Fixture]) FixtureRequest {
	return &fixtureRequest{
		config: config,
		requester: NewRequester[model.Fixture](config),
		repo: repo,
	}
}

func (r *fixtureRequest) Request(startDate, endDate time.Time, leagueIds... string) {
	for leagueId := range core.IdArrayToMap(leagueIds) {
		if fixtures, err := r.request(startDate, endDate, leagueId); err == nil {
			r.requestedData = append(r.requestedData, fixtures...)
		} else {
			core.Log.Errorf("Could not get fixtures for league ID %s: %v", leagueId, err)
		}
	}
}

func (r *fixtureRequest) request(startDate, endDate time.Time, leagueId string) ([]model.Fixture, error) {
	p := url.Values{}
	p.Add("league", leagueId)
	p.Add("from", startDate.Format(core.YYYY_MM_DD))
	p.Add("to", endDate.Format(core.YYYY_MM_DD))
	p.Add("season", strconv.Itoa(r.config.Season))

	resp, err := r.requester.Get(fixturesEndpoint, p)

	if err != nil {
		return nil, err
	}

	return resp.Response, nil
}

func (r *fixtureRequest) Persist() {
	rs := r.repo.Upsert(r.requestedData)
	if rs.HasErrors() {
		r.requestedData = nil
	} else {
		rs.LogSuccesses()
		r.postPersist()
	}
}

func (r *fixtureRequest) postPersist() {
	if r.fixtureMap == nil {
		r.fixtureMap = make(map[int]model.Fixture)
	}
	for _, fixture := range r.requestedData {
		r.fixtureIds = append(r.fixtureIds, fixture.Fixture.Id)
		r.fixtureMap[fixture.Fixture.Id] = fixture
	}
}

func (r *fixtureRequest) GetData() []model.Fixture {
	return r.requestedData
}

func (r *fixtureRequest) GetIds() []int {
	return r.fixtureIds
}

func (r *fixtureRequest) GetById(id int) *model.Fixture {
	if f, ok := r.fixtureMap[id]; ok {
		return &f
	} else {
		return nil
	}
}