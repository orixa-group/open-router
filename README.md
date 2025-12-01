# OpenRouter Go Client

A type-safe, structured output-oriented Go client for the [OpenRouter](https://openrouter.ai/) API.

## Features

- **Type-Safe Structured Outputs**: Leverage Go generics to define your expected response format. The client automatically generates the corresponding JSON schema for the LLM.
- **Builder Pattern**: Fluent API for constructing chat completion requests.
- **Support for Multi-modal Content**: Easily mix text and image content in messages.

## Installation

```bash
go get github.com/visiperf/open-router
```

## Usage

### Basic Example

Here's how to use the client to generate a structured response.

```go
package main

import (
	"fmt"
	"os"

	"github.com/visiperf/open-router"
)

// Define your desired output structure
type Response struct {
	Answer   string `json:"answer"`
	Confidence int    `json:"confidence"`
}

func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		panic("OPENROUTER_API_KEY is not set")
	}

	// Create a new chat completion request
	// The generic type argument specifies the expected response structure
	req := openrouter.ChatCompletion[Response]().
		Use(openrouter.ModelGemini2_5FlashLite).
		AppendMessages(
			openrouter.SystemMessage{
				Content: "You are a helpful assistant.",
			},
			openrouter.UserMessage{
				Content: []openrouter.Content{
					openrouter.TextContent{Text: "What is the capital of France?"},
				},
			},
		)

	// Generate content
	response, err := req.GenerateContent(apiKey)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Answer: %s\n", response.Answer)
	fmt.Printf("Confidence: %d%%\n", response.Confidence)
}
```

### Multi-modal Messages

You can include images in your user messages:

```go
openrouter.UserMessage{
    Content: []openrouter.Content{
        openrouter.TextContent{Text: "What is in this image?"},
        openrouter.ImageContent{
            URL: "https://example.com/image.jpg",
        },
    },
}
```

## License

MIT
