package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type Team struct {
	upsertRepository[model.Team]
}

func NewTeam(db core.Database) *Team {
	return &Team{
		upsertRepository: upsertRepository[model.Team]{repository: newRepo(db, "teams")},
	}
}