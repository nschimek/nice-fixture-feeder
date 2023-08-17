package cmd

import (
	"testing"

	"github.com/nschimek/nice-fixture-feeder/request/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestSeasonCommand(t *testing.T) {
	ctx = new(commandContext)

	seasonCmd.Run(&cobra.Command{}, []string{})

	assert.True(t, ctx.season)
}

func TestRunSeasonRequest(t *testing.T) {
	lr := &mocks.League{}
	tr := &mocks.Team{}

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