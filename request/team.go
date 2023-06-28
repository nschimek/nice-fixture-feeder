package request

import (
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

//go:generate mockery --name LeagueRequest
type TeamRequest interface {
	Request()
	Persist()
}

type teamRequest struct {
	config *core.Config
	requester Requester[model.Team]
	repo repository.UpsertRepository[model.Team]
	imageService service.ImageService
	requestedData []model.Team
}

func NewTeamRequest(config *core.Config, repo repository.TeamRepository, is service.ImageService) TeamRequest {
	return &teamRequest{
		config: config,
		requester: NewRequester[model.Team](config),
		imageService: is,
		repo: repo,
	}
}

func (r *teamRequest) Request() {
	for leagueId := range core.IdArrayToMap(r.config.Leagues) {
		if teams, err := r.request(leagueId); err == nil {
			r.requestedData = append(r.requestedData, teams...)
		} else {
			core.Log.Errorf("Could not get teams for league %d: %v", leagueId, err)
		}
	}
}

func (r *teamRequest) request(leagueId int) ([]model.Team, error) {
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

func (r *teamRequest) Persist() {
	var err error
	r.requestedData, err = r.repo.Upsert(r.requestedData)
	if err == nil {
		r.postPersist()
	}
}

func (r *teamRequest) postPersist() {
	for _, team := range r.requestedData {
		r.imageService.TransferURL(team.Team.Logo, r.config.AWS.BucketName, teamKeyFormat)
	}
}