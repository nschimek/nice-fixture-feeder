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
	teamsEndpoint = "leagues"
	teamKeyFormat = "images/logos/teams/%s"
)

//go:generate mockery --name LeagueRequest
type TeamRequest interface {
	Request(idMap map[string]struct{})
	Persist()
	PostPersist()
	GetData() []model.Team
}

type teamRequest struct {
	config *core.Config
	requester Requester[model.Team]
	repo repository.Repository[model.Team]
	imageService service.ImageService
	RequestedData []model.Team
}

func NewTeamRequest(config *core.Config, repo repository.Repository[model.Team], is service.ImageService) TeamRequest {
	return &teamRequest{
		config: config,
		requester: NewRequester[model.Team](config),
		imageService: is,
		repo: repo,
	}
}

func (r *teamRequest) Request(idMap map[string]struct{}) {
	for _, leagueId := range core.IdMapToArray(idMap) {
		if teams, err := r.request(leagueId); err == nil {
			r.RequestedData = append(r.RequestedData, teams...)
		} else {
			core.Log.Errorf("Could not get teams for league %s: %v", leagueId, err)
		}
	}
}

func (r *teamRequest) request(leagueId string) ([]model.Team, error) {
	p := url.Values{}
	p.Add("leagueId", leagueId)
	p.Add("season", strconv.Itoa(r.config.Season))

	resp, err := r.requester.Get(teamsEndpoint, p)

	if err != nil {
		return nil, err
	}

	lid, _ := strconv.Atoi(leagueId)
	var teams []model.Team
	for i, t := range resp.Response {
		t.SetTLS(lid, r.config.Season)
		teams[i] = t
	}

	return teams, nil
}

func (r *teamRequest) Persist() {
	rs := r.repo.Upsert(r.RequestedData)
	rs.LogErrors()
	rs.LogSuccesses()
}

func (r *teamRequest) PostPersist() {
	for _, team := range r.RequestedData {
		r.imageService.TransferURL(team.Team.Logo, r.config.AWS.BucketName, teamKeyFormat)
	}
}

func (r *teamRequest) GetData() []model.Team {
	return r.RequestedData
}