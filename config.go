package openrouter

// ReasoningEffort https://openrouter.ai/docs/api/api-reference/chat/send-chat-completion-request#request.body.reasoning
type ReasoningEffort string

const (
	ReasoningEffort_XHIGH   ReasoningEffort = "xhigh"
	ReasoningEffort_HIGH                    = "high"
	ReasoningEffort_MEDIUM                  = "medium"
	ReasoningEffort_LOW                     = "low"
	ReasoningEffort_MINIMAL                 = "minimal"
	ReasoningEffort_NONE                    = "none"
)
