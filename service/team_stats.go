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
	GetMinFixtureMap() map[model.TeamLeagueSeasonId]int
	Persist()
}

type teamStats struct {
	tsRepo repository.TeamStats
	tlsService TeamLeagueSeason
	statusService FixtureStatus
	statsMap map[model.TeamStatsId]*model.TeamStats // we use a pointer here because stats will be stored twice under different keys
	minFixtureMap map[model.TeamLeagueSeasonId]int 
}

func NewTeamStats(tsRepo repository.TeamStats, 
	tlsService TeamLeagueSeason,
	statusService FixtureStatus) *teamStats {
	return &teamStats{
		tsRepo: tsRepo, 
		tlsService: tlsService,
		statusService: statusService,
		statsMap: make(map[model.TeamStatsId]*model.TeamStats),
		minFixtureMap: make(map[model.TeamLeagueSeasonId]int),
	}
}

// Get Team Stats by ID.  Due to the way the stats are stored in the Map, populate either
// FixtureId or NextFixtureId - but not both.
func (s *teamStats) GetById(id model.TeamStatsId) (*model.TeamStats, error)  {
	core.Log.WithFields(logrus.Fields{
		"teamId": id.TeamId, "leagueId": id.LeagueId, "season": id.Season, "fixtureId": id.FixtureId, "nextFixtureId": id.NextFixtureId,
	}).Debug("Getting team stats by ID...")

	if id.FixtureId > 0 && id.NextFixtureId > 0 {
		return nil, errors.New("cannot have both FixtureId and NextFixtureId")
	}

	var stats *model.TeamStats
	if mv, ok := s.statsMap[id]; ok {
		stats = mv // use the map value, since we have it
	} else {
		stats, _ = s.tsRepo.GetById(id)
	}

	if stats == nil {
		return nil, errors.New("no stats for given ID")
	}

	s.addToStatsMap(stats)

	return stats, nil
}

// Get Stats by TLS.  This will use the max fixture ID from TLS to obtain the stats.
// Use the current boolean to match on FixtureId (true) or NextFixtureID (false).
func (s *teamStats) GetByIdWithTLS(id model.TeamStatsId, current bool) (*model.TeamStats, error) {
	tls, err := s.tlsService.GetById(id.GetTlsId())

	if err != nil {
		return nil, err
	}

	return s.GetById(*tls.GetTeamStatsId(current))
}

func (s *teamStats) GetMinFixtureMap() map[model.TeamLeagueSeasonId]int {
	return s.minFixtureMap
}

// Add stats to the Stats Map.
// But how we add these IDs to the stats map is important.  We have to support look-up via current fixture ID and next fixture ID.
// To accomplish this, we zero out the unused ID field - this way, only one is required to get the stats.
func (s *teamStats) addToStatsMap(stats *model.TeamStats) {
	if stats.Id.FixtureId > 0 {
		s.statsMap[stats.Id.GetCurrentId()] = stats
	}
	if stats.Id.NextFixtureId > 0 {
		s.statsMap[stats.Id.GetNextId()] = stats
	}
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

func (s *teamStats) Persist() {
	stats := []model.TeamStats{}

	// the map will potentially have two copies of stats (one keyed on current fixutre ID, the other on next fixture ID)
	// we only need to persist one copy
	for k, v := range s.statsMap {
		if k.FixtureId > 0 && k.NextFixtureId == 0 {
			stats = append(stats, *v)
		}
	}

	core.Log.WithField("team_stats", len(stats)).Info("Persisting Team Stats...")

	_, err := s.tsRepo.Upsert(stats)

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
		s.addToStatsMap(curr)
		s.addToMinFixtureMap(curr)
		// only persist non-zeroed previous stats (zeroed stats are used at start of season)
		if prev.Id.FixtureId > 0 {
			prev.Id.NextFixtureId = fixture.Fixture.Id
			s.addToStatsMap(prev)
		}
	} else if err != nil {
		// if we managed to get a previous stat, it will be in the map - so we need to remove it
		if prev != nil {
			core.Log.WithField("fixtureId", prev.Id.GetCurrentId()).Info("Removing prev stats from map due to error...")
			delete(s.statsMap, prev.Id.GetCurrentId())
		}
		// log that there were errors
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

