package request

type Response[T any] struct {
	Get string `json:"get"`
	Parameters map[string]string `json:"parameters"`
	Errors map[string]string `json:"errors"`
	Paging struct {
		Current int `json:"current"`
		Total int `json:"total"`	
	} `json:"paging"`
	Response []T `json:"response"`
}