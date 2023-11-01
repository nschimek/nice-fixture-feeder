package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

const futureFixtureSQL = "league_id = ? AND season = ? AND (team_home_id = ? OR team_away_id = ?) AND fixture_id >= ? AND fixture_id NOT IN ?"

type Fixture interface {
	UpsertRepository[model.Fixture]
	GetFutureFixturesByTLS(tlsId model.TeamLeagueSeasonId, minId int, notId []int) ([]model.Fixture, error)
}

type fixture struct {
	upsertRepository[model.Fixture]
	db core.Database
}

func NewFixture(db core.Database) *fixture {
	return &fixture{
		upsertRepository: upsertRepository[model.Fixture]{repository: newRepo(db, "fixtures")},
		db: db,
	}
}

// Find all fixtures with an ID GTE than Min ID with the same Season, League, and Team ID (home or away) using the TLS
// Excluding Fixture IDs provided in notId - this prevents returning Fixtures we already have
func (r *fixture) GetFutureFixturesByTLS(tlsId model.TeamLeagueSeasonId, minId int, notId []int) ([]model.Fixture, error) {
	var fixtures []model.Fixture
	if err := r.db.Where(&fixtures, futureFixtureSQL, tlsId.LeagueId, 
		tlsId.Season, tlsId.TeamId, tlsId.TeamId, minId, notId).Error; err != nil {
		return nil, err
	}
	return fixtures, nil
}
