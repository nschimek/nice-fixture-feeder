package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

//go:generate mockery --name TeamStatsRepository --filename team_stats_mock.go
type TeamStatsRepository interface {
	UpsertRepository[model.TeamStats]
	GetByIdRepository[model.TeamStats, model.TeamStats]
	SaveRepository[model.TeamStats]
}

type teamStatsRepository struct {
	upsertRepository[model.TeamStats]
	getByIdRepository[model.TeamStats, model.TeamStats]
	saveRepository[model.TeamStats]
}

func NewTeamStatsRepository(db core.Database) *teamStatsRepository {
	return &teamStatsRepository{
		upsertRepository: upsertRepository[model.TeamStats]{
			db: db,
			label: "team_stats",
		},
		getByIdRepository: getByIdRepository[model.TeamStats, model.TeamStats]{db: db},
		saveRepository: saveRepository[model.TeamStats]{db: db},
	}
}