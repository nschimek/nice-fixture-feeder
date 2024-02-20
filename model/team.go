package model

type Team struct {
	Team             TeamTeam         `gorm:"embedded"`
	TeamLeagueSeason TeamLeagueSeason `json:"-"`
	Venue            TeamVenue        `gorm:"embedded;embeddedPrefix:venue_"`
	Audit            `json:"-"`
}

func (t *Team) SetTLS(leagueId, season int) {
	t.TeamLeagueSeason = TeamLeagueSeason{
		Id: TeamLeagueSeasonId{
			TeamId:   t.Team.Id,
			LeagueId: leagueId,
			Season:   season,
		},
		MaxFixtureId: 0,
	}
}

func (t *Team) GetTLS(leagueId, season int) TeamLeagueSeason {
	return TeamLeagueSeason{
		Id: TeamLeagueSeasonId{
			TeamId:   t.Team.Id,
			LeagueId: leagueId,
			Season:   season,
		},
	}
}

type TeamTeam struct {
	Id                  int `gorm:"primaryKey"`
	Name, Code, Country string
	Logo                string `gorm:"-"`
	National            bool
}

type TeamLeagueSeason struct {
	Id           TeamLeagueSeasonId `gorm:"embedded"`
	MaxFixtureId int
	Audit
}

func (t TeamLeagueSeason) GetTeamStatsId(current bool) *TeamStatsId {
	tsid := &TeamStatsId{
		TeamId:   t.Id.TeamId,
		LeagueId: t.Id.LeagueId,
		Season:   t.Id.Season,
	}

	if current {
		tsid.FixtureId = t.MaxFixtureId
	} else {
		tsid.NextFixtureId = t.MaxFixtureId
	}

	return tsid
}

type TeamLeagueSeasonId struct {
	TeamId   int `gorm:"primaryKey"`
	LeagueId int `gorm:"primaryKey"`
	Season   int `gorm:"primaryKey"`
}

type TeamVenue struct {
	Name, Address, City string
	Capacity            int
}
