package main

import (
	"fmt"
	"os"

	"github.com/eduardolat/openroutergo"
)

// You can copy this code to https://play.go.dev modify the api key, model, and run it.

const apiKey = "sk......."
const model = "google/gemini-2.0-flash-exp:free"

func main() {
	client, err := openroutergo.NewClient().WithAPIKey(apiKey).Create()
	if err != nil {
		fmt.Printf("Failed to create client: %v", err)
		os.Exit(1)
	}

	// Create and execute a completion
	completion, resp, err := client.
		NewChatCompletion().
		WithDebug(true).  // Enable debug mode to see the request and response in the console
		WithModel(model). // Change the model if you want
		WithSystemMessage("You are a helpful assistant expert in geography.").
		WithUserMessage("What is the capital of France?").
		Execute()
	if err != nil {
		fmt.Printf("Failed to execute completion: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Response: %v", resp.Choices[0].Message.Content)

	// Reuse the completion to continue the conversation
	_, resp, err = completion.
		WithAssistantMessage(resp.Choices[0].Message.Content).
		WithUserMessage("Thanks! Now, what is the capital of Germany?").
		Execute()
	if err != nil {
		fmt.Printf("Failed to execute completion: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Response: %v", resp.Choices[0].Message.Content)
}
