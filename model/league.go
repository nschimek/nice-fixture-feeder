package model

type League struct {
	League struct {
		Id int `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		Logo string `json:"logo"`
	} `json:"league"`
	Country struct {
		Name string `json:"name"`
		Code string `json:"code"`
		Flag string `json:"flag"`
	} `json:"country"`
	Seasons []LeagueSeason`json:"seasons"`
}

type LeagueSeason struct {
	LeagueId int
	Season int `json:"year"`
	Start string `json:"start"`
	End string `json:"end"`
	Current bool `json:"current"`
}