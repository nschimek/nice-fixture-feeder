package service

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/repository"
)

type ServiceRegistry struct {
	FixtureStatus FixtureStatusService
	Image ImageService
	TeamStats TeamStatsService
}

func Setup(cfg *core.Config, s3 core.S3Client, repos *repository.RepositoryRegistry) *ServiceRegistry {
	return &ServiceRegistry{
		FixtureStatus: NewFixtureStatusService(repos.FixtureStatus),
		Image: NewImageService(s3),
		TeamStats: NewTeamStatsService(
			repos.TeamStats, 
			repos.TeamLeagueSeason, 
			NewFixtureStatusService(repos.FixtureStatus),
		),
	}
}