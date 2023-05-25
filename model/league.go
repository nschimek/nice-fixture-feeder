package model

type League struct {
	League  LeagueLeague   `json:"league" gorm:"embedded"`
	Country LeagueCountry  `json:"country" gorm:"embedded;embeddedPrefix:country_"`
	Seasons []LeagueSeason `json:"seasons"`
	Audit   `json:"-"`
}

type LeagueLeague struct {
	Id   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Type string `json:"type" gorm:"type:enum('league', 'cup');default:'league'"`
	Logo string `json:"logo" gorm:"-"`
}

type LeagueCountry struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Flag string `json:"flag" gorm:"-"`
}

type LeagueSeason struct {
	LeagueId int       `json:"-" gorm:"primaryKey"`
	Season   int       `json:"year" gorm:"primaryKey"`
	Start    CivilTime `json:""`
	End      CivilTime `json:""`
	Current  bool      `json:""`
	Audit
}
