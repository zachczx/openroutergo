package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/eduardolat/openroutergo"
)

// This example demonstrates how to use JSON Schema Mode to ensure that the model's response
// matches a specific JSON schema. This guarantees that the response is not only valid JSON but
// also adheres to the defined structure.
//
// You can copy this code modify the api key, model, and run it.

const apiKey = "sk......."
const model = "google/gemini-2.0-flash-exp:free"

func main() {
	client, err := openroutergo.
		NewClient().
		WithAPIKey(apiKey).
		WithRefererURL("https://my-app.com"). // Optional, for rankings on openrouter.ai
		WithRefererTitle("My App").           // Optional, for rankings on openrouter.ai
		Create()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	completion := client.
		NewChatCompletion().
		WithDebug(true).  // Enable debug mode to see the request and response in the console
		WithModel(model). // Change the model if you want

		// The following is the basic JSON Mode that only guarantees the message the model
		// generates is valid JSON but it doesn't guarantee that the JSON matches
		// any specific schema. 👇
		//
		// WithResponseFormat(map[string]any{"type": "json_object"}).
		//
		// --------------------------------------------------------------------------------
		//
		// However, if you want to guarantee that the JSON matches a specific schema, you
		// can use the JSON Schema Mode. 👇
		WithResponseFormat(map[string]any{
			"type": "json_schema",
			"json_schema": map[string]any{
				"name": "capital_response",
				"schema": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"country": map[string]any{
							"type": "string",
						},
						"capital": map[string]any{
							"type": "string",
						},
						"curious_fact": map[string]any{
							"type": "string",
						},
					},
					"required": []string{"country", "capital", "curious_fact"},
				},
			},
		}).
		WithSystemMessage("You are a helpful assistant expert in geography.").
		WithUserMessage("What is the capital of France?")

	_, resp, err := completion.Execute()
	if err != nil {
		log.Fatalf("Failed to execute completion: %v", err)
	}

	// You can unmarshal the response to a struct because the model
	// should return a valid JSON and match the provided schema.
	//
	// However, is recommended to validate the response yourself to
	// avoid surprises, remember that the model can hallucinate.
	var myResponse struct {
		Country     string `json:"country"`
		Capital     string `json:"capital"`
		CuriousFact string `json:"curious_fact"`
	}
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &myResponse); err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	fmt.Printf(
		"The capital of %s is %s and here's a curious fact: %s\n",
		myResponse.Country,
		myResponse.Capital,
		myResponse.CuriousFact,
	)
}
