package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/eduardolat/openroutergo"
)

// This example demonstrates how to use a model that supports tools to get the weather
// of a city and continue the conversation using the tool's response.
//
// This requires a model that supports tools, you can find the list here:
// https://openrouter.ai/models?supported_parameters=tools&order=top-weekly
//
// You can copy this code to https://play.go.dev modify the api key, model, and run it.

const apiKey = "sk......."
const model = "google/gemini-2.0-flash-exp:free"

func getWeather(city string) string {
	// This is a fake function that returns a string but you can
	// do calculations, api calls, database queries, etc.
	return "It's cold and -120 celsius degrees in " + city + " right now. Literally freezing."
}

func main() {
	client, err := openroutergo.NewClient().WithAPIKey(apiKey).Create()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	completion := client.
		NewChatCompletion().
		WithDebug(true).  // Enable debug mode to see the request and response in the console
		WithModel(model). // Change the model if you want
		WithTool(openroutergo.ChatCompletionTool{
			Name:        "getWeather",
			Description: "Get the weather of a city, use this every time the user asks for the weather",
			Parameters: map[string]any{
				// The parameters definition should be a JSON object
				"type": "object",
				"properties": map[string]any{
					"city": map[string]any{
						"type": "string",
					},
				},
				"required": []string{"city"},
			},
		}).
		WithSystemMessage("You are a helpful assistant expert in geography.").
		WithUserMessage("I want to know the weather in the capital of Brazil and a joke about it")

	completion, resp, err := completion.Execute()
	if err != nil {
		log.Fatalf("Failed to execute completion: %v", err)
	}

	if len(resp.Choices) == 0 || len(resp.Choices[0].Message.ToolCalls) == 0 {
		log.Fatalf("No tool calls returned")
	}

	toolCall := resp.Choices[0].Message.ToolCalls[0]
	toolName := toolCall.Function.Name
	if toolName != "getWeather" {
		log.Fatalf("Unexpected tool name: %s", toolName)
	}

	toolCallArguments := toolCall.Function.Arguments
	args := map[string]any{}
	if err := json.Unmarshal([]byte(toolCallArguments), &args); err != nil {
		log.Fatalf("Failed to unmarshal tool call arguments: %v", err)
	}

	// Call the function with the arguments provided by the model
	weather := getWeather(args["city"].(string))

	// Use the tool response to continue the conversation
	_, resp, err = completion.
		WithToolMessage(toolCall, weather).
		Execute()
	if err != nil {
		log.Fatalf("Failed to execute completion: %v", err)
	}

	fmt.Println("Response:", resp.Choices[0].Message.Content)
}
