package request

import (
	"github.com/nschimek/nice-fixture-feeder/core/util"
	"net/url"
	"strconv"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/service"
)

const (
	teamsEndpoint = "teams"
	teamKeyFormat = "images/logos/teams/%s"
)

//go:generate mockery --name Team --filename team_mock.go
type Team interface {
	Request()
	Persist()
}

type team struct {
	config        *core.Config
	requester     Requester[model.Team]
	repo          repository.UpsertRepository[model.Team]
	imageService  service.Image
	requestedData []model.Team
}

func NewTeam(config *core.Config, repo repository.UpsertRepository[model.Team], is service.Image) Team {
	return &team{
		config:       config,
		requester:    NewRequester[model.Team](config),
		imageService: is,
		repo:         repo,
	}
}

func (r *team) Request() {
	core.Log.WithField("leagues", r.config.Leagues).Info("Requesting teams for leagues...")
	for leagueId := range util.IdArrayToMap(r.config.Leagues) {
		if teams, err := r.request(leagueId); err == nil {
			r.requestedData = append(r.requestedData, teams...)
		} else {
			core.Log.Errorf("Could not get teams for league %d: %v", leagueId, err)
		}
	}
}

func (r *team) request(leagueId int) ([]model.Team, error) {
	p := url.Values{}
	p.Add("league", strconv.Itoa(leagueId))
	p.Add("season", strconv.Itoa(r.config.Season))

	resp, err := r.requester.Get(teamsEndpoint, p)

	if err != nil {
		return nil, err
	}

	teams := make([]model.Team, len(resp.Response))
	for i, t := range resp.Response {
		t.SetTLS(leagueId, r.config.Season)
		teams[i] = t
	}

	return teams, nil
}

func (r *team) Persist() {
	var err error
	r.requestedData, err = r.repo.Upsert(r.requestedData)
	if err == nil {
		r.postPersist()
	}
}

func (r *team) postPersist() {
	for _, team := range r.requestedData {
		r.imageService.TransferURL(team.Team.Logo, r.config.AWS.BucketName, teamKeyFormat)
	}
}
