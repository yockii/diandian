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

// TypeGenerator 输入操作生成器
type TypeGenerator struct {
	*BaseGenerator
}

// NewTypeGenerator 创建输入操作生成器
func NewTypeGenerator() *TypeGenerator {
	return &TypeGenerator{
		BaseGenerator: NewBaseGenerator(),
	}
}

// Generate 生成输入操作
func (g *TypeGenerator) Generate(contextInfo string) (*TypeOperation, error) {
	client, model, err := g.createTextClient()
	if err != nil {
		slog.Error("创建文本模型客户端失败", "error", err)
		return nil, err
	}

	// 构建消息
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: constant.PromptGenerateTypeOperation,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("上下文：%s", contextInfo),
		},
	}

	// 使用重试机制调用LLM
	content, err := g.retryLLMCall(
		func() (string, error) {
			return g.callLLM(client, model, messages)
		},
		func(content string) error {
			return g.validateTypeOperation(content)
		},
		3,
		"输入操作生成",
	)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var result TypeOperation
	schema, err := jsonschema.GenerateSchemaForType(result)
	if err != nil {
		return nil, fmt.Errorf("生成输入操作schema失败: %v", err)
	}

	if err := schema.Unmarshal(content, &result); err != nil {
		slog.Error("输入操作JSON解析失败",
			"error", err,
			"raw_content", content,
			"content_length", len(content))
		return nil, fmt.Errorf("JSON schema验证失败: %v", err)
	}

	slog.Info("输入操作验证成功", "text_length", len(result.Text), "text_preview", result.Text)
	return &result, nil
}

// callLLM 调用LLM
func (g *TypeGenerator) callLLM(client *openai.Client, model string, messages []openai.ChatCompletionMessage) (string, error) {
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

// validateTypeOperation 验证输入操作
func (g *TypeGenerator) validateTypeOperation(content string) error {
	var tempResult TypeOperation
	schema, err := jsonschema.GenerateSchemaForType(tempResult)
	if err != nil {
		return fmt.Errorf("生成schema失败: %v", err)
	}

	if err := schema.Unmarshal(content, &tempResult); err != nil {
		slog.Error("输入操作JSON解析失败",
			"error", err,
			"raw_content", content,
			"content_length", len(content))
		return fmt.Errorf("JSON schema验证失败: %v", err)
	}

	// 业务逻辑验证
	if tempResult.Text == "" {
		slog.Error("输入操作验证失败：文本内容为空")
		return fmt.Errorf("文本内容不能为空")
	}

	return nil
}
