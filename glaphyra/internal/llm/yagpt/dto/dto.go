package dto

// RequestDTO запрос к YaGPT API
type RequestDTO struct {
	Model         string  `json:"model"`
	Temperature   float32 `json:"temperature"`
	MaxTokens     int32   `json:"max_tokens"`
	SystemMessage string  `json:"system_message"`
	UserMessage   string  `json:"user_message"`
}

// ResponseDTO ответ от YaGPT API
type ResponseDTO struct {
	Result           string `json:"result"`
	InputTextTokens  string `json:"input_text_tokens"`
	CompletionTokens string `json:"completion_tokens"`
	TotalTokens      string `json:"total_tokens"`
	ModelVersion     string `json:"model_version"`
}
