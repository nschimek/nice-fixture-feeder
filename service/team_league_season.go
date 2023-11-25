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
	PersistOne(tls *model.TeamLeagueSeason)
}

type teamLeagueSeason struct {
	tlsRepo repository.TeamLeagueSeason
	cache core.Cache[model.TeamLeagueSeason]
}

func NewTeamLeagueSeason(tlsRepo repository.TeamLeagueSeason, cache core.Cache[model.TeamLeagueSeason]) *teamLeagueSeason {
	return &teamLeagueSeason{tlsRepo: tlsRepo, cache: cache}
}

func (s *teamLeagueSeason) GetById(id model.TeamLeagueSeasonId) (*model.TeamLeagueSeason, error) {
	core.Log.WithFields(logrus.Fields{
		"teamId": id.TeamId, "leagueId": id.LeagueId, "season": id.Season,
	}).Debug("Getting team league season (TLS)...")

	var tls *model.TeamLeagueSeason
	if mv, _ := s.cache.Get(id); mv != nil {
		tls = mv // use the cached value, since we have it
	} else {
		tls, _ = s.tlsRepo.GetById(id)
	}

	if tls == nil {
		return nil, errors.New("could not get TLS - league may not be setup, or newly promoted team")
	}

	s.cache.Set(id, tls)
	
	return tls, nil
}

func (s *teamLeagueSeason) PersistOne(tls *model.TeamLeagueSeason) {
	res, err := s.tlsRepo.UpsertOne(*tls)
	if err == nil {
		s.cache.Set(res.Id, &res)
	}
}



