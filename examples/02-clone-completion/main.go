package main

import (
	"fmt"
	"log"

	"github.com/eduardolat/openroutergo"
)

// In this example, we create a base chat completion that can be cloned and reused multiple times.
//
// This demonstrates how to set up a reusable chat completion configuration, clone it, and execute
// it with different user messages.
//
// You can copy this code to https://play.go.dev modify the api key, model, and run it.

const apiKey = "sk......."
const model = "google/gemini-2.0-flash-exp:free"

func main() {
	client, err := openroutergo.NewClient().WithAPIKey(apiKey).Create()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a base completion to be cloned and reused
	baseCompletion := client.
		NewChatCompletion().
		WithDebug(true).  // Enable debug mode to see the request and response in the console
		WithModel(model). // Change the model if you want
		WithSystemMessage("You are a helpful assistant expert in geography.")

	// Clone and execute the completion
	_, resp, err := baseCompletion.
		Clone().
		WithUserMessage("What is the capital of France?").
		Execute()
	if err != nil {
		log.Fatalf("Failed to execute completion: %v", err)
	}
	fmt.Println("Response:", resp.Choices[0].Message.Content)

	// Clone the original completion once again and execute it
	_, resp, err = baseCompletion.
		Clone().
		WithUserMessage("What is the capital of Germany?").
		Execute()
	if err != nil {
		log.Fatalf("Failed to execute completion: %v", err)
	}
	fmt.Println("Response:", resp.Choices[0].Message.Content)
}
