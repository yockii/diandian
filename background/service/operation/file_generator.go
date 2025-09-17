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

// FileGenerator 文件操作生成器
type FileGenerator struct {
	*BaseGenerator
}

// NewFileGenerator 创建文件操作生成器
func NewFileGenerator() *FileGenerator {
	return &FileGenerator{
		BaseGenerator: NewBaseGenerator(),
	}
}

// Generate 生成文件操作
func (g *FileGenerator) Generate(contextInfo string) (*FileOperation, error) {
	client, model, err := g.createTextClient()
	if err != nil {
		slog.Error("创建文本模型客户端失败", "error", err)
		return nil, err
	}

	// 构建消息
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: constant.PromptGenerateFileOperation,
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
			return g.validateFileOperation(content)
		},
		3,
		"文件操作生成",
	)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var result FileOperation
	schema, err := jsonschema.GenerateSchemaForType(result)
	if err != nil {
		return nil, fmt.Errorf("生成文件操作schema失败: %v", err)
	}

	if err := schema.Unmarshal(content, &result); err != nil {
		slog.Error("文件操作JSON解析失败",
			"error", err,
			"raw_content", content,
			"content_length", len(content))
		return nil, fmt.Errorf("JSON schema验证失败: %v", err)
	}

	slog.Info("文件操作验证成功",
		"operation", result.Operation,
		"source_path", result.SourcePath,
		"target_path", result.TargetPath,
		"has_content", len(result.Content) > 0)
	return &result, nil
}

// callLLM 调用LLM
func (g *FileGenerator) callLLM(client *openai.Client, model string, messages []openai.ChatCompletionMessage) (string, error) {
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

// validateFileOperation 验证文件操作
func (g *FileGenerator) validateFileOperation(content string) error {
	var tempResult FileOperation
	schema, err := jsonschema.GenerateSchemaForType(tempResult)
	if err != nil {
		return fmt.Errorf("生成schema失败: %v", err)
	}

	if err := schema.Unmarshal(content, &tempResult); err != nil {
		slog.Error("文件操作JSON解析失败",
			"error", err,
			"raw_content", content,
			"content_length", len(content))
		return fmt.Errorf("JSON schema验证失败: %v", err)
	}

	// 业务逻辑验证
	if tempResult.Operation == "" {
		slog.Error("文件操作验证失败：缺少operation")
		return fmt.Errorf("operation字段不能为空")
	}

	if tempResult.Operation != "create" && tempResult.Operation != "delete" && 
	   tempResult.Operation != "move" && tempResult.Operation != "copy" {
		slog.Error("文件操作验证失败：operation类型无效", "operation", tempResult.Operation)
		return fmt.Errorf("operation必须是 create、delete、move、copy 之一")
	}

	if tempResult.SourcePath == "" {
		slog.Error("文件操作验证失败：缺少source_path", "content", content, "context_info", "")
		return fmt.Errorf("source_path字段不能为空")
	}

	// create操作通常需要content（除非创建空文件）
	if tempResult.Operation == "create" && tempResult.Content == "" {
		slog.Error("文件操作验证失败：create操作缺少content", "content", content, "context_info", "")
		return fmt.Errorf("create操作通常需要content字段")
	}

	// move和copy操作需要target_path
	if (tempResult.Operation == "move" || tempResult.Operation == "copy") && tempResult.TargetPath == "" {
		slog.Error("文件操作验证失败：move/copy操作缺少target_path", "operation", tempResult.Operation)
		return fmt.Errorf("move和copy操作需要target_path字段")
	}

	return nil
}
