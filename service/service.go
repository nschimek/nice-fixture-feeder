package service

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/service/scores"
)

type ServiceRegistry struct {
	FixtureStatus FixtureStatus
	Image Image
	TeamLeagueSeason TeamLeagueSeason
	TeamStats TeamStats
}

func Setup(cfg *core.Config, s3 core.S3Client, repos *repository.RepositoryRegistry, scores *scores.ScoreRegistry) *ServiceRegistry {
	fixtureStatus := NewFixtureStatus(repos.FixtureStatus)
	teamLeagueSeason := NewTeamLeagueSeason(repos.TeamLeagueSeason)
	return &ServiceRegistry{
		FixtureStatus: fixtureStatus,
		Image: NewImage(s3),
		TeamLeagueSeason: teamLeagueSeason,
		TeamStats: NewTeamStats(
			repos.TeamStats, 
			teamLeagueSeason,
			fixtureStatus,
		),
	}
}