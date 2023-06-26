package service

import (
	"errors"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/sirupsen/logrus"
)

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

	if err == nil {
		// if there are no errors, make the ID updates and save these stats in the maps (they will get persisted later)
		tls.MaxFixtureId = fixture.Fixture.Id
		s.tlsMap[tls.Id] = *tls
		s.statsMap[curr.TeamStatsId] = *curr
		// only persist non-zeroed previous stats (zeroed stats are used at start of season)
		if prev.TeamStatsId.FixtureId > 0 {
			prev.NextFixtureId = fixture.Fixture.Id
			s.statsMap[prev.TeamStatsId] = *prev
		}
	} else {
		// just log that there were errors.  by not populating the maps, they will not be persisted.
		core.Log.Errorf("issues with fixture ID %d: %v", fixture.Fixture.Id, err)
	}
}

func (s *teamStatsService) getUpdatedStats(tsid *model.TeamStatsId, 
	fixture *model.Fixture) (tls *model.TeamLeagueSeason, curr, prev *model.TeamStats, err error) {
	tls, err = s.getTLS(tsid)

	if err != nil {
		return nil, nil, nil, err
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
	} else if tls.MaxFixtureId >= tsid.FixtureId {
		core.Log.WithField("tlsMaxFixtureId", tls.MaxFixtureId).Error("Max Fixture ID is invalid")
		return nil, errors.New("TLS max fixture ID is GTE the incoming fixture ID - out of order")
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
		return &model.TeamStats{TeamStatsId: id}, nil
	}

	var stats *model.TeamStats
	if mv, ok := s.statsMap[id]; ok {
		stats = &mv // use the map value, since we have it
	} else {
		stats, _ = s.tsRepo.GetById(model.TeamStats{TeamStatsId: id})
	}

	if stats == nil {
		return nil, errors.New("no stats for max fixture ID")
	}

	return stats, nil
}

func (s *teamStatsService) calculateCurrentStats(prev *model.TeamStats, fixture *model.Fixture) (*model.TeamStats, error) {
	core.Log.WithFields(logrus.Fields{
		"teamId": prev.TeamStatsId.TeamId, "leagueId": prev.TeamStatsId.LeagueId, "season": prev.TeamStatsId.Season,
		"prevFixtureId": prev.TeamStatsId.FixtureId, "currFixtureId": fixture.Fixture.Id,
	}).Debug("Calculating current stats...")
	copy := *prev // create a copy of the current
	curr := &copy // work with the pointer

	if prev.TeamStatsId.FixtureId >= fixture.Fixture.Id {
		return nil, errors.New("previous fixture ID is GTE incoming fixture ID - out of order")
	}

	curr.TeamStatsId.FixtureId = fixture.Fixture.Id // this is just a little important...
	rs := fixture.GetResultStats(curr.TeamStatsId.TeamId)
	
	curr.TeamStatsFixtures.FixturesPlayed.Increment(1, rs.Home)
	
	// result
	if (rs.Result == model.ResultWin) {
		curr.TeamStatsFixtures.FixturesWins.Increment(1, rs.Home)
	} else if (rs.Result == model.ResultLoss) {
		curr.TeamStatsFixtures.FixturesLosses.Increment(1, rs.Home)
	} else {
		curr.TeamStatsFixtures.FixturesDraws.Increment(1, rs.Home)
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

