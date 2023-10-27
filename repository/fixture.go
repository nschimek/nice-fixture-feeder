package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

const futureFixtureSQL = "league_id = ? AND season = ? and (team_home_id = ? or team_away_id = ?) and fixture_id >= ?"

type Fixture interface {
	UpsertRepository[model.Fixture]
	GetFutureFixturesByTLS(tlsId model.TeamLeagueSeasonId, minId int) ([]model.Fixture, error)
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
func (r *fixture) GetFutureFixturesByTLS(tlsId model.TeamLeagueSeasonId, minId int) ([]model.Fixture, error) {
	var fixtures []model.Fixture
	if err := r.db.Where(&fixtures, futureFixtureSQL, tlsId.LeagueId, 
		tlsId.Season, tlsId.TeamId, tlsId.TeamId, minId).Error; err != nil {
		return nil, err
	}
	return fixtures, nil
}
