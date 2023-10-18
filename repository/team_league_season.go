package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

//go:generate mockery --name TeamLeagueSeason --filename team_league_season_mock.go
type TeamLeagueSeason interface {
	UpsertRepository[model.TeamLeagueSeason]
	GetByIdRepository[model.TeamLeagueSeason, model.TeamLeagueSeasonId]
}

type teamLeagueSeason struct {
	upsertRepository[model.TeamLeagueSeason]
	getByIdRepository[model.TeamLeagueSeason, model.TeamLeagueSeasonId]
}

func NewTeamLeagueSeason(db core.Database) *teamLeagueSeason {
	r := newRepo(db, "team_league_season")
	return &teamLeagueSeason{
		upsertRepository: upsertRepository[model.TeamLeagueSeason]{repository: r},
		getByIdRepository: getByIdRepository[model.TeamLeagueSeason, model.TeamLeagueSeasonId]{repository: r},
	}
}