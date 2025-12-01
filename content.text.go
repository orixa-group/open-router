package openrouter

import "encoding/json"

type TextContent struct {
	Text string `json:"text"`
}

func (c TextContent) Type() string {
	return "text"
}

func (c TextContent) MarshalJSON() ([]byte, error) {
	type Alias TextContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  c.Type(),
		Alias: (Alias)(c),
	})
}
