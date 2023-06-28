package main

import (
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/request"
	"github.com/nschimek/nice-fixture-feeder/service"
	"github.com/spf13/viper"
)

const (
	defaultConfig = "./config/default.yaml"
)

func main() {
	if viper.GetBool("serverless") {
		lambda.Start(LambdaHandler)
	} else {
		test()
	}
}

func init() {
	core.SetupViper()
	core.SetupConfigFile(defaultConfig)
	core.SetupDatabase(core.Cfg)
	core.SetupS3(core.Cfg)
}

func LambdaHandler() (string, error) {
	test()
	return "success", nil
}

func test() {
	// imageService := service.NewImageService(core.S3)

	// req := request.NewLeagueRequest(core.Cfg, repository.NewLeagueRepository(core.DB), imageService)
	// req.Request("39", "140")
	// req.Persist()
	// req.PostPersist()

	// req := request.NewTeamRequest(core.Cfg, &repository.TeamRepository{DB: core.DB}, imageService)
	// req.Request("39")
	// req.Persist()
	// req.PostPersist()

	start, _ := time.Parse(core.YYYY_MM_DD, "2022-09-02")
	end, _ := time.Parse(core.YYYY_MM_DD, "2022-10-31")
	req := request.NewFixtureRequest(core.Cfg, *repository.NewFixtureRepository(core.DB))
	req.Request(start, end)
	req.Persist()

	// repo := repository.NewFixtureStatusRepository(core.DB)
	// s := repo.GetAll()
	// fmt.Printf("%+v\n", s)

	// repo := repository.NewTeamLeagueSeasonRepository(core.DB)
	// s := repo.GetById(model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 33, LeagueId: 39, Season: 2022}})
	// fmt.Printf("%+v\n", s)

	stsService := service.NewFixtureStatusService(repository.NewFixtureStatusRepository(core.DB))
	tlsRepo := repository.NewTeamLeagueSeasonRepository(core.DB)
	tsRepo := repository.NewTeamStatsRepository(core.DB)

	tsService := service.NewTeamStatsService(tsRepo, tlsRepo, stsService)
	tsService.MaintainStats(req.GetIds(), req.GetMap())
	tsService.Persist()
}