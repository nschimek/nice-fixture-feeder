package main

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/request"
	"github.com/nschimek/nice-fixture-feeder/service"
)

const (
	defaultConfig = "./config/default.yaml"
)

func main() {
	core.SetupViper()
	core.SetupConfigFile(defaultConfig)
	core.SetupDatabase(core.Cfg)
	core.SetupS3(core.Cfg)

	imageService := service.NewImageService(core.S3)

	req := request.NewLeagueRequest(core.Cfg, &repository.LeagueRepository{DB: core.DB}, imageService)
	req.Request(2022, 39)
	req.Persist()
	req.PostPersist()
}