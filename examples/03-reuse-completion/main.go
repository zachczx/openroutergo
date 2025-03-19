package main

import (
	"fmt"
	"log"

	"github.com/eduardolat/openroutergo"
)

// You can copy this code to https://play.go.dev modify the api key, model, and run it.

const apiKey = "sk......."
const model = "google/gemini-2.0-flash-exp:free"

func main() {
	client, err := openroutergo.NewClient().WithAPIKey(apiKey).Create()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
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
		log.Fatalf("Failed to execute completion: %v", err)
	}
	fmt.Println("Response:", resp.Choices[0].Message.Content)

	// Reuse the completion to continue the conversation, the assistant message of the previous
	// completion is automatically added so you can continue the conversation easily.
	_, resp, err = completion.
		WithUserMessage("Thanks! Now, what is the capital of Germany?").
		Execute()
	if err != nil {
		log.Fatalf("Failed to execute completion: %v", err)
	}
	fmt.Println("Response:", resp.Choices[0].Message.Content)
}
