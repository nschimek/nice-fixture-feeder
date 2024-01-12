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
	GetById(tsid model.TeamStatsId, current bool) (*model.TeamStats, error)
	GetByIdWithTLS(tsid model.TeamStatsId, current bool) (*model.TeamStats, error)
	MaintainStats(fixtureIds []int, fixtureMap map[int]model.Fixture)
	GetMinFixtureMap() map[model.TeamLeagueSeasonId]int
	PersistOne(stats *model.TeamStats)
}

type teamStats struct {
	tsRepo        repository.TeamStats
	cache         core.Cache[model.TeamStats]
	tlsService    TeamLeagueSeason
	statusService FixtureStatus
	minFixtureMap map[model.TeamLeagueSeasonId]int
}

func NewTeamStats(tsRepo repository.TeamStats,
	cache core.Cache[model.TeamStats],
	tlsService TeamLeagueSeason,
	statusService FixtureStatus) *teamStats {
	return &teamStats{
		tsRepo:        tsRepo,
		cache:         cache,
		tlsService:    tlsService,
		statusService: statusService,
		minFixtureMap: make(map[model.TeamLeagueSeasonId]int),
	}
}

// GetById gets the Team Stats by ID.
// If current is True, the Fixture ID will be used.  If current is false, the Next Fixture ID will be used.
func (s *teamStats) GetById(id model.TeamStatsId, current bool) (*model.TeamStats, error) {
	core.Log.WithFields(logrus.Fields{
		"teamId": id.TeamId, "leagueId": id.LeagueId, "season": id.Season, "fixtureId": id.FixtureId, "nextFixtureId": id.NextFixtureId,
	}).Debug("Getting team stats by ID...")

	// check to make sure we have ana ID based on the boolean value (this is really just a sanity check)
	if id.FixtureId == 0 && current {
		return nil, errors.New("current is true but FixtureId is 0")
	} else if id.NextFixtureId == 0 && !current {
		return nil, errors.New("current is false but NextFixtureId is 0")
	}

	// recreate ID using the current setting - this will zero out the unused FixtureId (just in case)
	id = id.FromCurrent(current)

	var stats *model.TeamStats
	if cv, _ := s.cache.Get(id); cv != nil {
		stats = cv // use the cached value, since we have it
	} else {
		stats, _ = s.tsRepo.GetById(id)
		if stats == nil {
			return nil, errors.New("no stats for given ID")
		}
		s.cache.Set(id, stats)
	}

	return stats, nil
}

// Get Stats by TLS.  This will use the max fixture ID from TLS to obtain the stats.
// Use the current boolean to match on FixtureId (true) or NextFixtureID (false).
func (s *teamStats) GetByIdWithTLS(id model.TeamStatsId, current bool) (*model.TeamStats, error) {
	tls, err := s.tlsService.GetById(id.GetTlsId())

	if err != nil {
		return nil, err
	}

	return s.GetById(*tls.GetTeamStatsId(current), current)
}

func (s *teamStats) GetMinFixtureMap() map[model.TeamLeagueSeasonId]int {
	return s.minFixtureMap
}

// Add to the Minimum Fixture Map.
// The goal here is to keep track of the lowest fixture ID for each TLS that had stats maintained.
// This fixture ID, along with all future ones for the given TLS, will need to be re-scored by the Scoring Service.
func (s *teamStats) addToMinFixtureMap(stats *model.TeamStats) {
	if fid, ok := s.minFixtureMap[stats.Id.GetTlsId()]; !ok || fid > stats.Id.FixtureId {
		s.minFixtureMap[stats.Id.GetTlsId()] = stats.Id.FixtureId
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

func (s *teamStats) PersistOne(stats *model.TeamStats) {
	core.Log.WithFields(logrus.Fields{
		"teamId": stats.Id.TeamId, "leagueId": stats.Id.LeagueId, "season": stats.Id.Season, "fixtureId": stats.Id.FixtureId, "nextFixtureId": stats.Id.NextFixtureId,
	}).Info("Persisting team stats...")

	if stats.Id.FixtureId > 0 {
		s.cache.Set(stats.Id.GetCurrentId(), stats)
	}
	if stats.Id.NextFixtureId > 0 {
		s.cache.Set(stats.Id.GetNextId(), stats)
	}

	s.tsRepo.UpsertOne(*stats)
}

func (s *teamStats) maintainFixture(fixture *model.Fixture, home bool) {
	tsid := fixture.GetTeamStatsId(home) // Team Stats ID (TSID) is set from the INCOMING fixture
	tls, curr, prev, err := s.getUpdatedStats(tsid, fixture)

	if err == nil && tls != nil {
		// if there are no errors, make the ID updates and save these stats in the maps (they will get persisted later)
		tls.MaxFixtureId = fixture.Fixture.Id
		s.tlsService.PersistOne(tls)
		s.PersistOne(curr)
		s.addToMinFixtureMap(curr)
		// only persist non-zeroed previous stats (zeroed stats are used at start of season)
		if prev.Id.FixtureId > 0 {
			prev.Id.NextFixtureId = fixture.Fixture.Id
			s.PersistOne(prev)
		}
	} else if err != nil {
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
			"fixtureId":    tsid.FixtureId,
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
		return nil, nil, prev, err
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

	return s.GetById(id, true)
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
	if rs.Result == model.ResultWin {
		curr.TeamStatsFixtures.FixturesWins.Increment(1, rs.Home)
		curr.Points = curr.Points + 3
	} else if rs.Result == model.ResultLoss {
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
	if rs.GoalsAgainst == 0 {
		curr.CleanSheets.Increment(1, rs.Home)
	}
	if rs.GoalsFor == 0 {
		curr.FailedToScore.Increment(1, rs.Home)
	}
	// goal differential
	curr.GoalDifferential = curr.GoalDifferential + (rs.GoalsFor - rs.GoalsAgainst)

	return curr, nil
}
