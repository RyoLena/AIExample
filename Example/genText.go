package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"google.golang.org/api/option"
)

const modelName = "gemini-2.0-flash"
const defaultPort = "9000"

type genaiServer struct {
	ctx   context.Context
	model *genai.GenerativeModel
}

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("错误: 无法加载 .env 文件")
	}
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		fmt.Println("警告: GOOGLE_API_KEY 环境变量未设置")
		return
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("创建Gemini client 失败: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	server := &genaiServer{
		ctx:   ctx,
		model: model,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /chat", server.chatHandler)
	mux.HandleFunc("POST /stream", server.streamingChatHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},                       // 允许所有来源
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},                      // 允许的 HTTP 方法
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "Content-Type"}, // 允许的请求头
		AllowCredentials: true,                                                    // 允许携带凭证（如果需要）
	})
	handler := c.Handler(mux)

	// 访问服务器必须监听的首选端口（如果没有提供环境变量，就用defaultPort）。
	port := cmp.Or(os.Getenv("PORT"), defaultPort)
	addr := "localhost:" + port
	log.Println("Listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

// part是一段模型内容或用户查询。它只能容纳文本片段。
//JSON编码的历史数组中的一个项目基于它所代表的角色（用户/模型），
//将单个模型响应/用户查询作为一个有序的文本块数组。
//这个数组中的每一项都必须符合部分.

type part struct {
	// 一段模型内容或用户查询。
	Text string
}

// content 是传入的 JSON 编码历史数组中的每一项都必须遵守的结构。
type content struct {
	// 内容的提供者。必须标明 是"用户 "或 "模型"。
	Role string
	// 有序的`Parts`构成一条信息。
	Parts []part
}

// chatRequest 是对响应正文中传入的 JSON 编码值进行解码的结构。
type chatRequest struct {
	// 用户向模型提出的查询。
	Chat string
	// 用户与模型在当前会话中的对话历史。
	History []content
}

// chatHandler 会将模型的完整响应返回给客户端。
// 有效的JSON请求的格式如下：
//
//	  Request：
//			chat: string
//			history：[]
//
// 向客户端发送包含模型响应的 JSON payload，格式如下。
// Response：
//
//	text: string
func (gs *genaiServer) chatHandler(writer http.ResponseWriter, request *http.Request) {
	cr := &chatRequest{}
	if err := parseRequestJSON(request, cr); err != nil {
		fmt.Println("解析失败")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	cs := gs.startChat(cr.History)
	res, err := cs.SendMessage(gs.ctx, genai.Text(cr.Chat))
	if err != nil {
		fmt.Println("发送失败")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	resTxt, err := responseString(res)
	if err != nil {
		fmt.Println("接收的响应失败")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	renderResponseJSON(writer, map[string]string{"text": resTxt})
}

// streamingChatHandler 会持续将模型的响应流式传输到客户端。
// 预计请求中的 JSON 有效载荷格式如下：
//
//	Request：
//		chat: string,
//		history：[],
//
// 模型的部分响应包含一段文本。
func (gs *genaiServer) streamingChatHandler(writer http.ResponseWriter, request *http.Request) {
	cr := &chatRequest{}
	if err := parseRequestJSON(request, cr); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	cs := gs.startChat(cr.History)

	iter := cs.SendMessageStream(gs.ctx, genai.Text(cr.Chat))

	writer.Header().Set("Content-Type", "text/event-stream")

	for {
		res, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		resTxt, err := responseString(res)
		if err != nil {
			log.Println(err)
			break
		}
		_, err = fmt.Fprint(writer, resTxt)
		if err != nil {
			return
		}
		if f, ok := writer.(http.Flusher); ok {
			f.Flush()
		}
	}
}

// startChat 使用给定的历史记录启动与模型的聊天会话。
func (gs *genaiServer) startChat(history []content) *genai.ChatSession {
	cs := gs.model.StartChat()

	cs.History = transform(history)
	return cs
}

// transform 会将 []content 转换为模型可接受的 []*genai.Content 格式。
func transform(cs []content) []*genai.Content {
	gcs := make([]*genai.Content, len(cs))
	for i, c := range cs {
		gcs[i] = c.transform()
	}
	return gcs
}

func (c content) transform() *genai.Content {
	gc := &genai.Content{}
	gc.Role = c.Role
	ps := make([]genai.Part, len(c.Parts))

	for i, p := range c.Parts {
		ps[i] = genai.Text(p.Text)
	}
	gc.Parts = ps
	return gc
}

// responseString 将 genai.GenerateContentResponse 类型的模型响应转换为字符串。
func responseString(res *genai.GenerateContentResponse) (string, error) {
	//只用取GenerationConfig.CandidateCount第一个元素即可，默认为1
	if len(res.Candidates) > 0 {
		if cs := contentString(res.Candidates[0].Content); cs != nil {
			return *cs, nil
		}
	}
	return "", fmt.Errorf("来自Gemini模型无效的响应")
}

// contentString 将 genai.Content 转换为字符串。
// 如果输入内容中的各部分都是 文本类型，则会在它们之间用新行连接起来，形成一个字符串。
func contentString(c *genai.Content) *string {
	if c == nil || c.Parts == nil {
		fmt.Println("内容为空")
		return nil
	}
	cStrs := make([]string, len(c.Parts))
	for i, p := range c.Parts {
		if pt, ok := p.(genai.Text); ok {
			cStrs[i] = string(pt)
		} else {
			return nil
		}
	}
	cStr := strings.Join(cStrs, "\n")
	return &cStr
}
