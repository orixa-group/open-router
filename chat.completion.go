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
	model     Model
	messages  []Message
	reasoning ReasoningEffort
}

func ChatCompletion[T any]() *ChatCompletionRequest[T] {
	return &ChatCompletionRequest[T]{}
}

func (r *ChatCompletionRequest[T]) WithReasoningEffort(value ReasoningEffort) *ChatCompletionRequest[T] {
	r.reasoning = value

	return r
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

	req := NewOpenRouterChatCompletionRequest(r.model, s, r.messages...)
	req.SetReasoningEffort(r.reasoning)

	return json.Marshal(req)
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

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var respErr apiError
		if err := json.Unmarshal(body, &respErr); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w", err)
		}

		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, respErr.Error.Message)
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
