package model

type Team struct {
	Team TeamTeam `gorm:"embedded"`
	TeamLeagueSeason TeamLeagueSeason `json:"-"`
	Venue TeamVenue `gorm:"embedded;embeddedPrefix:venue_"`
	Audit   `json:"-"`
}

func (t *Team) SetTLS(leagueId, season int) {
	t.TeamLeagueSeason = TeamLeagueSeason{
		Id: TeamLeagueSeasonId{
			TeamId: t.Team.Id,
			LeagueId: leagueId,
			Season: season,
		},
	}
}

type TeamTeam struct {
	Id int `gorm:"primaryKey"`
	Name, Code, Country string
	Logo string	`gorm:"-"`
	National bool
}

type TeamLeagueSeason struct {
	Id TeamLeagueSeasonId `gorm:"embedded"`
	MaxFixtureId int
	Audit
}

func (t TeamLeagueSeason) GetTeamStatsId() TeamStatsId {
	return TeamStatsId{
		TeamId: t.Id.TeamId,
		LeagueId: t.Id.LeagueId,
		Season: t.Id.Season,
		FixtureId: t.MaxFixtureId,
	}
}

type TeamLeagueSeasonId struct {
	TeamId int `gorm:"primaryKey"`
	LeagueId int `gorm:"primaryKey"`
	Season int `gorm:"primaryKey"`
}

type TeamVenue struct {
	Name, Address, City string
	Capacity int
}