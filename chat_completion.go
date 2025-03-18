package openroutergo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NewChatCompletion creates a new chat completion request builder for the OpenRouter API.
//
// Docs:
//   - Reference: https://openrouter.ai/docs/api-reference/chat-completion
//   - Request: http://openrouter.ai/docs/api-reference/overview#completions-request-format
//   - Parameters: https://openrouter.ai/docs/api-reference/parameters
//   - Response: https://openrouter.ai/docs/api-reference/overview#completionsresponse-format
func (c *Client) NewChatCompletion() *chatCompletionBuilder {
	return &chatCompletionBuilder{
		client:   c,
		ctx:      context.Background(),
		model:    "gpt-4o-mini",
		messages: []chatCompletionMessage{},
		timeout:  c.httpClient.Timeout,
	}
}

type chatCompletionBuilder struct {
	client   *Client
	ctx      context.Context
	model    string
	timeout  time.Duration
	messages []chatCompletionMessage
}

type chatCompletionMessage struct {
	Role    string `json:"role"`    // Who the message is from. Allowed values: system, user, assistant
	Content string `json:"content"` // The content of the message
}

// ChatCompletionResponse is the response from the OpenRouter API for a chat completion request.
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`    // Who the message is from. Allowed values: system, user, assistant
			Content string `json:"content"` // The content of the message
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

// WithTimeout sets the timeout for the chat completion request.
//
// If not set, the default timeout set in the client will be used.
func (b *chatCompletionBuilder) WithTimeout(timeout time.Duration) *chatCompletionBuilder {
	b.timeout = timeout
	return b
}

// AddSystemMessage adds a system message to the chat completion request.
//
// All messages are added to the request in the same order they are added.
func (b *chatCompletionBuilder) AddSystemMessage(message string) *chatCompletionBuilder {
	b.messages = append(b.messages, chatCompletionMessage{Role: "system", Content: message})
	return b
}

// AddUserMessage adds a user message to the chat completion request.
func (b *chatCompletionBuilder) AddUserMessage(message string) *chatCompletionBuilder {
	b.messages = append(b.messages, chatCompletionMessage{Role: "user", Content: message})
	return b
}

// AddAssistantMessage adds an assistant message to the chat completion request.
func (b *chatCompletionBuilder) AddAssistantMessage(message string) *chatCompletionBuilder {
	b.messages = append(b.messages, chatCompletionMessage{Role: "assistant", Content: message})
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

	req, err := http.NewRequestWithContext(b.ctx, http.MethodPost, b.client.baseURL+"/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+b.client.apiKey)

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
