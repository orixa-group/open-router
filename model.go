package openrouter

type Model string

const (
	ModelGemini2_5FlashLite = "google/gemini-2.5-flash-lite"
	ModelGemini3FlashLite   = "google/gemini-3-flash-preview"
	ModelGemini3Pro         = "google/gemini-3-pro-preview"
	ModelClaudeSonnet4_5    = "anthropic/claude-sonnet-4.5"
	ModelChatGpt5_2         = "openai/gpt-5.2"
)

func (m Model) String() string {
	return string(m)
}
