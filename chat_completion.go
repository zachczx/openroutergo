package openroutergo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

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
		tools:              []chatCompletionToolFunction{},
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
		toolChoice:         optional.String{IsSet: false},
		maxPromptPrice:     optional.Float64{IsSet: false},
		maxCompletionPrice: optional.Float64{IsSet: false},
	}
}

type chatCompletionBuilder struct {
	client             *Client
	ctx                context.Context
	model              optional.String
	messages           []chatCompletionMessage
	tools              []chatCompletionToolFunction
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
	toolChoice         optional.String
	maxPromptPrice     optional.Float64
	maxCompletionPrice optional.Float64
}

type chatCompletionMessage struct {
	Role    chatCompletionRole `json:"role"`    // Who the message is from.
	Content string             `json:"content"` // The content of the message
}

type chatCompletionToolFunction struct {
	Type     string             `json:"type"` // Always "function"
	Function ChatCompletionTool `json:"function"`
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
	b.tools = append(b.tools, chatCompletionToolFunction{Type: "function", Function: tool})
	return b
}

// WithTemperature sets the temperature for the chat completion request.
//
// This setting influences the variety in the model’s responses. Lower values lead
// to more predictable and typical responses, while higher values encourage more
// diverse and less common responses. At 0, the model always gives the same
// response for a given input.
//
//   - Default: 1.0
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
//   - Default: 1.0
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#top-p
//   - Explanation: https://youtu.be/wQP-im_HInk
func (b *chatCompletionBuilder) WithTopP(topP float64) *chatCompletionBuilder {
	b.topP = optional.Float64{IsSet: true, Value: topP}
	return b
}

// WithTopK sets the top-k value for the chat completion request.
//
// This limits the model's choice of tokens at each step, making it choose from
// a smaller set. A value of 1 means the model will always pick the most likely
// next token, leading to predictable results. By default this setting is disabled,
// making the model to consider all choices.
//
//   - Default: 0
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#top-k
//   - Explanation: https://youtu.be/EbZv6-N8Xlk
func (b *chatCompletionBuilder) WithTopK(topK int) *chatCompletionBuilder {
	b.topK = optional.Int{IsSet: true, Value: topK}
	return b
}

// WithFrequencyPenalty sets the frequency penalty for the chat completion request.
//
// This setting aims to control the repetition of tokens based on how often they appear
// in the input. It tries to use less frequently those tokens that appear more in the
// input, proportional to how frequently they occur. Token penalty scales with the number
// of occurrences. Negative values will encourage token reuse.
//
//   - Default: 0.0
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#frequency-penalty
//   - Explanation: https://youtu.be/p4gl6fqI0_w
func (b *chatCompletionBuilder) WithFrequencyPenalty(frequencyPenalty float64) *chatCompletionBuilder {
	b.frecuencyPenalty = optional.Float64{IsSet: true, Value: frequencyPenalty}
	return b
}

// WithPresencePenalty sets the presence penalty for the chat completion request.
//
// Adjusts how often the model repeats specific tokens already used in the input.
// Higher values make such repetition less likely, while negative values do the opposite.
// Token penalty does not scale with the number of occurrences. Negative values will
// encourage token reuse.
//
//   - Default: 0.0
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#presence-penalty
//   - Explanation: https://youtu.be/MwHG5HL-P74
func (b *chatCompletionBuilder) WithPresencePenalty(presencePenalty float64) *chatCompletionBuilder {
	b.presencePenalty = optional.Float64{IsSet: true, Value: presencePenalty}
	return b
}

// WithRepetitionPenalty sets the repetition penalty for the chat completion request.
//
// Helps to reduce the repetition of tokens from the input. A higher value makes the
// model less likely to repeat tokens, but too high a value can make the output less
// coherent (often with run-on sentences that lack small words). Token penalty scales
// based on original token's probability.
//
//   - Default: 1.0
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#repetition-penalty
//   - Explanation: https://youtu.be/LHjGAnLm3DM
func (b *chatCompletionBuilder) WithRepetitionPenalty(repetitionPenalty float64) *chatCompletionBuilder {
	b.repetitionPenalty = optional.Float64{IsSet: true, Value: repetitionPenalty}
	return b
}

// WithMinP sets the min-p value for the chat completion request.
//
// Represents the minimum probability for a token to be considered, relative to
// the probability of the most likely token. If your Min-P is set to 0.1, that
// means it will only allow for tokens that are at least 1/10th as probable as
// the best possible option.
//
//   - Default: 0.0
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#min-p
func (b *chatCompletionBuilder) WithMinP(minP float64) *chatCompletionBuilder {
	b.minP = optional.Float64{IsSet: true, Value: minP}
	return b
}

// WithTopA sets the top-a value for the chat completion request.
//
// Consider only the top tokens with "sufficiently high" probabilities based on
// the probability of the most likely token. Think of it like a dynamic Top-P.
// A lower Top-A value focuses the choices based on the highest probability token
// but with a narrower scope.
//
//   - Default: 0.0
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#top-a
func (b *chatCompletionBuilder) WithTopA(topA float64) *chatCompletionBuilder {
	b.topA = optional.Float64{IsSet: true, Value: topA}
	return b
}

