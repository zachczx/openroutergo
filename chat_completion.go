package openroutergo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eduardolat/openroutergo/internal/optional"
	"github.com/orsinium-labs/enum"
)

// chatCompletionRole is an enum for the role of a message in a chat completion.
type chatCompletionRole enum.Member[string]

// MarshalJSON implements the json.Marshaler interface for chatCompletionRole.
func (ccr chatCompletionRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(ccr.Value)
}

// UnmarshalJSON implements the json.Unmarshaler interface for chatCompletionRole.
func (ccr *chatCompletionRole) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*ccr = chatCompletionRole{Value: value}
	return nil
}

var (
	// RoleSystem is the role of a system message in a chat completion.
	RoleSystem = chatCompletionRole{"system"}
	// RoleUser is the role of a user message in a chat completion.
	RoleUser = chatCompletionRole{"user"}
	// RoleAssistant is the role of an assistant message in a chat completion.
	RoleAssistant = chatCompletionRole{"assistant"}
)

// NewChatCompletion creates a new chat completion request builder for the OpenRouter API.
//
// Docs:
//   - Reference: https://openrouter.ai/docs/api-reference/chat-completion
//   - Request: https://openrouter.ai/docs/api-reference/overview#completions-request-format
//   - Parameters: https://openrouter.ai/docs/api-reference/parameters
//   - Response: https://openrouter.ai/docs/api-reference/overview#completionsresponse-format
func (c *Client) NewChatCompletion() *chatCompletionBuilder {
	return &chatCompletionBuilder{
		client:             c,
		ctx:                context.Background(),
		model:              optional.String{IsSet: false},
		messages:           []chatCompletionMessage{},
		tools:              []ChatCompletionTool{},
		temperature:        optional.Float64{IsSet: false},
		topP:               optional.Float64{IsSet: false},
		topK:               optional.Int{IsSet: false},
		frecuencyPenalty:   optional.Float64{IsSet: false},
		presencePenalty:    optional.Float64{IsSet: false},
		repetitionPenalty:  optional.Float64{IsSet: false},
		minP:               optional.Float64{IsSet: false},
		topA:               optional.Float64{IsSet: false},
		seed:               optional.Int{IsSet: false},
		maxTokens:          optional.Int{IsSet: false},
		responseFormat:     optional.MapAny{IsSet: false},
		structuredOutputs:  optional.Bool{IsSet: false},
		maxPromptPrice:     optional.Float64{IsSet: false},
		maxCompletionPrice: optional.Float64{IsSet: false},
	}
}

type chatCompletionBuilder struct {
	client             *Client
	ctx                context.Context
	model              optional.String
	messages           []chatCompletionMessage
	tools              []ChatCompletionTool
	temperature        optional.Float64
	topP               optional.Float64
	topK               optional.Int
	frecuencyPenalty   optional.Float64
	presencePenalty    optional.Float64
	repetitionPenalty  optional.Float64
	minP               optional.Float64
	topA               optional.Float64
	seed               optional.Int
	maxTokens          optional.Int
	responseFormat     optional.MapAny
	structuredOutputs  optional.Bool
	maxPromptPrice     optional.Float64
	maxCompletionPrice optional.Float64
}

type chatCompletionMessage struct {
	Role    chatCompletionRole `json:"role"`    // Who the message is from.
	Content string             `json:"content"` // The content of the message
}

type ChatCompletionTool struct {
	Description string         `json:"description,omitempty,omitzero"`
	Name        string         `json:"name"`
	Parameters  map[string]any `json:"parameters"`
}

// ChatCompletionResponse is the response from the OpenRouter API for a chat completion request.
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			// Who the message is from. Must be one of openroutergo.RoleSystem, openroutergo.RoleUser, or openroutergo.RoleAssistant.
			Role chatCompletionRole `json:"role"`
			// The content of the message
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// WithContext sets the context for the chat completion request.
//
// If not set, a context.Background() context will be used.
func (b *chatCompletionBuilder) WithContext(ctx context.Context) *chatCompletionBuilder {
	b.ctx = ctx
	return b
}

// WithModel sets the model for the chat completion request.
//
// If not set, the default model configured in the OpenRouter user's account will be used.
//
// You can search for models here: https://openrouter.ai/models
func (b *chatCompletionBuilder) WithModel(model string) *chatCompletionBuilder {
	b.model = optional.String{IsSet: true, Value: model}
	return b
}

// WithSystemMessage adds a system message to the chat completion request.
//
// All messages are added to the request in the same order they are added.
func (b *chatCompletionBuilder) WithSystemMessage(message string) *chatCompletionBuilder {
	b.messages = append(b.messages, chatCompletionMessage{Role: RoleSystem, Content: message})
	return b
}

// WithUserMessage adds a user message to the chat completion request.
func (b *chatCompletionBuilder) WithUserMessage(message string) *chatCompletionBuilder {
	b.messages = append(b.messages, chatCompletionMessage{Role: RoleUser, Content: message})
	return b
}

// WithAssistantMessage adds an assistant message to the chat completion request.
func (b *chatCompletionBuilder) WithAssistantMessage(message string) *chatCompletionBuilder {
	b.messages = append(b.messages, chatCompletionMessage{Role: RoleAssistant, Content: message})
	return b
}

// WithTool adds a tool to the chat completion request so the model can return a tool call.
//
// See models supporting tool calling: https://openrouter.ai/models?supported_parameters=tools
func (b *chatCompletionBuilder) WithTool(tool ChatCompletionTool) *chatCompletionBuilder {
	b.tools = append(b.tools, tool)
	return b
}

// WithTemperature sets the temperature for the chat completion request.
//
// This setting influences the variety in the model’s responses. Lower values lead
// to more predictable and typical responses, while higher values encourage more
// diverse and less common responses. At 0, the model always gives the same
// response for a given input.
//
// Default: 1.0
//
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#temperature
//   - Explanation: https://youtu.be/ezgqHnWvua8
func (b *chatCompletionBuilder) WithTemperature(temperature float64) *chatCompletionBuilder {
	b.temperature = optional.Float64{IsSet: true, Value: temperature}
	return b
}

// WithTopP sets the top-p value for the chat completion request.
//
// This setting limits the model’s choices to a percentage of likely tokens: only the
// top tokens whose probabilities add up to P. A lower value makes the model’s responses
// more predictable, while the default setting allows for a full range of token choices.
// Think of it like a dynamic Top-K.
//
// Default: 1.0
//
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#top-p
//   - Explanation: https://youtu.be/wQP-im_HInk
func (b *chatCompletionBuilder) WithTopP(topP float64) *chatCompletionBuilder {
	b.topP = optional.Float64{IsSet: true, Value: topP}
	return b
}

// Execute executes the chat completion request with the configured parameters.
func (b *chatCompletionBuilder) Execute() (ChatCompletionResponse, error) {
	if len(b.messages) == 0 {
		return ChatCompletionResponse{}, ErrMessagesRequired
	}

	requestBodyMap := map[string]any{
		"messages": b.messages,
	}
	if b.model.IsSet {
		requestBodyMap["model"] = b.model.Value
	}

	requestBodyBytes, err := json.Marshal(requestBodyMap)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := b.client.newRequest(b.ctx, http.MethodPost, "/chat/completions", requestBodyBytes)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := b.client.httpClient.Do(req)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ChatCompletionResponse{}, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	var response ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}
