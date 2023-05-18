package model

type League struct {
	League struct {
		Id int `json:"id" gorm:"primaryKey"`
		Name string `json:"name"`
		Type string `json:"type" gorm:"type:enum('league', 'cup');default:'league'"`
		Logo string `json:"logo" gorm:"-"`
	} `json:"league" gorm:"embedded"`
	Country struct {
		Name string `json:"name"`
		Code string `json:"code"`
		Flag string `json:"flag" gorm:"-"`
	} `json:"country" gorm:"embedded;embeddedPrefix:country_"`
	Seasons []LeagueSeason`json:"seasons"`
	Audit `json:"-"`
}

type LeagueSeason struct {
	LeagueId int `json:"-" gorm:"primaryKey"`
	Season int `json:"year" gorm:"primaryKey"`
	Start CivilTime `json:"start"`
	End CivilTime `json:"end"`
	Current bool `json:"current"`
	Audit `json:"-"`
}