package cmd

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version of the Nice Fixture Feeder",
	Long:  `Display the current version of the Nice Fixture Feeder`,
	Run: func(cmd *cobra.Command, args []string) {
		core.Log.Infof("Current Nice Fixture Feeder version: %s", core.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}