package openrouter

type ToolType string

const (
	ToolTypeWebSearch ToolType = "web_search"
)

type Tool struct {
	Type ToolType `json:"type"`
}

func NewWebSearch() *Tool {
	return &Tool{
		Type: ToolTypeWebSearch,
	}
}
