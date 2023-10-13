package service

import (
	"errors"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/sirupsen/logrus"
)

//go:generate mockery --name TeamLeagueSeason --filename team_league_season_mock.go
type TeamLeagueSeason interface {
	GetById(tsid model.TeamLeagueSeasonId) (*model.TeamLeagueSeason, error)
	AddToMap(tls *model.TeamLeagueSeason) 
	Persist()
}

type teamLeagueSeason struct {
	tlsRepo repository.TeamLeagueSeason
	tlsMap map[model.TeamLeagueSeasonId]model.TeamLeagueSeason
}

func NewTeamLeagueSeason(tlsRepo repository.TeamLeagueSeason) *teamLeagueSeason {
	return &teamLeagueSeason{
		tlsRepo: tlsRepo,
		tlsMap: make(map[model.TeamLeagueSeasonId]model.TeamLeagueSeason),
	}
}

func (s *teamLeagueSeason) GetById(id model.TeamLeagueSeasonId) (*model.TeamLeagueSeason, error) {
	core.Log.WithFields(logrus.Fields{
		"teamId": id.TeamId, "leagueId": id.LeagueId, "season": id.Season,
	}).Debug("Getting team league season (TLS)...")

	var tls *model.TeamLeagueSeason
	if mv, ok := s.tlsMap[id]; ok {
		tls = &mv // use the map value, since we have it
	} else {
		tls, _ = s.tlsRepo.GetById(model.TeamLeagueSeason{Id: id})
	}

	if tls == nil {
		return nil, errors.New("could not get TLS, was the league setup?")
	}

	s.AddToMap(tls) // cache in the map
	
	return tls, nil
}

func (s *teamLeagueSeason) AddToMap(tls *model.TeamLeagueSeason) {
	s.tlsMap[tls.Id] = *tls
}

func (s *teamLeagueSeason) Persist() {
	s.tlsRepo.Upsert(core.MapToArray[model.TeamLeagueSeasonId, model.TeamLeagueSeason](s.tlsMap))
}

