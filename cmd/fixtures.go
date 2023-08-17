package cmd

import (
	"time"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/request"
	"github.com/nschimek/nice-fixture-feeder/service"
	"github.com/spf13/cobra"
)

var (
	startDateString, endDateString  string
	fixtureCmd = &cobra.Command{
		Use: "fixtures",
		Short: "Request fixtures by date range",
		Long: `Request fixtures by date range, without setting up a new season.  Team stats will be maintained.
This should primarily be used if there are issues with the daily request that result in catch-up runs being necessary.
DO NOT ATTEMPT TO RUN DATE RANGES OUT OF ORDER.  This will fail.  In order re-runs will work, but stats will not be maintained.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if startDate, endDate, err := dateRangeFromStrings(startDateString, endDateString); err == nil {
				core.Log.Info(startDate, endDate)
				ctx.startDate = startDate
				ctx.endDate = endDate
			} else {
				return err
			}
			return nil
		},
	}
)

func init() {
	fixtureCmd.Flags().StringVarP(&startDateString, "start-date", "s", "", "start date for date range, use YYYY-MM-DD format")
	fixtureCmd.Flags().StringVarP(&endDateString, "end-date", "e", "", "end date for date range, use YYYY-MM-DD format")

	fixtureCmd.MarkFlagsRequiredTogether("start-date", "end-date")

	rootCmd.AddCommand(fixtureCmd)
}

func dateRangeFromStrings(startDateString, endDateString string) (startDate, endDate time.Time, err error) {
	startDate, err = stringToDate(startDateString)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err = stringToDate(endDateString)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return
}

func stringToDate(dateString string) (time.Time, error) {
	if dateString != "" {
		return time.ParseInLocation("2006-01-02", dateString, core.CST)
	} else {
		return time.Time{}, nil
	}
}

func runFixturesRequest(fixtureRequest request.Fixture, teamStatsService service.TeamStats) {
	if ctx.season {
		fixtureRequest.Request()
	} else {
		fixtureRequest.RequestDateRange(ctx.startDate, ctx.endDate)
	}
	fixtureRequest.Persist()

	teamStatsService.MaintainStats(fixtureRequest.GetIds(), fixtureRequest.GetMap())
	teamStatsService.Persist()
}