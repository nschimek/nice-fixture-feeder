package model

import "time"

type TeamStatsId struct {
	TeamId int
	LeagueId int
	Season int
	FixtureId int
	Date time.Time
}