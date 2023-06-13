package service

import (
	"strings"

	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
)

//go:generate mockery --name FixtureStatusService  --filename fixture_status_mock.go
type FixtureStatusService interface {
	GetType(id string) string
	IsFinished(id string) bool
	IsScheduled(id string) bool
}

type fixtureStatusService struct {
	repo repository.FixtureStatusRepository
	idMap map[string]string
}

func NewFixtureStatusService(repo repository.FixtureStatusRepository) FixtureStatusService {
	return &fixtureStatusService{repo: repo}
}

func (s *fixtureStatusService) GetType(id string) string {
	if (s.idMap == nil) {
		s.initializeMap()
	}
	return s.idMap[strings.ToUpper(id)]
}

func (s *fixtureStatusService) IsFinished(id string) bool {
	return s.GetType(id) == model.StatusTypeFinished
}

func (s *fixtureStatusService) IsScheduled(id string) bool {
	return s.GetType(id) == model.StatusTypeScheduled
}

func (s *fixtureStatusService) initializeMap() {
	s.idMap = make(map[string]string)
	for _, fs := range s.repo.GetAll() {
		s.idMap[fs.Id] = fs.Type
	}
}

