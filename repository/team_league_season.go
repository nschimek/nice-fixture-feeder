package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

//go:generate mockery --name TeamLeagueSeasonRepository --filename team_league_season_mock.go
type TeamLeagueSeasonRepository interface {
	UpsertRepository[model.TeamLeagueSeason]
	GetByIdRepository[model.TeamLeagueSeason, model.TeamLeagueSeason]
}

type teamLeagueSeasonRepository struct {
	upsertRepository[model.TeamLeagueSeason]
	getByIdRepository[model.TeamLeagueSeason, model.TeamLeagueSeason]
}

func NewTeamLeagueSeasonRepository(db core.Database) *teamLeagueSeasonRepository {
	r := newRepo(db, "team_league_season")
	return &teamLeagueSeasonRepository{
		upsertRepository: upsertRepository[model.TeamLeagueSeason]{repository: r},
		getByIdRepository: getByIdRepository[model.TeamLeagueSeason, model.TeamLeagueSeason]{repository: r},
	}
}