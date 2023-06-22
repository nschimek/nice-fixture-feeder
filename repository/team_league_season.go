package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

//go:generate mockery --name TeamLeagueSeasonRepository --filename team_league_season_mock.go
type TeamLeagueSeasonRepository interface {
	GetByIdRepository[model.TeamLeagueSeason, model.TeamLeagueSeason]
	SaveRepository[model.TeamLeagueSeason]
	UpsertRepository[model.TeamLeagueSeason]
}

type teamLeagueSeasonRepository struct {
	getByIdRepository[model.TeamLeagueSeason, model.TeamLeagueSeason]
	saveRepository[model.TeamLeagueSeason]
	upsertRepository[model.TeamLeagueSeason]
}

func NewTeamLeagueSeasonRepository(db core.Database) *teamLeagueSeasonRepository {
	return &teamLeagueSeasonRepository{
		getByIdRepository: getByIdRepository[model.TeamLeagueSeason, model.TeamLeagueSeason]{db: db},
		saveRepository: saveRepository[model.TeamLeagueSeason]{db: db},
		upsertRepository: upsertRepository[model.TeamLeagueSeason]{
			db: db,
			label: "team_league_season",
		},
	}
}