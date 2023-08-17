package service

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/repository"
)

type ServiceRegistry struct {
	FixtureStatus FixtureStatus
	Image Image
	TeamStats TeamStats
}

func Setup(cfg *core.Config, s3 core.S3Client, repos *repository.RepositoryRegistry) *ServiceRegistry {
	return &ServiceRegistry{
		FixtureStatus: NewFixtureStatus(repos.FixtureStatus),
		Image: NewImage(s3),
		TeamStats: NewTeamStats(
			repos.TeamStats, 
			repos.TeamLeagueSeason, 
			NewFixtureStatus(repos.FixtureStatus),
		),
	}
}