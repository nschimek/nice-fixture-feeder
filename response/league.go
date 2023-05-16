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
	} `json:"response"`
}