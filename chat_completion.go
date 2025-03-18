package openroutergo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
		client:   c,
		ctx:      context.Background(),
		model:    "gpt-4o-mini",
		messages: []chatCompletionMessage{},
	}
}

type chatCompletionBuilder struct {
	client   *Client
	ctx      context.Context
	model    string
	messages []chatCompletionMessage
}

type chatCompletionMessage struct {
	Role    chatCompletionRole `json:"role"`    // Who the message is from.
	Content string             `json:"content"` // The content of the message
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
// You can search for models here: https://openrouter.ai/models
func (b *chatCompletionBuilder) WithModel(model string) *chatCompletionBuilder {
	b.model = model
	return b
}

// AddSystemMessage adds a system message to the chat completion request.
//
// All messages are added to the request in the same order they are added.
func (b *chatCompletionBuilder) AddSystemMessage(message string) *chatCompletionBuilder {
	b.messages = append(b.messages, chatCompletionMessage{Role: RoleSystem, Content: message})
	return b
}

// AddUserMessage adds a user message to the chat completion request.
func (b *chatCompletionBuilder) AddUserMessage(message string) *chatCompletionBuilder {
	b.messages = append(b.messages, chatCompletionMessage{Role: RoleUser, Content: message})
	return b
}

// AddAssistantMessage adds an assistant message to the chat completion request.
func (b *chatCompletionBuilder) AddAssistantMessage(message string) *chatCompletionBuilder {
	b.messages = append(b.messages, chatCompletionMessage{Role: RoleAssistant, Content: message})
	return b
}

// Execute executes the chat completion request with the configured parameters.
func (b *chatCompletionBuilder) Execute() (ChatCompletionResponse, error) {
	if b.model == "" {
		return ChatCompletionResponse{}, ErrModelRequired
	}

	if len(b.messages) == 0 {
		return ChatCompletionResponse{}, ErrMessagesRequired
	}

	requestBody, err := json.Marshal(struct {
		Model    string                  `json:"model"`
		Messages []chatCompletionMessage `json:"messages"`
	}{
		Model:    b.model,
		Messages: b.messages,
	})
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := b.client.newRequest(b.ctx, http.MethodPost, "/chat/completions", requestBody)
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
