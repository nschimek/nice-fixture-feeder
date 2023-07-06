package service

import (
	"errors"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/sirupsen/logrus"
)

//go:generate mockery --name TeamStatsService --filename team_stats_mock.go
type TeamStatsService interface {
	MaintainStats(fixtureIds []int, fixtureMap map[int]model.Fixture)
	Persist()
}

type teamStatsService struct {
	tsRepo repository.TeamStatsRepository
	tlsRepo repository.TeamLeagueSeasonRepository
	statusService FixtureStatusService
	tlsMap map[model.TeamLeagueSeasonId]model.TeamLeagueSeason
	statsMap map[model.TeamStatsId]model.TeamStats
}

func NewTeamStatsService(tsRepo repository.TeamStatsRepository, 
	tlsRepo repository.TeamLeagueSeasonRepository,
	statusService FixtureStatusService) *teamStatsService {
	return &teamStatsService{
		tsRepo: tsRepo, 
		tlsRepo: tlsRepo, 
		statusService: statusService,
		tlsMap: make(map[model.TeamLeagueSeasonId]model.TeamLeagueSeason),
		statsMap: make(map[model.TeamStatsId]model.TeamStats),
	}
}

// fixtureIds MUST be sorted ascending and all fixture IDs must be present in the fixture map!
func (s *teamStatsService) MaintainStats(fixtureIds []int, fixtureMap map[int]model.Fixture) {
	for _, id := range fixtureIds {
		if fixture, ok := fixtureMap[id]; ok {
			if s.statusService.IsFinished(fixture.Fixture.Status.Id) {
				core.Log.Infof("Maintaining stats for finished fixture %d...", fixture.Fixture.Id)
				s.maintainFixture(&fixture, true)
				s.maintainFixture(&fixture, false)
			}
		} else {
			core.Log.Errorf("Fixture %d not found in Fixture Map!", id)
		}
	}
}

func (s *teamStatsService) Persist() {
	core.Log.WithFields(logrus.Fields{
		"team_stats": len(s.statsMap), "team_league_seasons": len(s.tlsMap),
	}).Info("Persisting Team Stats...")

	_, err := s.tsRepo.Upsert(core.MapToArray[model.TeamStatsId, model.TeamStats](s.statsMap))

	if err == nil {
		s.tlsRepo.Upsert(core.MapToArray[model.TeamLeagueSeasonId, model.TeamLeagueSeason](s.tlsMap))
	}
}

func (s *teamStatsService) maintainFixture(fixture *model.Fixture, home bool) {
	tsid := fixture.GetTeamStatsId(home) // Team Stats ID (TSID) is set from the INCOMING fixture
	tls, curr, prev, err := s.getUpdatedStats(tsid, fixture)

	if err == nil && tls != nil {
		// if there are no errors, make the ID updates and save these stats in the maps (they will get persisted later)
		tls.MaxFixtureId = fixture.Fixture.Id
		s.tlsMap[tls.Id] = *tls
		s.statsMap[curr.Id] = *curr
		// only persist non-zeroed previous stats (zeroed stats are used at start of season)
		if prev.Id.FixtureId > 0 {
			prev.NextFixtureId = fixture.Fixture.Id
			s.statsMap[prev.Id] = *prev
		}
	} else if err != nil {
		// just log that there were errors.  by not populating the maps, they will not be persisted.
		core.Log.Errorf("issues maintaing stats for fixture ID %d: %v", fixture.Fixture.Id, err)
	}
}

func (s *teamStatsService) getUpdatedStats(tsid *model.TeamStatsId, 
	fixture *model.Fixture) (tls *model.TeamLeagueSeason, curr, prev *model.TeamStats, err error) {
	tls, err = s.getTLS(tsid)

	if err != nil {
		return nil, nil, nil, err
	} else if tls.MaxFixtureId >= tsid.FixtureId {
		// this will allow completed games to be safely re-requested without causing errors or currupting the team_stats table
		core.Log.WithFields(logrus.Fields{
			"fixtureId": tsid.FixtureId,
			"maxFixtureId": tls.MaxFixtureId,
		}).Warn("Max Fixture ID is GTE incoming fixture ID (likely a re-run) - cannot maintain stats, skipping...")
		return nil, nil, nil, nil
	}

	prev, err = s.getPreviousStats(tls)

	if err != nil {
		return nil, nil, nil, err
	}

	curr, err = s.calculateCurrentStats(prev, fixture)

	if err != nil {
		return nil, nil, nil, err
	}

	return
}

