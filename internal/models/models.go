package models

// ParseRequest 表示URL解析请求
type ParseRequest struct {
	URL string `json:"url" binding:"required"`
}

// ParseResponse 表示URL解析响应
type ParseResponse struct {
	Success bool   `json:"success"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	URL     string `json:"url,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ChatRequest 表示聊天请求
type ChatRequest struct {
	URL            string `json:"url"`
	Message        string `json:"message" binding:"required"`
	Model          string `json:"model,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
}

// ChatResponse 表示聊天响应
type ChatResponse struct {
	Success        bool   `json:"success"`
	Response       string `json:"response,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	Model          string `json:"model,omitempty"`
	Error          string `json:"error,omitempty"`
}

// ErrorResponse 表示API错误响应
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
