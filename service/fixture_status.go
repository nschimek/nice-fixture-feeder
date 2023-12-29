package service

import (
	"strings"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
)

//go:generate mockery --name FixtureStatus --filename fixture_status_mock.go
type FixtureStatus interface {
	GetType(id string) string
	IsFinished(id string) bool
	IsScheduled(id string) bool
}

type fixtureStatus struct {
	repo repository.FixtureStatus
	cache core.Cache[model.FixtureStatus]
}

func NewFixtureStatus(repo repository.FixtureStatus, cache core.Cache[model.FixtureStatus]) *fixtureStatus {
	s := &fixtureStatus{repo: repo, cache: cache}
	s.initialize()
	return s
}

func (s *fixtureStatus) GetType(id string) string {
	sts, _ := s.cache.Get(strings.ToUpper(id))
	return sts.Type
}

func (s *fixtureStatus) IsFinished(id string) bool {
	return s.GetType(id) == model.StatusTypeFinished
}

func (s *fixtureStatus) IsScheduled(id string) bool {
	return s.GetType(id) == model.StatusTypeScheduled
}

func (s *fixtureStatus) initialize() {
	all, _ := s.repo.GetAll()
	for _, fs := range all {
		s.cache.Set(fs.Id, &fs)
	}
}

