package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

//go:generate mockery --name TeamStatsRepository --filename team_stats_mock.go
type TeamStatsRepository interface {
	UpsertRepository[model.TeamStats]
	GetByIdRepository[model.TeamStats, model.TeamStats]
}

type teamStatsRepository struct {
	upsertRepository[model.TeamStats]
	getByIdRepository[model.TeamStats, model.TeamStats]
}

func NewTeamStatsRepository(db core.Database) *teamStatsRepository {
	r := newRepo(db, "team_stats")
	return &teamStatsRepository{
		upsertRepository: upsertRepository[model.TeamStats]{repository: r},
		getByIdRepository: getByIdRepository[model.TeamStats, model.TeamStats]{repository: r},
	}
}