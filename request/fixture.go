package request

import (
	"net/url"
	"sort"
	"strconv"
	"sync"
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
	done chan struct{}
}

func NewFixture(config *core.Config, repo repository.UpsertRepository[model.Fixture]) Fixture {
	return &fixture{
		config: config,
		requester: NewRequester[model.Fixture](config),
		fixtureMap: make(map[int]model.Fixture),
		repo: repo,
		done: make(chan struct{}),
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

	leagues := r.produceLeagues(r.config.Leagues)

	fixtures := make(chan []model.Fixture)
	errc := make(chan error, 1)

	// start-up one goroutine for each league (TODO: set a configurable max)
	var wg sync.WaitGroup
	for i := 0; i < len(r.config.Leagues); i++ {
		wg.Add(1)
		go func() {
			// listen for league Ids and request fixtures
			for leagueId := range leagues {
				res, err := r.request(startDate, endDate, leagueId)
				select {
				case fixtures <- res:
				case errc <- err:
				case <-r.done: return
				}
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(fixtures)
	}()

	for f := range fixtures {
		go func(f []model.Fixture) {
			rd, err := r.repo.Upsert(f)
			if err == nil {
				r.requestedData = append(r.requestedData, rd...)
			}
			select {		
			case errc <- err:
			case <-r.done: return
			}
		}(f)
	}
	
	if err := <-errc; err != nil {
		core.Log.Errorf("Could not get fixtures: %v", err)
	}

	// for leagueId := range core.IdArrayToMap(r.config.Leagues) {
	// 	if fixtures, err := r.request(startDate, endDate, leagueId); err == nil {
	// 		r.requestedData = append(r.requestedData, fixtures...)
	// 	} else {
	// 		core.Log.Errorf("Could not get fixtures for league ID %d: %v", leagueId, err)
	// 	}
	// }

}

func (r *fixture) produceLeagues(leagueIds []int) (<- chan int) {
	leagues := make(chan int)
	go func() {
		defer close(leagues)
		for leagueId := range core.IdArrayToMap(leagueIds) {
			select {
			case leagues <- leagueId:
			case <-r.done: return
			}
		}
	}()
	return leagues
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