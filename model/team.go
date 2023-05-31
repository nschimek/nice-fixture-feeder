package model

type Team struct {
	Team TeamTeam `json:"league" gorm:"embedded"`
	TeamLeagueSeason TeamLeagueSeason `json:"-"`
	TeamVenue TeamVenue `json:"venue" gorm:"embedded;embeddedPrefix:venue_"`
}

func (t *Team) SetTLS(leagueId, season int) {
	t.TeamLeagueSeason = TeamLeagueSeason{
		TeamId: t.Team.Id,
		LeagueId: leagueId,
		Season: season,
	}
}

type TeamTeam struct {
	Id int `gorm:"primaryKey"`
	Name, Code, Country string
	Logo string	`gorm:"-"`
	National bool
}

type TeamLeagueSeason struct {
	TeamId int `gorm:"primaryKey"`
	LeagueId int `gorm:"primaryKey"`
	Season int `gorm:"primaryKey"`
}

type TeamVenue struct {
	Name, Address, City string
	Capacity int
}