// WithSeed sets the seed value for the chat completion request.
//
// If specified, the inferencing will sample deterministically, such that repeated
// requests with the same seed and parameters should return the same result.
// Determinism is not guaranteed for some models.
//
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#seed
func (b *chatCompletionBuilder) WithSeed(seed int) *chatCompletionBuilder {
	b.seed = optional.Int{IsSet: true, Value: seed}
	return b
}

// WithMaxTokens sets the maximum number of tokens to generate for the chat completion request.
//
// This sets the upper limit for the number of tokens the model can generate in response.
// It won't produce more than this limit. The maximum value is the context length minus
// the prompt length.
//
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#max-tokens
func (b *chatCompletionBuilder) WithMaxTokens(maxTokens int) *chatCompletionBuilder {
	b.maxTokens = optional.Int{IsSet: true, Value: maxTokens}
	return b
}

// WithResponseFormat sets the response format for the chat completion request.
//
// Forces the model to produce specific output format. Setting to { "type": "json_object" }
// enables JSON mode, which guarantees the message the model generates is valid JSON.
//
// Note: when using JSON mode, you should also instruct the model to produce JSON
// yourself via a system or user message.
//
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#response-format
func (b *chatCompletionBuilder) WithResponseFormat(responseFormat map[string]any) *chatCompletionBuilder {
	b.responseFormat = optional.MapAny{IsSet: true, Value: responseFormat}
	return b
}

// WithStructuredOutputs sets whether the model can return structured outputs.
//
// If the model can return structured outputs using response_format json_schema.
//
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#structured-outputs
func (b *chatCompletionBuilder) WithStructuredOutputs(structuredOutputs bool) *chatCompletionBuilder {
	b.structuredOutputs = optional.Bool{IsSet: true, Value: structuredOutputs}
	return b
}

// WithToolChoice controls which (if any) tool is called by the model.
//
//   - none: The model will not call any tool and instead generates a message.
//   - auto: The model can pick between generating a message or calling one or more tools.
//   - required: The model must call one or more tools.
//
// If you want to force the model to call a specific tool, set the toolChoice parameter
// to the name of the tool you want to call and this will send the tool in the correct
// format to the model.
//
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#tool-choice
func (b *chatCompletionBuilder) WithToolChoice(toolChoice string) *chatCompletionBuilder {
	b.toolChoice = optional.String{IsSet: true, Value: toolChoice}
	return b
}

// WithMaxPrice sets the maximum price accepted for the chat completion request for both prompt and completion tokens.
//
// For example, the value (1, 2) will route to any provider with a price of <= $1/m prompt tokens and <= $2/m completion tokens.
//
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#max-price
func (b *chatCompletionBuilder) WithMaxPrice(maxPromptPrice float64, maxCompletionPrice float64) *chatCompletionBuilder {
	b.maxPromptPrice = optional.Float64{IsSet: true, Value: maxPromptPrice}
	b.maxCompletionPrice = optional.Float64{IsSet: true, Value: maxCompletionPrice}
	return b
}

// Execute executes the chat completion request with the configured parameters.
func (b *chatCompletionBuilder) Execute() (ChatCompletionResponse, error) {
	if len(b.messages) == 0 {
		return ChatCompletionResponse{}, ErrMessagesRequired
	}

	requestBodyMap := map[string]any{}
	if len(b.messages) > 0 {
		requestBodyMap["messages"] = b.messages
	}
	if len(b.tools) > 0 {
		requestBodyMap["tools"] = b.tools
	}
	if b.model.IsSet {
		requestBodyMap["model"] = b.model.Value
	}
	if b.temperature.IsSet {
		requestBodyMap["temperature"] = b.temperature.Value
	}
	if b.topP.IsSet {
		requestBodyMap["top_p"] = b.topP.Value
	}
	if b.topK.IsSet {
		requestBodyMap["top_k"] = b.topK.Value
	}
	if b.frecuencyPenalty.IsSet {
		requestBodyMap["frequency_penalty"] = b.frecuencyPenalty.Value
	}
	if b.presencePenalty.IsSet {
		requestBodyMap["presence_penalty"] = b.presencePenalty.Value
	}
	if b.repetitionPenalty.IsSet {
		requestBodyMap["repetition_penalty"] = b.repetitionPenalty.Value
	}
	if b.minP.IsSet {
		requestBodyMap["min_p"] = b.minP.Value
	}
	if b.topA.IsSet {
		requestBodyMap["top_a"] = b.topA.Value
	}
	if b.seed.IsSet {
		requestBodyMap["seed"] = b.seed.Value
	}
	if b.maxTokens.IsSet {
		requestBodyMap["max_tokens"] = b.maxTokens.Value
	}
	if b.responseFormat.IsSet {
		requestBodyMap["response_format"] = b.responseFormat.Value
	}
	if b.structuredOutputs.IsSet {
		requestBodyMap["structured_outputs"] = b.structuredOutputs.Value
	}
	if b.toolChoice.IsSet {
		if slices.Contains([]string{"none", "auto", "required"}, b.toolChoice.Value) {
			requestBodyMap["tool_choice"] = b.toolChoice.Value
		} else {
			requestBodyMap["tool_choice"] = map[string]any{
				"type": "function",
				"function": map[string]string{
					"name": b.toolChoice.Value,
				},
			}
		}
	}
	if b.maxPromptPrice.IsSet && b.maxCompletionPrice.IsSet {
		requestBodyMap["max_price"] = map[string]float64{
			"prompt":     b.maxPromptPrice.Value,
			"completion": b.maxCompletionPrice.Value,
		}
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
