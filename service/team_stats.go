package service

import (
	"errors"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/sirupsen/logrus"
)

//go:generate mockery --name TeamStats --filename team_stats_mock.go
type TeamStats interface {
	GetById(tsid model.TeamStatsId) (*model.TeamStats, error)
	GetByIdWithTLS(tsid model.TeamStatsId, current bool) (*model.TeamStats, error)
	MaintainStats(fixtureIds []int, fixtureMap map[int]model.Fixture)
	Persist()
}

type teamStats struct {
	tsRepo repository.TeamStats
	tlsService TeamLeagueSeason
	statusService FixtureStatus
	tlsMap map[model.TeamLeagueSeasonId]model.TeamLeagueSeason
	statsMap map[model.TeamStatsId]model.TeamStats
}

func NewTeamStats(tsRepo repository.TeamStats, 
	tlsService TeamLeagueSeason,
	statusService FixtureStatus) *teamStats {
	return &teamStats{
		tsRepo: tsRepo, 
		tlsService: tlsService,
		statusService: statusService,
		tlsMap: make(map[model.TeamLeagueSeasonId]model.TeamLeagueSeason),
		statsMap: make(map[model.TeamStatsId]model.TeamStats),
	}
}

func (s *teamStats) GetById(id model.TeamStatsId) (*model.TeamStats, error)  {
	core.Log.WithFields(logrus.Fields{
		"teamId": id.TeamId, "leagueId": id.LeagueId, "season": id.Season,
	}).Debug("Getting team stats by ID...")

	var stats *model.TeamStats
	if mv, ok := s.statsMap[id]; ok {
		stats = &mv // use the map value, since we have it
	} else {
		stats, _ = s.tsRepo.GetById(id)
	}

	if stats == nil {
		return nil, errors.New("no stats for given ID")
	}

	s.AddToMap(stats)

	return stats, nil
}

func (s *teamStats) GetByIdWithTLS(id model.TeamStatsId, current bool) (*model.TeamStats, error) {
	tls, err := s.tlsService.GetById(id.GetTlsId())

	if err != nil {
		return nil, err
	}

	return s.GetById(*tls.GetTeamStatsId(current))
}

/*
	Add stats to the stats map
	how we add these IDs to the map is important.  we have to support look-up via current fixture ID and next fixture ID.
	to accomplish this, we zero out the unused ID field - this way, only one is required to get the stats
*/
func (s *teamStats) AddToMap(stats *model.TeamStats) {
	if stats.Id.FixtureId > 0 {
		s.statsMap[stats.Id.GetCurrentId()] = *stats
	}
	if stats.Id.NextFixtureId > 0 {
		s.statsMap[stats.Id.GetNextId()] = *stats
	}
}

// fixtureIds MUST be sorted ascending and all fixture IDs must be present in the fixture map!
func (s *teamStats) MaintainStats(fixtureIds []int, fixtureMap map[int]model.Fixture) {
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

func (s *teamStats) Persist() {
	core.Log.WithFields(logrus.Fields{
		"team_stats": len(s.statsMap), "team_league_seasons": len(s.tlsMap),
	}).Info("Persisting Team Stats...")

	_, err := s.tsRepo.Upsert(core.MapToArray[model.TeamStatsId, model.TeamStats](s.statsMap))

	if err == nil {
		s.tlsService.Persist()
	}
}

func (s *teamStats) maintainFixture(fixture *model.Fixture, home bool) {
	tsid := fixture.GetTeamStatsId(home) // Team Stats ID (TSID) is set from the INCOMING fixture
	tls, curr, prev, err := s.getUpdatedStats(tsid, fixture)

	if err == nil && tls != nil {
		// if there are no errors, make the ID updates and save these stats in the maps (they will get persisted later)
		tls.MaxFixtureId = fixture.Fixture.Id
		s.tlsService.AddToMap(tls)
		s.AddToMap(curr)
		// only persist non-zeroed previous stats (zeroed stats are used at start of season)
		if prev.Id.FixtureId > 0 {
			prev.Id.NextFixtureId = fixture.Fixture.Id
			s.AddToMap(prev)
		}
	} else if err != nil {
		// just log that there were errors.  by not populating the maps, they will not be persisted.
		core.Log.Errorf("issues maintaing stats for fixture ID %d: %v", fixture.Fixture.Id, err)
	}
}

func (s *teamStats) getUpdatedStats(tsid *model.TeamStatsId, 
	fixture *model.Fixture) (tls *model.TeamLeagueSeason, curr, prev *model.TeamStats, err error) {
	// get Team League Season (TLS), which contains the current max fixture ID (which will be the last played fixture)
	tls, err = s.tlsService.GetById(tsid.GetTlsId())

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

func (s *teamStats) getPreviousStats(tls *model.TeamLeagueSeason) (*model.TeamStats, error) {
	core.Log.WithFields(logrus.Fields{
		"teamId": tls.Id.TeamId, "leagueId": tls.Id.LeagueId, "season": tls.Id.Season, "prevFixtureId": tls.MaxFixtureId,
	}).Debug("Getting previous stats using TLS...")

	id := *tls.GetTeamStatsId(true)

	if tls.MaxFixtureId == 0 {
		core.Log.Debug("max fixture ID was 0, using zeroed stats as previous!")
		return &model.TeamStats{Id: id}, nil
	}

	return s.GetByIdWithTLS(id, true)
}

func (s *teamStats) calculateCurrentStats(prev *model.TeamStats, fixture *model.Fixture) (*model.TeamStats, error) {
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

