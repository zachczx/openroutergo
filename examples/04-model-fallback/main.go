package main

import (
	"fmt"
	"log"

	"github.com/eduardolat/openroutergo"
)

// In this example, we set up three fallback models. The idea is to use free models
// first, and only if they fail, a paid model is used as the last fallback.
//
// This demonstrates how to configure fallback models to ensure that the request
// is handled by a model and when it fails because the rate limits, context window,
// or other reasons, the request is automatically retried with a different model.
//
// You can copy this code to https://play.go.dev modify the api key, models, and run it.

const apiKey = "sk......."
const baseModel = "google/gemini-2.0-flash-exp:free"
const firstFallbackModel = "google/gemini-2.0-flash-thinking-exp-1219:free"
const secondFallbackModel = "deepseek/deepseek-r1-zero:free"
const thirdFallbackModel = "google/gemini-2.0-flash-001"

func main() {
	client, err := openroutergo.NewClient().WithAPIKey(apiKey).Create()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	completion := client.
		NewChatCompletion().
		WithDebug(true). // Enable debug mode to see the request and response in the console
		WithModel(baseModel).
		WithModelFallback(firstFallbackModel).
		WithModelFallback(secondFallbackModel).
		WithModelFallback(thirdFallbackModel).
		WithSystemMessage("You are a helpful assistant expert in geography.").
		WithUserMessage("What is the capital of France?")

	_, resp, err := completion.Execute()
	if err != nil {
		log.Fatalf("Failed to execute completion: %v", err)
	}

	fmt.Println("Model used:", resp.Model)
	fmt.Println("Provider used:", resp.Provider)
	fmt.Println("Response:", resp.Choices[0].Message.Content)
}
