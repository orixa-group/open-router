package openrouter

import "encoding/json"

type UserMessage struct {
	Content []Content `json:"content"`
	Name    string    `json:"name,omitempty"`
}

func (m UserMessage) Role() string {
	return "user"
}

func (m UserMessage) MarshalJSON() ([]byte, error) {
	type Alias UserMessage
	return json.Marshal(&struct {
		Role string `json:"role"`
		Alias
	}{
		Role:  m.Role(),
		Alias: (Alias)(m),
	})
}
