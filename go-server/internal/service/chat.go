package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RyoLena/Adventure/go-server/internal/models"
	"io"
	"net/http"
	"time"
)

// ChatService 结构体，包含python提供服务的URL和HTTP客户端
type ChatService struct {
	PythonAIServiceURL string
	httpClient         *http.Client //使用 HTTP 客户端，可以配置超时等参数
}

// NewChatService 创建一个新的chatService实例
func NewChatService(paServiceURL string) *ChatService {
	return &ChatService{
		PythonAIServiceURL: paServiceURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ChatRequest 定义发送给 Python AI 服务的请求结构体
type ChatRequest struct {
	Message        string `json:"message"`
	ConversationID string `json:"conversation_id"` // 可选的对话 ID
}

// ChatResponse 定义从 Python AI 服务接收的响应结构体
//type ChatResponse struct {
//	Reply          string `json:"reply"`
//	ImageURL       string `json:"image_url"`
//	ConversationID string `json:"conversation_id"`
//}

// GetAIResponse 发送消息给 Python AI 服务，并接收回复
func (s *ChatService) GetAIResponse(message string, conversationID string) (models.ChatResponse, error) {
	// 1. 构造请求体
	requestBody, err := json.Marshal(ChatRequest{
		Message:        message,
		ConversationID: conversationID,
	})
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 2. 创建 HTTP 请求
	req, err := http.NewRequest("POST", s.PythonAIServiceURL+"/chat", bytes.NewBuffer(requestBody))
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json") // 设置 Content-Type

	// 3. 发送 HTTP 请求
	resp, err := s.httpClient.Do(req) // 使用配置好的 HTTP 客户端\
	fmt.Println("python服务的URL位置", s.PythonAIServiceURL)
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// 4. 处理 HTTP 响应
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // 读取错误信息，方便调试
		return models.ChatResponse{}, fmt.Errorf("received non-OK status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// 5. 解析 HTTP 响应
	var response models.ChatResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return models.ChatResponse{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return response, nil
}
