package agentic

type AgenticCommandRequest struct {
	Command string `json:"command" binding:"required"`
}

type AgenticCommandResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ParsedPrompt struct {
	Food     string
	Location string
	Radius   string
}
