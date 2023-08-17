package service

import (
	"strings"

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
	idMap map[string]string
}

func NewFixtureStatus(repo repository.FixtureStatus) *fixtureStatus {
	return &fixtureStatus{repo: repo}
}

func (s *fixtureStatus) GetType(id string) string {
	if (s.idMap == nil) {
		s.initializeMap()
	}
	return s.idMap[strings.ToUpper(id)]
}

func (s *fixtureStatus) IsFinished(id string) bool {
	return s.GetType(id) == model.StatusTypeFinished
}

func (s *fixtureStatus) IsScheduled(id string) bool {
	return s.GetType(id) == model.StatusTypeScheduled
}

func (s *fixtureStatus) initializeMap() {
	s.idMap = make(map[string]string)
	all, _ := s.repo.GetAll()
	for _, fs := range all {
		s.idMap[fs.Id] = fs.Type
	}
}

