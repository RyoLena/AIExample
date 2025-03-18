package main

import (
	"fmt"
	"github.com/RyoLena/Adventure/go-server/internal/config"
	"github.com/RyoLena/Adventure/go-server/internal/handle"
	"github.com/RyoLena/Adventure/go-server/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	cfg, err := config.LoadConfig("go-server/internal/config/config.yml")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}
	fmt.Printf("Loaded configuration: %+v\n", cfg)

	router := gin.Default()
	paServiceURL := "http://127.0.0.1:9001"
	chatService := service.NewChatService(paServiceURL)
	chatHandle := handle.NewChatHandler(chatService)
	fmt.Printf("pythonURL：%s", chatService.PythonAIServiceURL)

	router.Use(cors.Default())

	//注册路由
	router.POST("/chat", chatHandle.Chat)
	router.GET("/health", handle.HealthCheck)

	//启动服务器
	port := ":9000"
	fmt.Printf("Server listening on port %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
