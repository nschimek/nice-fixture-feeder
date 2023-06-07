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
	GetTeamStatsMap() map[model.TeamStatsId]struct{}
}

type fixtureRequest struct {
	config *core.Config
	requester Requester[model.Fixture]
	repo repository.Repository[model.Fixture]
	RequestedData []model.Fixture
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
			r.RequestedData = append(r.RequestedData, fixtures...)
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
	rs := r.repo.Upsert(r.RequestedData)
	rs.LogErrors()
	rs.LogSuccesses()
}

func (r *fixtureRequest) GetData() []model.Fixture {
	return r.RequestedData
}

// we need to query the team stats endpoint for both the home and away teams stats as of the day before this game
func (r *fixtureRequest) GetTeamStatsMap() map[model.TeamStatsId]struct{} {
	m := make(map[model.TeamStatsId]struct{})

	for _, fixture := range r.RequestedData {
		m[fixture.GetTeamStatsId(true)] = core.Exists
		m[fixture.GetTeamStatsId(false)] = core.Exists
	}

	return m
}