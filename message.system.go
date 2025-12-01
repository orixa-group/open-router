package openrouter

import "encoding/json"

type SystemMessage struct {
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

func (m SystemMessage) Role() string {
	return "system"
}

func (m SystemMessage) MarshalJSON() ([]byte, error) {
	type Alias SystemMessage
	return json.Marshal(&struct {
		Role string `json:"role"`
		Alias
	}{
		Role:  m.Role(),
		Alias: (Alias)(m),
	})
}
