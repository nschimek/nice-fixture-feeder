package main

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/request"
)

const (
	defaultConfig = "./config/default.yaml"
)

func main() {
	core.SetupViper()
	core.SetupConfigFile(defaultConfig)
	core.SetupDatabase(core.Cfg)

	req := request.NewLeagueRequest(core.Cfg)
	req.Request(39, 2022)
}