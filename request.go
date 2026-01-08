package openrouter

import "github.com/orixa-group/open-router/schema"

type reasoningConfig struct {
	Effort ReasoningEffort `json:"effort"`
}

type jsonSchema struct {
	Name   string         `json:"name"`
	Strict bool           `json:"strict"`
	Schema *schema.Schema `json:"schema"`
}

type responseFormat struct {
	Type       string     `json:"type"`
	JsonSchema jsonSchema `json:"json_schema"`
}

func newJSONResponseFormat(schema *schema.Schema) responseFormat {
	return responseFormat{
		Type: "json_schema",
		JsonSchema: jsonSchema{
			Name:   "response",
			Strict: true,
			Schema: schema,
		},
	}
}

type openRouterChatCompletionRequest struct {
	Model          Model            `json:"model"`
	Messages       []Message        `json:"messages"`
	ResponseFormat responseFormat   `json:"response_format"`
	Reasoning      *reasoningConfig `json:"reasoning,omitempty"`
}

func (o *openRouterChatCompletionRequest) SetReasoningEffort(value ReasoningEffort) {
	if len(value) > 0 {
		o.Reasoning = &reasoningConfig{Effort: value}
	} else {
		o.Reasoning = nil
	}
}

func NewOpenRouterChatCompletionRequest(model Model, schema *schema.Schema, messages ...Message) *openRouterChatCompletionRequest {
	return &openRouterChatCompletionRequest{
		Model:          model,
		Messages:       messages,
		ResponseFormat: newJSONResponseFormat(schema),
	}
}