// get Team League Season (TLS), which contains the current max fixture ID (which will be the last played fixture)
func (s *teamStatsService) getTLS(tsid *model.TeamStatsId) (*model.TeamLeagueSeason, error) {
	core.Log.WithFields(logrus.Fields{
		"teamId": tsid.TeamId, "leagueId": tsid.LeagueId, "season": tsid.Season,
	}).Debug("Getting team league season (TLS)...")

	id := tsid.GetTlsId() // use the incoming Team Stats ID to get the TLS ID, which is just one less field

	var tls *model.TeamLeagueSeason
	if mv, ok := s.tlsMap[id]; ok {
		tls = &mv // use the map value, since we have it
	} else {
		tls, _ = s.tlsRepo.GetById(model.TeamLeagueSeason{Id: id})
	}

	if tls == nil {
		return nil, errors.New("could not get TLS, was the league setup?")
	}
	
	return tls, nil
}

func (s *teamStatsService) getPreviousStats(tls *model.TeamLeagueSeason) (*model.TeamStats, error) {
	core.Log.WithFields(logrus.Fields{
		"teamId": tls.Id.TeamId, "leagueId": tls.Id.LeagueId, "season": tls.Id.Season, "prevFixtureId": tls.MaxFixtureId,
	}).Debug("Getting previous stats using TLS...")

	id := tls.GetTeamStatsId()

	if tls.MaxFixtureId == 0 {
		core.Log.Debug("max fixture ID was 0, using zeroed stats as previous!")
		return &model.TeamStats{Id: id}, nil
	}

	var stats *model.TeamStats
	if mv, ok := s.statsMap[id]; ok {
		stats = &mv // use the map value, since we have it
	} else {
		stats, _ = s.tsRepo.GetById(model.TeamStats{Id: id})
	}

	if stats == nil {
		return nil, errors.New("no stats for max fixture ID")
	}

	return stats, nil
}

func (s *teamStatsService) calculateCurrentStats(prev *model.TeamStats, fixture *model.Fixture) (*model.TeamStats, error) {
	core.Log.WithFields(logrus.Fields{
		"teamId": prev.Id.TeamId, "leagueId": prev.Id.LeagueId, "season": prev.Id.Season,
		"prevFixtureId": prev.Id.FixtureId, "currFixtureId": fixture.Fixture.Id,
	}).Debug("Calculating current stats...")
	copy := *prev // create a copy of the current
	curr := &copy // work with the pointer

	if prev.Id.FixtureId >= fixture.Fixture.Id {
		return nil, errors.New("previous fixture ID is GTE incoming fixture ID - out of order")
	}

	curr.Id.FixtureId = fixture.Fixture.Id // this is just a little important...
	rs := fixture.GetResultStats(curr.Id.TeamId)
	
	curr.TeamStatsFixtures.FixturesPlayed.Increment(1, rs.Home)
	
	// result (and points)
	if (rs.Result == model.ResultWin) {
		curr.TeamStatsFixtures.FixturesWins.Increment(1, rs.Home)
		curr.Points = curr.Points + 3
	} else if (rs.Result == model.ResultLoss) {
		curr.TeamStatsFixtures.FixturesLosses.Increment(1, rs.Home)
	} else {
		curr.TeamStatsFixtures.FixturesDraws.Increment(1, rs.Home)
		curr.Points = curr.Points + 1
	}
	// form
	curr.Form = curr.Form + rs.Result
	// goals 
	curr.TeamStatsGoals.GoalsFor.Increment(rs.GoalsFor, rs.Home)
	curr.TeamStatsGoals.GoalsAgainst.Increment(rs.GoalsAgainst, rs.Home)
	// clean sheets and failed to score
	if (rs.GoalsAgainst == 0) {
		curr.CleanSheets.Increment(1, rs.Home)
	}
	if (rs.GoalsFor == 0) {
		curr.FailedToScore.Increment(1, rs.Home)
	}
	// goal differential
	curr.GoalDifferential = curr.GoalDifferential + (rs.GoalsFor - rs.GoalsAgainst)

	return curr, nil
}

