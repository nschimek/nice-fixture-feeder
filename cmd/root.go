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
	services *service.ServiceRegistry
	requests *request.RequestRegistry
	runSeasonFunc = runSeasonRequest
	runFixturesFunc = runFixturesRequest
	exitFunc = os.Exit
	rootCmd = &cobra.Command{
		Use: "nice-fixture-feeder",
		Short: "Feed data into Nice Fixture by querying API-Football",
		Long: `Queries API-Football for configured leagues, teams, and fixtures and loads all relevant data into the database.  
Maintains team stats as it loads in order played fixtures for the first time.  Re-runs of played fixtures will not have stats maintained.
If previously loaded played fixtures need stats re-maintained, use the season command to re-initialize the season instead.
Running without parameters will limit the query to yesterday and today's fixtures only.`,
		Run: func(cmd *cobra.Command, args []string) {
			core.Log.Info("Started without commands or parameters, defaulting to fixtures for yesterday and today")
			ctx.startDate = time.Now().AddDate(0, 0, -1)
			ctx.endDate = time.Now()
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == "version" {
				Exit(0) // the version command should do nothing more
				return nil
			}

			// season mode requests leagues and teams
			if ctx.season {
				core.Log.Warn("RUNNING IN INITIALIZE SEASON MODE")
				runSeasonFunc(requests.League, requests.Team)
			}
			// always run fixtures
			runFixturesFunc(requests.Fixture, services.TeamStats)

			return nil
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
)

// Application entry point from main.go
func Execute() error {
	return rootCmd.Execute()
}

func Exit(code int) {
	exitFunc(code)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", defaultConfig, "optional override of config file")
	rootCmd.PersistentFlags().IntP("season", "n", 0, "optional override of season set in config file")
	rootCmd.PersistentFlags().IntSliceP("leagues", "l", []int{0}, "optional override of league IDs set in config file")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode - use for more detailed logging")

	viper.BindPFlag("season", rootCmd.PersistentFlags().Lookup("season"))
	viper.BindPFlag("leagues", rootCmd.PersistentFlags().Lookup("leagues"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	cobra.OnInitialize(setup)
}

func setup() {
	core.Setup(configFile)
	
	repos := repository.Setup(core.DB)
	services = service.Setup(core.Cfg, core.S3, repos)
	requests = request.Setup(core.Cfg, repos, services)

	ctx = new(commandContext)
}