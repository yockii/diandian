package operation

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"diandian/background/constant"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// ClickGenerator 点击操作生成器
type ClickGenerator struct {
	*BaseGenerator
}

// NewClickGenerator 创建点击操作生成器
func NewClickGenerator() *ClickGenerator {
	return &ClickGenerator{
		BaseGenerator: NewBaseGenerator(),
	}
}

// Generate 生成点击操作
func (g *ClickGenerator) Generate(contextInfo string, screenAnalysis *VisualAnalysisResponse) (*ClickOperation, error) {
	client, model, err := g.createTextClient()
	if err != nil {
		slog.Error("创建文本模型客户端失败", "error", err)
		return nil, err
	}

	// 构建消息
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: constant.PromptGenerateClickOperation,
		},
	}

	// 构建用户消息（只包含上下文和屏幕分析结果）
	userMessage := fmt.Sprintf("上下文：%s", contextInfo)

	// 如果有屏幕分析结果，添加到用户消息中
	if screenAnalysis != nil {
		analysisText := fmt.Sprintf("\n\n屏幕分析结果：找到 %d 个可交互元素", len(screenAnalysis.ElementsFound))
		for _, element := range screenAnalysis.ElementsFound {
			if element.Clickable {
				analysisText += fmt.Sprintf("\n- %s: (%d,%d) %dx%d",
					element.Description, element.Coordinates.X, element.Coordinates.Y,
					element.Coordinates.Width, element.Coordinates.Height)
			}
		}
		userMessage += analysisText
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})

	// 使用重试机制调用LLM
	content, err := g.retryLLMCall(
		func() (string, error) {
			return g.callLLM(client, model, messages)
		},
		func(content string) error {
			return g.validateClickOperation(content)
		},
		3,
		"点击操作生成",
	)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var result ClickOperation
	schema, err := jsonschema.GenerateSchemaForType(result)
	if err != nil {
		return nil, fmt.Errorf("生成点击操作schema失败: %v", err)
	}

	if err := schema.Unmarshal(content, &result); err != nil {
		slog.Error("点击操作JSON解析失败",
			"error", err,
			"raw_content", content,
			"content_length", len(content))
		return nil, fmt.Errorf("JSON schema验证失败: %v", err)
	}

	slog.Info("点击操作验证成功", "x", result.X, "y", result.Y, "button", result.Button)
	return &result, nil
}

// callLLM 调用LLM
func (g *ClickGenerator) callLLM(client *openai.Client, model string, messages []openai.ChatCompletionMessage) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("LLM返回空响应")
	}

	return resp.Choices[0].Message.Content, nil
}

// validateClickOperation 验证点击操作
func (g *ClickGenerator) validateClickOperation(content string) error {
	var tempResult ClickOperation
	schema, err := jsonschema.GenerateSchemaForType(tempResult)
	if err != nil {
		return fmt.Errorf("生成schema失败: %v", err)
	}

	if err := schema.Unmarshal(content, &tempResult); err != nil {
		slog.Error("点击操作JSON解析失败",
			"error", err,
			"raw_content", content,
			"content_length", len(content))
		return fmt.Errorf("JSON schema验证失败: %v", err)
	}

	// 业务逻辑验证
	if tempResult.X <= 0 || tempResult.Y <= 0 {
		slog.Error("点击操作验证失败：坐标无效", "x", tempResult.X, "y", tempResult.Y)
		return fmt.Errorf("坐标必须为正数")
	}

	if tempResult.Button != "left" && tempResult.Button != "right" && tempResult.Button != "middle" {
		slog.Error("点击操作验证失败：按钮类型无效", "button", tempResult.Button)
		return fmt.Errorf("按钮类型必须是 left、right 或 middle")
	}

	return nil
}
