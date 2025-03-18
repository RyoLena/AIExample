package models

// ChatInput 定义 POST /chat 接口接收的请求体
type ChatInput struct {
	Message        string `json:"message" binding:"required"`
	ConversationID string `json:"conversation_id"`
}

type ChatResponse struct {
	Reply          string `json:"reply"`
	ImageURL       string `json:"image_url"`
	ConversationID string `json:"conversation_id"`
}
