package request

import (
	"context"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/service"
)

type RequestRegistry struct {
	Fixture Fixture
	League  League
	Team    Team
}

func Setup(cfg *core.Config, repos *repository.RepositoryRegistry, svcs *service.ServiceRegistry) *RequestRegistry {
	ctx := context.Background()
	return &RequestRegistry{
		Fixture: NewFixture(ctx, cfg, repos.Fixture),
		League:  NewLeague(cfg, repos.League, svcs.Image),
		Team:    NewTeam(cfg, repos.Team, svcs.Image),
	}
}
