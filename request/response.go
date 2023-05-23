package request

import "encoding/json"

type Response[T any] struct {
	Get string `json:""`
	Parameters MapOrEmptyArray `json:",omitempty"`
	Errors MapOrEmptyArray `json:",omitempty"`
	Paging struct {
		Current int `json:""`
		Total int `json:""`	
	} `json:"paging"`
	Response []T `json:""`
}

// The API returns an empty array for the Parameters and Errors fields, but a key-value object when populated.
type MapOrEmptyArray map[string]string

func (m *MapOrEmptyArray) UnmarshalJSON(data []byte) error {
	if string(data) == `[]` {
		return nil
	}

	type tmp MapOrEmptyArray
	return json.Unmarshal(data, (*tmp)(m))
}