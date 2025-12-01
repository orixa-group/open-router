package openrouter

import "encoding/json"

type ImageContent struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

func (c ImageContent) Type() string {
	return "image_url"
}

func (c ImageContent) MarshalJSON() ([]byte, error) {
	type Alias ImageContent
	return json.Marshal(&struct {
		Type     string `json:"type"`
		ImageURL Alias  `json:"image_url"`
	}{
		Type:     c.Type(),
		ImageURL: (Alias)(c),
	})
}
