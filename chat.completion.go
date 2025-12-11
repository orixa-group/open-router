package openrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/orixa-group/open-router/schema"
)

const (
	baseURL = "https://openrouter.ai/api/v1"
)

type ChatCompletionRequest[T any] struct {
	model    Model
	messages []Message
}

func ChatCompletion[T any]() *ChatCompletionRequest[T] {
	return &ChatCompletionRequest[T]{}
}

func (r *ChatCompletionRequest[T]) Use(model Model) *ChatCompletionRequest[T] {
	r.model = model
	return r
}

func (r *ChatCompletionRequest[T]) AppendMessages(messages ...Message) *ChatCompletionRequest[T] {
	r.messages = append(r.messages, messages...)
	return r
}

func (r ChatCompletionRequest[T]) GenerateContent(apiKey string) (*T, error) {
	return createChatCompletion(apiKey, r)
}

func (r ChatCompletionRequest[T]) MarshalJSON() ([]byte, error) {
	var t T
	s, err := schema.Generate(t)
	if err != nil {
		return nil, fmt.Errorf("error generating schema: %w", err)
	}

	return json.Marshal(struct {
		Model          Model     `json:"model"`
		Messages       []Message `json:"messages"`
		ResponseFormat struct {
			Type       string `json:"type"`
			JsonSchema struct {
				Name   string         `json:"name"`
				Strict bool           `json:"strict"`
				Schema *schema.Schema `json:"schema"`
			} `json:"json_schema"`
		} `json:"response_format"`
	}{
		Model:    r.model,
		Messages: r.messages,
		ResponseFormat: struct {
			Type       string `json:"type"`
			JsonSchema struct {
				Name   string         `json:"name"`
				Strict bool           `json:"strict"`
				Schema *schema.Schema `json:"schema"`
			} `json:"json_schema"`
		}{
			Type: "json_schema",
			JsonSchema: struct {
				Name   string         `json:"name"`
				Strict bool           `json:"strict"`
				Schema *schema.Schema `json:"schema"`
			}{
				Name:   "response",
				Strict: true,
				Schema: s,
			},
		},
	})
}

func createChatCompletion[T any](apiKey string, params ChatCompletionRequest[T]) (*T, error) {
	url := baseURL + "/chat/completions"

	payload, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var respErr map[string]Error
		if err := json.Unmarshal(body, &respErr); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w", err)
		}

		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, respErr["error"].Message)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			}
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	var t T
	if err := json.Unmarshal([]byte(result.Choices[0].Message.Content), &t); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &t, nil
}
