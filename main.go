package main

import (
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/request"
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

	// req := request.NewLeagueRequest(core.Cfg, &repository.LeagueRepository{DB: core.DB}, imageService)
	// req.Request(core.IdArrayToMap([]string{"39"}))
	// req.Persist()
	// req.PostPersist()

	// req := request.NewTeamRequest(core.Cfg, &repository.TeamRepository{DB: core.DB}, imageService)
	// req.Request("39")
	// req.Persist()
	// req.PostPersist()

	start, _ := time.Parse(core.YYYY_MM_DD, "2022-08-05")
	end, _ := time.Parse(core.YYYY_MM_DD, "2022-08-06")
	req := request.NewFixtureRequest(core.Cfg, &repository.FixtureRepository{DB: core.DB})
	req.Request(start, end, "39")
	req.Persist()
}