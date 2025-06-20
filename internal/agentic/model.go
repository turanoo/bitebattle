package agentic

type AgenticRequest struct {
	Command string `json:"command" binding:"required"`
}

type AgenticResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ParsedPrompt struct {
	Food     string
	Location string
	Radius   string
}
