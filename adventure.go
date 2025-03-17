package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/googleapis/gax-go/v2/apierror"
	"google.golang.org/api/googleapi"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// 设置系统prompt
const systemInstructionsFile = "system-instructions.md"

var sleepTime = struct {
	character time.Duration
	sentence  time.Duration
}{
	character: time.Millisecond * 30,
	sentence:  time.Millisecond * 300,
}

// // 流输出列的位置。
var col = 0

// getBytes 返回文件内容的字节数。
func getBytes(path string) []byte {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file bytes %v: %v\n", path, err)
	}
	return bytes
}

// 从文件读取 API 密钥
func getAPIKeyFromFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// newClient 通过api-key权限新建一个客户端
func newClient(ctx context.Context) *genai.Client {
	apiKey, exists := os.LookupEnv("API_KEY")
	if !exists {
		var err error
		apiKey, err = getAPIKeyFromFile(".api_key")
		if err != nil {
			log.Fatalf("无法获取 API 密钥：环境变量中没有 API_KEY，且无法读取 .api_key 文件: %v\n"+
				"可以去 Google 的 AI Studio 获取，浏览 https://aistudio.google.com/，选择 'Get API key'。\n", err)
		}
	}

	// 确保 API 密钥不为空
	if apiKey == "" {
		log.Fatalf("API 密钥为空\n" +
			"可以去 Google 的 AI Studio 获取，浏览 https://aistudio.google.com/，选择 'Get API key'。\n")
	}

	// 使用 apiKey 进行后续操作
	log.Println("成功获取 API 密钥")

	//新建一个client，使用api key来批准
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}
	return client
}

func main() {
	ctx := context.Background()
	client := newClient(ctx)
	defer client.Close()

	//选择模型
	model := client.GenerativeModel("gemini-2.0-flash")

	//初始化新的聊天会话。
	session := model.StartChat()

	//设置系统命令
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(getBytes(systemInstructionsFile))},
		Role:  "system",
	}

	model.SetTemperature(0.7)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(1024)

	dreamQuestion := "What do want to dream about?"

	//建立聊天记录
	session.History = []*genai.Content{{
		Role:  "model",
		Parts: []genai.Part{genai.Text(dreamQuestion)},
	}}
	printRuneFormatted('\n')

	topic := askUser(dreamQuestion)
	sendAndPrintResponse(ctx, session, topic)

	chat(ctx, session)

}

// chat 简单的聊天循环
func chat(ctx context.Context, session *genai.ChatSession) {
	for {
		fmt.Println()
		userInput := askUser(">>")
		sendAndPrintResponse(ctx, session, userInput)
	}
}

// askUser 提示用户输入.
func askUser(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		printStringFormatted(fmt.Sprintf("%v", prompt))
		action, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input:%v\n", err)
		}
		action = strings.TrimSpace(action)
		if len(action) == 0 {
			continue
		}
		return action
	}
}

// sendAndPrintResponse 向模型发送信息并打印响应。
func sendAndPrintResponse(ctx context.Context, session *genai.ChatSession, text string) {
	it := session.SendMessageStream(ctx, genai.Text(text))
	printRuneFormatted('\n')
	printRuneFormatted('\n')

	for {
		resp, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Printf("错误类型: %T", err)
			log.Printf("完整错误: %+v", err)
			printStringFormatted("\n\nYou feel a jolt of electricity as you realize you're being unplugged from the matrix.\n\n")
			log.Printf("Error sending message: err=%v\n", err)

			var aerr *googleapi.Error
			if errors.As(err, &aerr) {
				log.Printf("Google API 错误码: %d", aerr.Code)
				log.Printf("Google API 错误信息: %s", aerr.Message)
				log.Printf("Google API 错误详情: %+v", aerr.Errors)
			}

			var ae *apierror.APIError
			if errors.As(err, &ae) {
				log.Printf("ae.Reason(): %v\n", ae.Reason())
				log.Printf("ae.Details().Help.GetLinks(): %v\n", ae.Details().Help.GetLinks())
			}
			if s, ok := status.FromError(err); ok {
				log.Printf("s.Message: %v\n", s.Message())
				for _, d := range s.Proto().Details {
					log.Printf("- Details: %v\n", d)
				}
			}
			os.Exit(1)
		}
		for _, cand := range resp.Candidates {
			streamPartialResponse(cand.Content.Parts)
		}
	}
	printRuneFormatted('\n')
}

// streamPartialResponse 打印部分回复。
func streamPartialResponse(parts []genai.Part) {
	for _, part := range parts {
		printStringFormatted(fmt.Sprintf("\n%v", part))
	}
}

// printStringFormatted 打印字符串并将其格式化，同时为达到效果而延迟.
func printStringFormatted(text string) {
	for _, c := range text {
		printRuneFormatted(c)
	}
}

// printRuneFormatted 打印符文并将其格式化，效果延迟.
func printRuneFormatted(c rune) {
	switch c {
	case '.':
		fmt.Print(string(c))
		col++
		time.Sleep(sleepTime.sentence)
	case '\n':
		fmt.Print(string(c))
		col++
	case ' ':
		if col == 0 {
			//还没开始不用管
		} else if col > 80 {
			//换行
			fmt.Print("\n")
			col = 0
		} else {
			fmt.Print(string(c))
			col++
		}
	default:
		fmt.Print(string(c))
		col++
	}
	time.Sleep(sleepTime.character)
}
