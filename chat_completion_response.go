package openroutergo

import (
	"encoding/json"

	"github.com/orsinium-labs/enum"
)

// chatCompletionFinishReason is an enum for the reason the model stopped generating tokens.
//
//   - https://openrouter.ai/docs/api-reference/overview#finish-reason
type chatCompletionFinishReason enum.Member[string]

// MarshalJSON implements the json.Marshaler interface for chatCompletionFinishReason.
func (cfr chatCompletionFinishReason) MarshalJSON() ([]byte, error) {
	return json.Marshal(cfr.Value)
}

// UnmarshalJSON implements the json.Unmarshaler interface for chatCompletionFinishReason.
func (cfr *chatCompletionFinishReason) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*cfr = chatCompletionFinishReason{Value: value}
	return nil
}

var (
	// FinishReasonStop is when the model hit a natural stop point or a provided stop sequence.
	FinishReasonStop = chatCompletionFinishReason{"stop"}
	// FinishReasonLength is when the maximum number of tokens specified in the request was reached.
	FinishReasonLength = chatCompletionFinishReason{"length"}
	// FinishReasonContentFilter is when content was omitted due to a flag from our content filters.
	FinishReasonContentFilter = chatCompletionFinishReason{"content_filter"}
	// FinishReasonToolCalls is when the model called a tool.
	FinishReasonToolCalls = chatCompletionFinishReason{"tool_calls"}
	// FinishReasonError is when the model returned an error.
	FinishReasonError = chatCompletionFinishReason{"error"}
)

// ChatCompletionResponse is the response from the OpenRouter API for a chat completion request.
//
//   - https://openrouter.ai/docs/api-reference/overview#responses
//   - https://platform.openai.com/docs/api-reference/chat/object
type ChatCompletionResponse struct {
	// A unique identifier for the chat completion.
	ID string `json:"id"`
	// A list of chat completion choices (the responses from the model).
	Choices []ChatCompletionResponseChoice `json:"choices"`
	// Usage statistics for the completion request.
	Usage ChatCompletionResponseUsage `json:"usage"`
	// The Unix timestamp (in seconds) of when the chat completion was created.
	Created int `json:"created"`
	// The model used for the chat completion.
	Model string `json:"model"`
	// The object type, which is always "chat.completion"
	Object string `json:"object"`
}

type ChatCompletionResponseChoice struct {
	// The reason the model stopped generating tokens. This will be `stop` if the model hit a
	// natural stop point or a provided stop sequence, `length` if the maximum number of
	// tokens specified in the request was reached, `content_filter` if content was omitted
	// due to a flag from our content filters, `tool_calls` if the model called a tool, or
	// `error` if the model returned an error.
	FinishReason chatCompletionFinishReason `json:"finish_reason"`
	// A chat completion message generated by the model.
	Message ChatCompletionResponseChoiceMessage `json:"message"`
}

type ChatCompletionResponseChoiceMessage struct {
	// Who the message is from. Must be one of openroutergo.RoleSystem, openroutergo.RoleUser, or openroutergo.RoleAssistant.
	Role chatCompletionRole `json:"role"`
	// The content of the message
	Content string `json:"content"`
	// When the model decided to call a tool
	ToolCalls []ChatCompletionResponseChoiceMessageToolCall `json:"tool_calls,omitempty,omitzero"`
}

type ChatCompletionResponseChoiceMessageToolCall struct {
	// The ID of the tool call.
	ID string `json:"id"`
	// The type of tool call. Always "function".
	Type string `json:"type"`
	// Function is the function that the model wants to call.
	Function ChatCompletionResponseChoiceMessageToolCallFunction `json:"function,omitempty,omitzero"`
}

type ChatCompletionResponseChoiceMessageToolCallFunction struct {
	// The name of the function to call.
	Name string `json:"name"`
	// The arguments to call the function with, as generated by the model in JSON
	// format. Note that the model does not always generate valid JSON, and may
	// hallucinate parameters not defined by your function schema. Validate the
	// arguments in your code before calling your function.
	//
	// You have to unmarshal the arguments to the correct type yourself.
	Arguments string `json:"arguments"`
}

type ChatCompletionResponseUsage struct {
	// The number of tokens in the prompt.
	PromptTokens int `json:"prompt_tokens"`
	// The number of tokens in the generated completion.
	CompletionTokens int `json:"completion_tokens"`
	// The total number of tokens used in the request (prompt + completion).
	TotalTokens int `json:"total_tokens"`
}
