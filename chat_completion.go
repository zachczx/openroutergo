package openroutergo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/eduardolat/openroutergo/internal/debug"
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
		debug:              false,
		ctx:                context.Background(),
		model:              optional.String{IsSet: false},
		fallbackModels:     []string{},
		messages:           []chatCompletionMessage{},
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
		logitBias:          optional.MapIntInt{IsSet: false},
		logprobs:           optional.Bool{IsSet: false},
		topLogprobs:        optional.Int{IsSet: false},
		responseFormat:     optional.MapStringAny{IsSet: false},
		structuredOutputs:  optional.Bool{IsSet: false},
		stop:               []string{},
		tools:              []chatCompletionToolFunction{},
		toolChoice:         optional.String{IsSet: false},
		maxPromptPrice:     optional.Float64{IsSet: false},
		maxCompletionPrice: optional.Float64{IsSet: false},
	}
}

type chatCompletionBuilder struct {
	client             *Client
	debug              bool
	ctx                context.Context
	model              optional.String
	fallbackModels     []string
	messages           []chatCompletionMessage
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
	logitBias          optional.MapIntInt
	logprobs           optional.Bool
	topLogprobs        optional.Int
	responseFormat     optional.MapStringAny
	structuredOutputs  optional.Bool
	stop               []string
	tools              []chatCompletionToolFunction
	toolChoice         optional.String
	maxPromptPrice     optional.Float64
	maxCompletionPrice optional.Float64
}

// Clone returns a completely new chat completion builder with the same configuration as the current
// builder.
//
// This is useful if you want to reuse the same configuration for multiple requests.
func (b *chatCompletionBuilder) Clone() *chatCompletionBuilder {
	return &chatCompletionBuilder{
		client:             b.client,
		debug:              b.debug,
		ctx:                b.ctx,
		messages:           b.messages,
		model:              b.model,
		fallbackModels:     b.fallbackModels,
		temperature:        b.temperature,
		topP:               b.topP,
		topK:               b.topK,
		frecuencyPenalty:   b.frecuencyPenalty,
		presencePenalty:    b.presencePenalty,
		repetitionPenalty:  b.repetitionPenalty,
		minP:               b.minP,
		topA:               b.topA,
		seed:               b.seed,
		maxTokens:          b.maxTokens,
		logitBias:          b.logitBias,
		logprobs:           b.logprobs,
		topLogprobs:        b.topLogprobs,
		responseFormat:     b.responseFormat,
		structuredOutputs:  b.structuredOutputs,
		stop:               b.stop,
		tools:              b.tools,
		toolChoice:         b.toolChoice,
		maxPromptPrice:     b.maxPromptPrice,
		maxCompletionPrice: b.maxCompletionPrice,
	}
}

type chatCompletionMessage struct {
	Role    chatCompletionRole `json:"role"`    // Who the message is from.
	Content string             `json:"content"` // The content of the message
}

type chatCompletionToolFunction struct {
	Type     string             `json:"type"` // Always "function"
	Function ChatCompletionTool `json:"function"`
}

// ChatCompletionTool is a tool that can be used in a chat completion request.
//
//   - Models supporting tool calling: https://openrouter.ai/models?supported_parameters=tools
//   - JSON Schema reference: https://json-schema.org/understanding-json-schema/reference
//   - Tool calling example: https://platform.openai.com/docs/guides/function-calling
type ChatCompletionTool struct {
	// The name of the tool, when the model calls this tool, it will return this name so
	// you can identify it.
	Name string `json:"name"`
	// The description of the tool, make sure to give a good description so the model knows
	// when to use it.
	Description string `json:"description,omitempty,omitzero"`
	// Make sure to define your tool's parameters using map[string]any and following the
	// JSON Schema format.
	//
	//   - Format example: https://platform.openai.com/docs/guides/function-calling
	//   - JSON Schema reference: https://json-schema.org/understanding-json-schema/reference
	Parameters map[string]any `json:"parameters"`
}

