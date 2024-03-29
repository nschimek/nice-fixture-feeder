package request

import (
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/sirupsen/logrus"
)

const fixturesEndpoint = "fixtures"

//go:generate mockery --name Fixture --filename fixture_mock.go
type Fixture interface {
	Request()
	RequestDateRange(startDate, endDate time.Time)
	Persist()
	GetMap() map[int]model.Fixture
	GetIds() []int
}

type fixture struct {
	config *core.Config
	requester Requester[model.Fixture]
	repo repository.UpsertRepository[model.Fixture]
	requestedData []model.Fixture
	fixtureMap map[int]model.Fixture
	fixtureIds []int
}

func NewFixture(config *core.Config, repo repository.UpsertRepository[model.Fixture]) Fixture {
	return &fixture{
		config: config,
		requester: NewRequester[model.Fixture](config),
		fixtureMap: make(map[int]model.Fixture),
		repo: repo,
	}
}

func (r *fixture) Request() {
	r.RequestDateRange(time.Time{}, time.Time{})
}

func (r *fixture) RequestDateRange(startDate, endDate time.Time) {
	core.Log.WithFields(logrus.Fields{
		"leagues": r.config.Leagues,
		"startDate": startDate.Format(core.YYYY_MM_DD),
		"endDate": startDate.Format(core.YYYY_MM_DD),
	}).Info("Requesting fixtures for leagues...")
	for leagueId := range core.IdArrayToMap(r.config.Leagues) {
		if fixtures, err := r.request(startDate, endDate, leagueId); err == nil {
			r.requestedData = append(r.requestedData, fixtures...)
		} else {
			core.Log.Errorf("Could not get fixtures for league ID %d: %v", leagueId, err)
		}
	}
}


func (r *fixture) request(startDate, endDate time.Time, leagueId int) ([]model.Fixture, error) {
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

func (r *fixture) Persist() {
	var err error
	r.requestedData, err = r.repo.Upsert(r.requestedData)
	if err == nil {
		r.postPersist()
	}
}

func (r *fixture) postPersist() {
	for _, fixture := range r.requestedData {
		r.fixtureIds = append(r.fixtureIds, fixture.Fixture.Id)
		r.fixtureMap[fixture.Fixture.Id] = fixture
	}
	// TODO: test this
	if !sort.IntsAreSorted(r.fixtureIds) {
		sort.Ints(r.fixtureIds)
	}
}

func (r *fixture) GetIds() []int {
	return r.fixtureIds
}

func (r *fixture) GetMap() map[int]model.Fixture {
	return r.fixtureMap
}