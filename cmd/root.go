package cmd

import (
	"os"
	"time"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/request"
	"github.com/nschimek/nice-fixture-feeder/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfig = "./config/default.yaml"
)

// stores context variables tied to the command
type commandContext struct {	
	startDate, endDate time.Time
	season bool
}

var (
	ctx *commandContext
	configFile string
	rootCmd = &cobra.Command{
		Use: "nice-fixture-feeder",
		Short: "Feed data into Nice Fixture by querying API-Football",
		Run: func(cmd *cobra.Command, args []string) {
			core.Log.Info("Started without commands or parameters, defaulting to Fixtures for today")
			ctx.startDate = time.Now()
			ctx.endDate = time.Now()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if cmd.Use == "version" {
				os.Exit(0) // the version command should do nothing more
			}

			// season mode requests leagues and teams
			if ctx.season {
				runSeasonRequest(request.Requests.League, request.Requests.Fixture)
			}
			// always run fixtures
			runFixturesRequest(request.Requests.Fixture, service.Services.TeamStats)
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", defaultConfig, "optional override of config file")
	rootCmd.PersistentFlags().IntP("season", "s", 0, "optional override of season set in config file")
	rootCmd.PersistentFlags().IntSliceP("leagues", "l", []int{0}, "optional override of league IDs set in config file")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode - use for more detailed logging")

	viper.BindPFlag("season", rootCmd.PersistentFlags().Lookup("season"))
	viper.BindPFlag("leagues", rootCmd.PersistentFlags().Lookup("leagues"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	cobra.OnInitialize(setup)
}

func setup() {
	core.Setup(configFile)
	repository.Setup(core.DB)
	service.Setup(core.Cfg, core.S3, repository.Repositories)
	request.Setup(core.Cfg, repository.Repositories, service.Services)

	ctx = new(commandContext)
}