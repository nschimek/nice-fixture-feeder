package cmd

import (
	"github.com/nschimek/nice-fixture-feeder/request"
	"github.com/spf13/cobra"
)

var (
	seasonCmd = &cobra.Command{
		Use: "season",
		Short: "Initialize a new season by requesting league(s), teams, and fixtures",
		Long: `Setup a new season by requesting configured or specified league(s) and the Teams that play in those leagues.
This can be run even if the season is already setup - it will simply update existing record(s).
All fixtures for the leagues within the given season will also be requested.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx.season = true
		},
	}
)

func init() {
	rootCmd.AddCommand(seasonCmd)
}

func runSeasonRequest(leagueRequest request.LeagueRequest, teamRequest request.TeamRequest) {
	leagueRequest.Request()
	leagueRequest.Persist()
	teamRequest.Request()
	teamRequest.Persist()
}