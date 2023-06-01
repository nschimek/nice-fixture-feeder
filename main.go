package main

import (
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
	imageService := service.NewImageService(core.S3)

	// req := request.NewLeagueRequest(core.Cfg, &repository.LeagueRepository{DB: core.DB}, imageService)
	// req.Request(core.IdArrayToMap([]string{"39"}))
	// req.Persist()
	// req.PostPersist()

	req := request.NewTeamRequest(core.Cfg, &repository.TeamRepository{DB: core.DB}, imageService)
	req.Request(core.IdArrayToMap([]string{"39"}))
	req.Persist()
	req.PostPersist()
}