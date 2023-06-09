package request

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/service"
)

type RequestRegistry struct {
	Fixture FixtureRequest
	League LeagueRequest
	Team TeamRequest
}

func Setup(cfg *core.Config, repos *repository.RepositoryRegistry, svcs *service.ServiceRegistry) *RequestRegistry {
	return &RequestRegistry{
		Fixture: NewFixtureRequest(cfg, repos.Fixture),
		League: NewLeagueRequest(cfg, repos.League, svcs.Image),
		Team: NewTeamRequest(cfg, repos.Team, svcs.Image),
	}
}