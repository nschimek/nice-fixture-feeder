package model


type TeamStatsId struct {
	TeamId int `gorm:"primaryKey"`
	LeagueId int `gorm:"primaryKey"`
	Season int `gorm:"primaryKey"`
	FixtureId int `gorm:"primaryKey"`
}

type TeamStats struct {
	TeamStatsId TeamStatsId `gorm:"embedded"`
	TeamStatsFixtures TeamStatsFixtures `gorm:"embedded;embeddedPrefix:fixtures_"`
	TeamStatsGoals TeamStatsGoals `gorm:"embedded;embeddedPrefix:goals_"`
	GoalDifferential int
}

type TeamStatsFixtures struct {
	FixturesPlayed TeamStatsFixturesTotals `gorm:"embedded;embeddedPrefix:played_"`
	FixturesWins TeamStatsFixturesTotals `gorm:"embedded;embeddedPrefix:wins_"`
	FixturesDraws TeamStatsFixturesTotals `gorm:"embedded;embeddedPrefix:draws_"`
	FixturesLosses TeamStatsFixturesTotals `gorm:"embedded;embeddedPrefix:losses_"`
}

type TeamStatsFixturesTotals struct {
	Home, Away, Total int
}

type TeamStatsGoals struct {
	GoalsHome TeamStatsGoalsTotals `gorm:"embedded;embeddedPrefix:home_"`
	GoalsAway TeamStatsGoalsTotals `gorm:"embedded;embeddedPrefix:away_"`
}

type TeamStatsGoalsTotals struct {
	Home, Away, Total int
}

