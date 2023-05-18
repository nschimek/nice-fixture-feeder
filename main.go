package main

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/request"
)

const (
	defaultConfig = "./config/default.yaml"
)

func main() {
	core.SetupViper()
	core.SetupConfigFile(defaultConfig)
	core.SetupDatabase(core.Cfg)

	req := request.NewLeagueRequest(core.Cfg, &repository.LeagueRepository{DB: core.DB})
	req.Request(2022, 39)
	req.Persist()
}