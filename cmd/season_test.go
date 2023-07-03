package cmd

import (
	"testing"

	"github.com/nschimek/nice-fixture-feeder/request"
)

func TestRunSeasonRequest(t *testing.T) {
	lr := &request.MockLeagueRequest{}
	tr := &request.MockTeamRequest{}

	lr.EXPECT().Request()
	lr.EXPECT().Persist()
	tr.EXPECT().Request()
	tr.EXPECT().Persist()

	runSeasonRequest(lr, tr)

	lr.AssertCalled(t, "Request")
	lr.AssertCalled(t, "Persist")
	tr.AssertCalled(t, "Request")
	tr.AssertCalled(t, "Persist")
}