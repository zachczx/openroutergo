package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eduardolat/openroutergo"
)

// You can copy this code to https://play.go.dev modify the api key, model, and run it.

const apiKey = "sk......."
const model = "google/gemini-2.0-flash-exp:free"
const modelFallback = "deepseek/deepseek-r1-zero:free"

func main() {
	client, err := openroutergo.
		NewClient().
		WithAPIKey(apiKey).
		WithTimeout(10 * time.Minute).        // Set a timeout for the client requests
		WithRefererURL("https://my-app.com"). // Optional, for rankings on openrouter.ai
		WithRefererTitle("My App").           // Optional, for rankings on openrouter.ai
		Create()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// You can use your code editor to help you explore all the available options but
	// here are some of the most useful ones
	completion := client.
		NewChatCompletion().
		WithContext(context.Background()).
		WithDebug(true).
		WithModel(model).
		WithModelFallback(modelFallback).
		WithSeed(1234567890).
		WithTemperature(0.5).
		WithMaxPrice(0.5, 2).
		WithMaxTokens(1000).
		WithSystemMessage("You are a helpful assistant expert in geography.").
		WithUserMessage("What is the capital of France?")

	_, resp, err := completion.Execute()
	if err != nil {
		log.Fatalf("Failed to execute completion: %v", err)
	}

	fmt.Println("Response:", resp.Choices[0].Message.Content)
}
