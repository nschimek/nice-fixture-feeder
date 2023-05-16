package response

type League struct {
	Base
	Response []struct{
		League struct{
			Id int `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			Logo string `json:"logo"`
		} `json:"league"`
		Country struct{
			Name string `json:"name"`
			Code string `json:"code"`
			Flag string `json:"flag"`
		} `json:"country"`
		Seasons []struct{
			Year int `json:"year"`
			Start string `json:"start"`
			End string `json:"end"`
			Current bool `json:"current"`
			Coverage struct {
				Fixtures struct {
					Events bool `json:"events"`
					Lineups bool `json:"lineups"`
					StatisticsFixtures bool `json:"statistics_fixtures"`
					StatisticsPlayers bool `json:"statistics_players"`
				} `json:"fixtures"`
				Standings bool `json:"standings"`
				Players bool `json:"players"`
				TopScorers bool `json:"top_scorers"`
				TopAssists bool `json:"top_assists"`
				TopCards bool `json:"top_cards"`
				Injuries bool `json:"injuries"`
				Predictions bool `json:"predictions"`
				Odds bool `json:"odds"`
			}
		} `json:"seasons"`
	} `json:"response"`
}