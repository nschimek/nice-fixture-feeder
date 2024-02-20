package request

import (
	"net/url"
	"strconv"

	"github.com/nschimek/nice-fixture-feeder/core/util"

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
	tlsRepo       repository.UpsertRepository[model.TeamLeagueSeason]
	imageService  service.Image
	requestedData []model.Team
	tlsData       []model.TeamLeagueSeason
}

func NewTeam(config *core.Config,
	repo repository.UpsertRepository[model.Team],
	tlsRepo repository.UpsertRepository[model.TeamLeagueSeason],
	is service.Image) Team {
	return &team{
		config:       config,
		requester:    NewRequester[model.Team](config),
		imageService: is,
		repo:         repo,
		tlsRepo:      tlsRepo,
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

	for _, t := range resp.Response {
		r.tlsData = append(r.tlsData, t.GetTLS(leagueId, r.config.Season))
	}

	return resp.Response, nil
}

func (r *team) Persist() {
	var err error

	// persist teams first...
	r.requestedData, err = r.repo.Upsert(r.requestedData)

	if err != nil {
		core.Log.Errorf("Could not persist Teams from Team request: %v", err)
		return
	}

	// ... then do TLS (gorm should have been able to do this automatically but it wouldn't work)
	_, err = r.tlsRepo.Upsert(r.tlsData)

	if err == nil {
		r.postPersist()
	} else {
		core.Log.Errorf("Could not persist TLS from Team request: %v", err)
	}
}

func (r *team) postPersist() {
	for _, team := range r.requestedData {
		r.imageService.TransferURL(team.Team.Logo, r.config.AWS.BucketName, teamKeyFormat)
	}
}
