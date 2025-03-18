package handle

import (
	"encoding/json"
	"fmt"
	"github.com/RyoLena/Adventure/go-server/internal/models"
	"github.com/RyoLena/Adventure/go-server/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ChatHandler struct {
	ChatService *service.ChatService
}

func NewChatHandler(s *service.ChatService) *ChatHandler {
	return &ChatHandler{
		ChatService: s,
	}
}

// ChatInput 定义 POST /chat 接口接收的请求体

func (h *ChatHandler) Chat(c *gin.Context) {
	//1.解析请求
	var input models.ChatInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		// 请求体格式错误，返回 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//2.调用ChatService 获取AI回复
	response, err := h.ChatService.GetAIResponse(input.Message, input.ConversationID)
	if err != nil {
		// 调用 ChatService 失败，返回 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get AI response: " + err.Error(),
		})
		return
	}
	// 3. 返回 AI 回复
	c.JSON(http.StatusOK, models.ChatResponse{
		Reply:          response.Reply,
		ImageURL:       response.ImageURL,
		ConversationID: response.ConversationID,
	}) // 返回 200 OK 和 AI 回复
}

type HealthResponse struct {
	Status string `json:"status"`
}

func HealthCheck(c *gin.Context) {
	pythonServiceURL := "http://127.0.0.1:9001/health"

	// 发起 HTTP GET 请求
	resp, err := http.Get(pythonServiceURL)
	if err != nil {
		log.Printf("Failed to call Python service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to call Python service: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("Python service returned status code: %d", resp.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Python service returned status code: %d", resp.StatusCode),
		})
		return
	}

	// 解析 JSON 响应
	var healthResponse HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResponse); err != nil {
		log.Printf("Failed to decode Python service response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to decode Python service response: %v", err),
		})
		return
	}

	// 返回 Python 服务的健康状态
	c.JSON(http.StatusOK, gin.H{
		"python_service_status": healthResponse.Status, // 返回 Python 服务的状态
		"message":               "OK",                  // Go 服务的状态
	})
}
