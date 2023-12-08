package cmd

import (
	"testing"
	"time"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	req_mocks "github.com/nschimek/nice-fixture-feeder/request/mocks"
	svc_mocks "github.com/nschimek/nice-fixture-feeder/service/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Fixture Command Tests will also test the date range methods
func TestFixturesCommandSuccess(t *testing.T) {
	ctx = new(commandContext)
	startDateString = "2023-07-04"
	endDateString = "2023-07-05"
	err := fixtureCmd.RunE(&cobra.Command{}, []string{})

	assert.Equal(t, time.Date(2023, 7, 4, 0, 0, 0, 0, core.CST), ctx.startDate)
	assert.Equal(t, time.Date(2023, 7, 5, 0, 0, 0, 0, core.CST), ctx.endDate)
	assert.Nil(t, err)
}

func TestFixturesCommandError(t *testing.T) {
	ctx = new(commandContext)
	// test invalid start date branch
	startDateString = "20230704"
	endDateString = "2023-07-05"
	err1 := fixtureCmd.RunE(&cobra.Command{}, []string{})

	assert.ErrorContains(t, err1, "cannot parse")

	// test invalid end date branch
	startDateString = "2023-07-04"
	endDateString = "20230705"
	err2 := fixtureCmd.RunE(&cobra.Command{}, []string{})

	assert.ErrorContains(t, err2, "cannot parse")
}

func TestRunFixturesRequestSeason(t *testing.T) {
	ids := []int{100, 101, 102}
	fixtures := map[int]model.Fixture{100: {}, 101: {}, 102: {}}
	ctx = &commandContext{season: true}

	fr := &req_mocks.Fixture{}
	ts := &svc_mocks.TeamStats{}
	ss := &svc_mocks.Scoring{}

	fr.EXPECT().Request()
	fr.EXPECT().Persist()
	fr.EXPECT().GetIds().Return(ids)
	fr.EXPECT().GetMap().Return(fixtures)
	ts.EXPECT().MaintainStats(ids, fixtures)

	runFixturesRequest(fr, ts, ss)

	fr.AssertCalled(t, "Request")
	fr.AssertCalled(t, "Persist")
	fr.AssertCalled(t, "GetIds")
	fr.AssertCalled(t, "GetMap")
	ts.AssertCalled(t, "MaintainStats", ids, fixtures)
	ts.AssertCalled(t, "Persist")
}

func TestRunFixturesRequestNoSeason(t *testing.T) {
	sd, ed := time.Date(2023, 7, 4, 0, 0, 0, 0, core.CST), time.Date(2023, 7, 5, 0, 0, 0, 0, core.CST)
	ids := []int{100, 101, 102}
	fixtures := map[int]model.Fixture{100: {}, 101: {}, 102: {}}
	ctx = &commandContext{season: false, startDate: sd, endDate: ed}

	fr := &req_mocks.Fixture{}
	ts := &svc_mocks.TeamStats{}
	ss := &svc_mocks.Scoring{}

	fr.EXPECT().RequestDateRange(sd, ed)
	fr.EXPECT().Persist()
	fr.EXPECT().GetIds().Return(ids)
	fr.EXPECT().GetMap().Return(fixtures)
	ts.EXPECT().MaintainStats(ids, fixtures)

	runFixturesRequest(fr, ts, ss)

	fr.AssertCalled(t, "RequestDateRange", sd, ed)
	fr.AssertCalled(t, "Persist")
	fr.AssertCalled(t, "GetIds")
	fr.AssertCalled(t, "GetMap")
	ts.AssertCalled(t, "MaintainStats", ids, fixtures)
	ts.AssertCalled(t, "Persist")
}
