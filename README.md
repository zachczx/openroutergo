# OpenRouterGo

A powerful, developer-friendly Go SDK for
[OpenRouter.ai](https://openrouter.ai) - the platform that gives you unified
access to 100+ AI models from OpenAI, Anthropic, Google, and more through a
single consistent API.

<p>
  <a href="https://pkg.go.dev/github.com/eduardolat/openroutergo">
    <img src="https://pkg.go.dev/badge/github.com/eduardolat/openroutergo" alt="Go Reference"/>
  </a>
  <a href="https://goreportcard.com/report/eduardolat/openroutergo">
    <img src="https://goreportcard.com/badge/eduardolat/openroutergo" alt="Go Report Card"/>
  </a>
  <a href="https://github.com/eduardolat/openroutergo/releases/latest">
    <img src="https://img.shields.io/github/release/eduardolat/openroutergo.svg" alt="Release Version"/>
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/github/license/eduardolat/openroutergo.svg" alt="License"/>
  </a>
  <a href="https://github.com/eduardolat/openroutergo">
    <img src="https://img.shields.io/github/stars/eduardolat/openroutergo?style=flat&label=github+stars"/>
  </a>
</p>

> [!WARNING]
> This client is not yet stable and the API signature may change in the future
> until it reaches version 1.0.0, so be careful when upgrading. However, the API
> signature should not change too much.

## Features

- 🚀 **Simple & Intuitive API** - Fluent builder pattern with method chaining
  for clean, readable code
- 🔄 **Smart Fallbacks** - Automatically retry with alternative models if your
  first choice fails or is rate-limited
- 🛠️ **Function Calling** - Let AI models access your tools and functions when
  needed
- 📊 **Structured Outputs** - Force responses in valid JSON format with schema
  validation
- 🧠 **Complete Control** - Fine-tune model behavior with temperature, top-p,
  frequency penalty and more
- 🔍 **Debug Mode** - Instantly see the exact requests and responses for easier
  development

## Installation

Go version 1.22 or higher is required.

```go
go get github.com/eduardolat/openroutergo
```

## Quick Start Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/eduardolat/openroutergo"
)

func main() {
	// Create a client with your API key
	client, err := openroutergo.
		NewClient().
		WithAPIKey("your-api-key").
		Create()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Build and execute your request with a fluent API
	completion, resp, err := client.
		NewChatCompletion().
		WithModel("google/gemini-2.0-flash-exp:free").
		WithSystemMessage("You are a helpful assistant expert in geography.").
		WithUserMessage("What is the capital of France?").
		Execute()
	if err != nil {
		log.Fatalf("Failed to execute completion: %v", err)
	}

	// Print the model's first response
	fmt.Println("Response:", resp.Choices[0].Message.Content)

	// Continue the conversation seamlessly, the last response is
	// automatically added to the conversation history
	_, resp, err = completion.
		WithUserMessage("That's great! What about Japan?").
		Execute()
	if err != nil {
		log.Fatalf("Failed to execute completion: %v", err)
	}

	// Print the model's second response
	fmt.Println("Response:", resp.Choices[0].Message.Content)
}
```

## More Examples

We've included several examples to help you get started quickly:

- [Basic Usage](examples/01-basic/main.go) - Simple chat completion with any
  model
- [Clone Completion](examples/02-clone-completion/main.go) - Reuse
  configurations for multiple requests
- [Conversation](examples/03-reuse-completion/main.go) - Build multi-turn
  conversations with context
- [Model Fallbacks](examples/04-model-fallback/main.go) - Gracefully handle rate
  limits and save money with alternative models
- [Function Calling](examples/05-function-calling/main.go) - Allow AI to call
  your application functions
- [Advanced Options](examples/06-other-options/main.go) - Explore additional
  parameters for fine-tuning
- [JSON Responses](examples/07-force-response-format/main.go) - Get structured,
  validated outputs

## Get Started

1. Get your API key from [OpenRouter.ai](https://openrouter.ai/keys)
2. Install the package: `go get github.com/eduardolat/openroutergo`
3. Start building with the examples above

## About me

I'm Eduardo, if you like my work please ⭐ star the repo and find me on the
following platforms:

- [X](https://x.com/eduardoolat)
- [GitHub](https://github.com/eduardolat)
- [LinkedIn](https://www.linkedin.com/in/eduardolat)
- [My Website](https://eduardo.lat)
- [Buy me a coffee](https://buymeacoffee.com/eduardolat)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.
