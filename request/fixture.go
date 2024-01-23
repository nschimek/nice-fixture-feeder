package request

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/nschimek/nice-fixture-feeder/core/util"

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
	config        *core.Config
	requester     Requester[model.Fixture]
	repo          repository.UpsertRepository[model.Fixture]
	requestedData []model.Fixture
	fixtureMap    map[int]model.Fixture
	fixtureIds    []int
	ctx           context.Context
}

func NewFixture(ctx context.Context, config *core.Config, repo repository.UpsertRepository[model.Fixture]) Fixture {
	return &fixture{
		config:     config,
		requester:  NewRequester[model.Fixture](config),
		fixtureMap: make(map[int]model.Fixture),
		repo:       repo,
		ctx:        ctx,
	}
}

func (r *fixture) Request() {
	r.RequestDateRange(time.Time{}, time.Time{})
}

func (r *fixture) RequestDateRange(startDate, endDate time.Time) {
	core.Log.WithFields(logrus.Fields{
		"leagues":   r.config.Leagues,
		"startDate": startDate.Format(core.YYYY_MM_DD),
		"endDate":   startDate.Format(core.YYYY_MM_DD),
	}).Info("Requesting fixtures for leagues...")

	var cancel context.CancelFunc
	r.ctx, cancel = context.WithTimeout(r.ctx, time.Second*5)
	defer cancel()

	leagues := make(chan int)

	// 2-stage pipeline: 1. request fixtures, 2. persist fixtures
	fixtures := r.concurrentRequest(leagues, startDate, endDate)
	persisted := r.concurrentPersist(fixtures)

	r.produceLeagues(r.config.Leagues, leagues)

	for result := range persisted {
		r.requestedData = append(r.requestedData, result...)
	}
}

func (r *fixture) produceLeagues(leagueIds []int, leagues chan int) {
	go func() {
		defer close(leagues)
		for leagueId := range util.IdArrayToMap(leagueIds) {
			select {
			case leagues <- leagueId:
			case <-r.ctx.Done():
				return
			}
		}
	}()
}

func (r *fixture) concurrentRequest(lc <-chan int, startDate, endDate time.Time) <-chan []model.Fixture {
	fixtures := make(chan []model.Fixture)
	pool := util.NewWorkerPool(len(r.config.Leagues))

	pool.Go(func(worker int) error {
		for leagueId := range lc {
			core.Log.WithFields(logrus.Fields{"leagueId": leagueId, "worker": worker}).Debug("Requesting fixtures...")
			res, err := r.request(startDate, endDate, leagueId)
			if err != nil {
				return fmt.Errorf("(worker %d, league %d) requesting: %v", worker, leagueId, err)
			}
			select {
			case fixtures <- res:
			case <-r.ctx.Done():
				return r.ctx.Err()
			}
		}
		return nil
	})
	pool.Wait(func() {
		close(fixtures)
		pool.LogErrors("Errors occurred while requesting fixtures!")
	})

	return fixtures
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

func (r *fixture) concurrentPersist(fc <-chan []model.Fixture) <-chan []model.Fixture {
	persisted := make(chan []model.Fixture)
	pool := util.NewWorkerPool(len(r.config.Leagues))

	pool.Go(func(worker int) error {
		for fixtures := range fc {
			core.Log.WithField("worker", worker).Debug("Persisting fixtures...")
			pd, err := r.repo.Upsert(fixtures)
			if err != nil {
				return fmt.Errorf("(worker %d) persisting: %v", worker, err)
			}
			select {
			case persisted <- pd:
			case <-r.ctx.Done():
				return r.ctx.Err()
			}
		}
		return nil
	})
	pool.Wait(func() {
		close(persisted)
		pool.LogErrors("Errors occurred while persisting fixtures!")
	})
	return persisted
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
