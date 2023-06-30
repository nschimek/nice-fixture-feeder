package model

type TeamStats struct {
	Id TeamStatsId `gorm:"embedded"`
	NextFixtureId int
	TeamStatsFixtures TeamStatsFixtures `gorm:"embedded;embeddedPrefix:fixtures_"`
	TeamStatsGoals TeamStatsGoals `gorm:"embedded;embeddedPrefix:goals_"`
	GoalDifferential int
	Form string
	CleanSheets TeamStatsHomeAwayTotal `gorm:"embedded;embeddedPrefix:cs_"`
	FailedToScore TeamStatsHomeAwayTotal `gorm:"embedded;embeddedPrefix:fts_"`
	Points int
	Audit
}

type TeamStatsId struct {
	TeamId int `gorm:"primaryKey"`
	LeagueId int `gorm:"primaryKey"`
	Season int `gorm:"primaryKey"`
	FixtureId int `gorm:"primaryKey"`
}

func (t TeamStatsId) GetTlsId() TeamLeagueSeasonId {
	return TeamLeagueSeasonId{
		TeamId: t.TeamId,
		LeagueId: t.LeagueId,
		Season: t.Season,
	}
}


type TeamStatsFixtures struct {
	FixturesPlayed TeamStatsHomeAwayTotal `gorm:"embedded;embeddedPrefix:played_"`
	FixturesWins TeamStatsHomeAwayTotal `gorm:"embedded;embeddedPrefix:wins_"`
	FixturesDraws TeamStatsHomeAwayTotal `gorm:"embedded;embeddedPrefix:draws_"`
	FixturesLosses TeamStatsHomeAwayTotal `gorm:"embedded;embeddedPrefix:losses_"`
}

type TeamStatsGoals struct {
	GoalsFor TeamStatsHomeAwayTotal `gorm:"embedded;embeddedPrefix:for_"`
	GoalsAgainst TeamStatsHomeAwayTotal `gorm:"embedded;embeddedPrefix:against_"`
}

type TeamStatsHomeAwayTotal struct {
	Home, Away, Total int
}

func (t *TeamStatsHomeAwayTotal) Increment(amount int, home bool) {
	if home {
		t.Home = t.Home + amount
 	} else {
		t.Away = t.Away + amount
	}
	t.Total = t.Total + amount
}