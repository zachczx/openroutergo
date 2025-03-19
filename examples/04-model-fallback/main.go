package main

import (
	"fmt"
	"os"

	"github.com/eduardolat/openroutergo"
)

// You can copy this code to https://play.go.dev modify the api key, model, and run it.

const apiKey = "sk......."
const baseModel = "google/gemini-2.0-flash-exp:free"
const firstFallbackModel = "google/gemini-2.0-flash-thinking-exp-1219:free"
const secondFallbackModel = "deepseek/deepseek-r1-zero:free"

// Paid model as last fallback if all free models fail
const thirdFallbackModel = "google/gemini-2.0-flash-001"

func main() {
	client, err := openroutergo.NewClient().WithAPIKey(apiKey).Create()
	if err != nil {
		fmt.Printf("Failed to create client: %v", err)
		os.Exit(1)
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
		fmt.Printf("Failed to execute completion: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Model used: %v", resp.Model)
	fmt.Printf("Response: %v", resp.Choices[0].Message.Content)
}
