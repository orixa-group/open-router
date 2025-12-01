package openrouter

type Model string

const (
	ModelGemini2_5FlashLite = "google/gemini-2.5-flash-lite"
)

func (m Model) String() string {
	return string(m)
}
