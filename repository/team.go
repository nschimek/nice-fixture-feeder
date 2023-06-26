package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type TeamRepository struct {
	upsertRepository[model.Team]
}

func NewTeamRepository(db core.Database) *TeamRepository {
	return &TeamRepository{
		upsertRepository: upsertRepository[model.Team]{repository: newRepo(db, "teams")},
	}
}