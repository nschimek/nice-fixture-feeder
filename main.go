package main

import (
	"github.com/nschimek/nice-fixture-feeder/core"
)

const (
	defaultConfig = "./config/default.yaml"
)

func main() {
	core.SetupViper()
	core.SetupConfigFile(defaultConfig)
}