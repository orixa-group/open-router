package openrouter

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateContent(t *testing.T) {
	type testResponse struct {
		Summary string   `json:"summary"`
		Tags    []string `json:"tags"`
	}

	type testCase[T any] struct {
		name          string
		roundTripFunc func(req *http.Request) (*http.Response, error)
		want          *T
		wantErr       error
	}

	testCases := []testCase[testResponse]{{
		name: "success",
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Status:     http.StatusText(200),
				Body: io.NopCloser(bytes.NewBufferString(`{
				"choices": [
					{
						"message": {
							"content": "{\"summary\": \"It works\", \"tags\": [\"go\", \"test\"]}"
						}
					}
				]
			}`)),
				Header: make(http.Header),
			}, nil
		},
		want: &testResponse{
			Summary: "It works",
			Tags:    []string{"go", "test"},
		},
	}, {
		name: "internal server error",
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Status:     http.StatusText(500),
				Body: io.NopCloser(bytes.NewBufferString(`{
					"error": {
						"message": "Internal server error"
					}
				}`)),
				Header: make(http.Header),
			}, nil
		},
		wantErr: errors.New("API error: 500 - Internal server error"),
	}, {
		name: "json unmarshal error",
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Status:     http.StatusText(200),
				Body: io.NopCloser(bytes.NewBufferString(`{
				"choices": [
					{
						"message": {
							"content": "{\"summary\": 123, \"tags\": [\"go\", \"test\"]}"
						}
					}
				]
			}`)),
				Header: make(http.Header),
			}, nil
		},
		wantErr: errors.New("error unmarshaling response"),
	}}

	originalTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = originalTransport
	}()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			http.DefaultClient.Transport = &MockTransport{
				RoundTripFunc: tt.roundTripFunc,
			}

			req := ChatCompletion[testResponse]().
				Use(ModelGemini2_5FlashLite).
				AppendMessages(UserMessage{Content: []Content{TextContent{Text: "Test"}}})

			result, err := req.GenerateContent("test-api-key")

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.want, result)
			}
		})
	}

	t.Run("schema generation error", func(t *testing.T) {
		res, err := ChatCompletion[map[string]any]().GenerateContent("test-api-key")

		assert.ErrorContains(t, err, "error generating schema")
		assert.Nil(t, res)
	})
}

func TestChatCompletionRequest_MarshalJSON(t *testing.T) {
	req := ChatCompletion[struct {
		Summary string   `json:"summary"`
		Tags    []string `json:"tags"`
	}]().
		Use(ModelGemini2_5FlashLite).
		AppendMessages(
			SystemMessage{Content: "System prompt"},
			UserMessage{Content: []Content{TextContent{Text: "User prompt"}}},
		)

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var resultMap map[string]any
	err = json.Unmarshal(data, &resultMap)
	assert.NoError(t, err)

	assert.Equal(t, string(ModelGemini2_5FlashLite), resultMap["model"])
	assert.Len(t, resultMap["messages"], 2)

	respFormat, ok := resultMap["response_format"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "json_schema", respFormat["type"])

	jsonSchema, ok := respFormat["json_schema"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "response", jsonSchema["name"])
	assert.Equal(t, true, jsonSchema["strict"])

	schemaObj, ok := jsonSchema["schema"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "object", schemaObj["type"])

	props, ok := schemaObj["properties"].(map[string]any)
	assert.True(t, ok)
	assert.Contains(t, props, "summary")
	assert.Contains(t, props, "tags")
}
