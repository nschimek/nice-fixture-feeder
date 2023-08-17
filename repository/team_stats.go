package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

//go:generate mockery --name TeamStats --filename team_stats_mock.go
type TeamStats interface {
	UpsertRepository[model.TeamStats]
	GetByIdRepository[model.TeamStats, model.TeamStats]
}

type teamStats struct {
	upsertRepository[model.TeamStats]
	getByIdRepository[model.TeamStats, model.TeamStats]
}

func NewTeamStats(db core.Database) *teamStats {
	r := newRepo(db, "team_stats")
	return &teamStats{
		upsertRepository: upsertRepository[model.TeamStats]{repository: r},
		getByIdRepository: getByIdRepository[model.TeamStats, model.TeamStats]{repository: r},
	}
}