// WithDebug sets the debug flag for the chat completion request.
//
// If true, the JSON request and response will be printed to the console for debugging purposes.
func (b *chatCompletionBuilder) WithDebug(debug bool) *chatCompletionBuilder {
	b.debug = debug
	return b
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

// WithModelFallback adds a model to the fallback list for the chat completion request.
//
// You can call this method up to 3 times to add more than one fallback model.
//
// This lets you automatically try other models if the primary model’s providers are down,
// rate-limited, or refuse to reply due to content moderation.
//
// If the primary model is not available, all the fallback models will be tried in the
// same order they were added.
//
//   - Docs: https://openrouter.ai/docs/features/model-routing#the-models-parameter
//   - Example: https://openrouter.ai/docs/features/model-routing#using-with-openai-sdk
func (b *chatCompletionBuilder) WithModelFallback(modelFallback string) *chatCompletionBuilder {
	b.fallbackModels = append(b.fallbackModels, modelFallback)
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

// WithLogitBias Accepts a JSON object that maps tokens (specified by their token ID in the tokenizer) to
// an associated bias value from -100 to 100. Mathematically, the bias is added to the logits generated
// by the model prior to sampling. The exact effect will vary per model, but values between -1 and 1 should
// decrease or increase likelihood of selection; values like -100 or 100 should result in a ban or
// exclusive selection of the relevant token.
//
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#logit-bias
func (b *chatCompletionBuilder) WithLogitBias(logitBias map[int]int) *chatCompletionBuilder {
	b.logitBias = optional.MapIntInt{IsSet: true, Value: logitBias}
	return b
}

// WithLogprobs Whether to return log probabilities of the output tokens or not. If true, returns the
// log probabilities of each output token returned.
//
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#logprobs
func (b *chatCompletionBuilder) WithLogprobs(logprobs bool) *chatCompletionBuilder {
	b.logprobs = optional.Bool{IsSet: true, Value: logprobs}
	return b
}

// WithTopLogprobs An integer between 0 and 20 specifying the number of most likely tokens to return
// at each token position, each with an associated log probability. logprobs must be set to true if
// this parameter is used.
//
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#top-logprobs
func (b *chatCompletionBuilder) WithTopLogprobs(topLogprobs int) *chatCompletionBuilder {
	b.topLogprobs = optional.Int{IsSet: true, Value: topLogprobs}
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
	b.responseFormat = optional.MapStringAny{IsSet: true, Value: responseFormat}
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

// WithStop Stop generation immediately if the model encounter any token specified in the stop array.
//
//   - Docs: https://openrouter.ai/docs/api-reference/parameters#stop
func (b *chatCompletionBuilder) WithStop(stop []string) *chatCompletionBuilder {
	b.stop = stop
	return b
}

// WithTool adds a tool to the chat completion request so the model can return a tool call.
//
// If your tool requires parameters, read the [ChatCompletionTool] type documentation
// for more information on how to define the parameters using JSON Schema.
//
//   - Models supporting tool calling: https://openrouter.ai/models?supported_parameters=tools
//   - JSON Schema reference: https://json-schema.org/understanding-json-schema/reference
//   - Tool calling example: https://platform.openai.com/docs/guides/function-calling
func (b *chatCompletionBuilder) WithTool(tool ChatCompletionTool) *chatCompletionBuilder {
	b.tools = append(b.tools, chatCompletionToolFunction{Type: "function", Function: tool})
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

// errorResponse is a struct that represents an error response when there is an error
// in the response from the OpenRouter API.
//
//   - Docs: https://openrouter.ai/docs/api-reference/errors
type errorResponse struct {
	Error struct {
		Code     int            `json:"code"`
		Message  string         `json:"message"`
		Metadata map[string]any `json:"metadata"`
	} `json:"error"`
}

// Execute the chat completion request with the configured parameters.
//
// Returns:
//
//   - The chat completion builder in the same state as before calling this method.
//   - The response from the OpenRouter API.
//   - An error if the request fails.
//
// IMPORTANT: The first return value (the builder) does not include the new assistant message content.
// To continue the conversation with the assistant's response, you must explicitly add it using
// the [WithAssistantMessage] method.
//
// Example:
//
//	completion := client.
//		NewChatCompletion().
//		WithModel("...").
//		WithSystemMessage("You are a helpful assistant expert in geography.").
//		WithUserMessage("What is the capital of France?")
//
//	completion, resp, err := completion.Execute()
//	if err != nil {
//		// handle error
//	}
//
//	// Use the response, add the response to the builder to continue the conversation
//	completion = completion.WithAssistantMessage(
//		resp.Choices[0].Message.Content,
//	)
//
//	// Use the same builder for another request
//	completion = completion.WithUserMessage("Thank you!! Now, what is the capital of Germany?")
//	_, resp, err = completion.Execute()
//	if err != nil {
//		// handle error
//	}
func (b *chatCompletionBuilder) Execute() (*chatCompletionBuilder, ChatCompletionResponse, error) {
	clone := b.Clone()

	if len(b.messages) == 0 {
		return clone, ChatCompletionResponse{}, ErrMessagesRequired
	}

	requestBodyMap := map[string]any{}
	if len(b.messages) > 0 {
		requestBodyMap["messages"] = b.messages
	}
	if b.model.IsSet {
		requestBodyMap["model"] = b.model.Value
	}
	if len(b.fallbackModels) > 0 {
		requestBodyMap["models"] = b.fallbackModels
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
	if b.logitBias.IsSet {
		requestBodyMap["logit_bias"] = b.logitBias.Value
	}
	if b.logprobs.IsSet {
		requestBodyMap["logprobs"] = b.logprobs.Value
	}
	if b.topLogprobs.IsSet {
		requestBodyMap["top_logprobs"] = b.topLogprobs.Value
	}
	if b.responseFormat.IsSet {
		requestBodyMap["response_format"] = b.responseFormat.Value
	}
	if b.structuredOutputs.IsSet {
		requestBodyMap["structured_outputs"] = b.structuredOutputs.Value
	}
	if len(b.stop) > 0 {
		requestBodyMap["stop"] = b.stop
	}
	if len(b.tools) > 0 {
		requestBodyMap["tools"] = b.tools
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

	if b.debug {
		fmt.Println()
		fmt.Println("---------------------------")
		fmt.Println("-- Request to OpenRouter --")
		fmt.Println("---------------------------")
		debug.PrintAsJSON(requestBodyMap)
		fmt.Println()
	}

	requestBodyBytes, err := json.Marshal(requestBodyMap)
	if err != nil {
		return clone, ChatCompletionResponse{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := b.client.newRequest(b.ctx, http.MethodPost, "/chat/completions", requestBodyBytes)
	if err != nil {
		return clone, ChatCompletionResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := b.client.httpClient.Do(req)
	if err != nil {
		return clone, ChatCompletionResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return clone, ChatCompletionResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var tempResp map[string]any
	if err := json.Unmarshal(bodyBytes, &tempResp); err != nil {
		return clone, ChatCompletionResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if b.debug {
		fmt.Println()
		fmt.Println("------------------------------")
		fmt.Println("-- Response from OpenRouter --")
		fmt.Println("------------------------------")
		fmt.Printf("Status code: %d\n", resp.StatusCode)
		debug.PrintAsJSON(tempResp)
		fmt.Println()
	}

	if tempResp["error"] != nil {
		var errorResponse errorResponse
		if err := json.Unmarshal(bodyBytes, &errorResponse); err != nil {
			return clone, ChatCompletionResponse{}, fmt.Errorf("failed to decode error response: %w", err)
		}
		return clone, ChatCompletionResponse{}, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, errorResponse.Error.Message)
	}

	var response ChatCompletionResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return clone, ChatCompletionResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return clone, response, nil
}
