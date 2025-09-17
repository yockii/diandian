package operation

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"diandian/background/constant"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// VisionGenerator 视觉分析生成器
type VisionGenerator struct {
	*BaseGenerator
}

// NewVisionGenerator 创建视觉分析生成器
func NewVisionGenerator() *VisionGenerator {
	return &VisionGenerator{
		BaseGenerator: NewBaseGenerator(),
	}
}

// Analyze 分析屏幕截图
func (g *VisionGenerator) Analyze(imageData []byte, analysisRequest string) (*VisualAnalysisResponse, error) {
	// 首先尝试使用JSON格式
	result, err := g.analyzeWithJSONFormat(imageData, analysisRequest)
	if err != nil {
		slog.Warn("JSON格式分析失败，尝试降级到文本格式", "error", err)
		// 降级到文本格式，然后转换为JSON
		return g.analyzeWithTextFallback(imageData, analysisRequest)
	}
	return result, nil
}

// analyzeWithJSONFormat 使用JSON格式进行分析
func (g *VisionGenerator) analyzeWithJSONFormat(imageData []byte, analysisRequest string) (*VisualAnalysisResponse, error) {
	client, model, err := g.createVisionClient()
	if err != nil {
		slog.Error("创建视觉模型客户端失败", "error", err)
		return nil, err
	}

	// 将图片转换为base64
	imageBase64 := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(imageData))

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: constant.PromptVisualAnalysis,
		},
		{
			Role: openai.ChatMessageRoleUser,
			MultiContent: []openai.ChatMessagePart{
				{
					Type: openai.ChatMessagePartTypeText,
					Text: analysisRequest,
				},
				{
					Type: openai.ChatMessagePartTypeImageURL,
					ImageURL: &openai.ChatMessageImageURL{
						URL: imageBase64,
					},
				},
			},
		},
	}

	// 使用重试机制调用LLM
	content, err := g.retryLLMCall(
		func() (string, error) {
			return g.callVisionLLM(client, model, messages, true) // 尝试JSON格式
		},
		func(content string) error {
			return g.validateVisualAnalysis(content)
		},
		2, // JSON格式只重试2次
		"视觉分析(JSON格式)",
	)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var result VisualAnalysisResponse
	schema, err := jsonschema.GenerateSchemaForType(result)
	if err != nil {
		return nil, fmt.Errorf("生成视觉分析schema失败: %v", err)
	}

	if err := schema.Unmarshal(content, &result); err != nil {
		slog.Error("视觉分析JSON解析失败",
			"error", err,
			"raw_content", content,
			"content_length", len(content))
		return nil, fmt.Errorf("JSON schema验证失败: %v", err)
	}

	slog.Info("视觉分析验证成功", "elements_count", len(result.ElementsFound))
	return &result, nil
}

// analyzeWithTextFallback 降级到文本格式分析
func (g *VisionGenerator) analyzeWithTextFallback(imageData []byte, analysisRequest string) (*VisualAnalysisResponse, error) {
	slog.Info("使用文本降级模式进行视觉分析")

	// 第一步：使用视觉模型生成文本描述
	textDescription, err := g.generateTextDescription(imageData, analysisRequest)
	if err != nil {
		return nil, fmt.Errorf("生成文本描述失败: %v", err)
	}

	// 第二步：使用文本模型将描述转换为JSON
	return g.convertTextToJSON(textDescription, analysisRequest)
}

// generateTextDescription 生成文本描述
func (g *VisionGenerator) generateTextDescription(imageData []byte, analysisRequest string) (string, error) {
	client, model, err := g.createVisionClient()
	if err != nil {
		return "", err
	}

	// 将图片转换为base64
	imageBase64 := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(imageData))

	// 使用简化的文本提示词
	textPrompt := fmt.Sprintf(`请详细描述这个屏幕截图中的所有可交互元素，包括：
1. 按钮的位置和文字
2. 输入框的位置
3. 文本内容
4. 图标和链接
5. 窗口和对话框

用户请求：%s

请用自然语言详细描述，不要使用JSON格式。`, analysisRequest)

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleUser,
			MultiContent: []openai.ChatMessagePart{
				{
					Type: openai.ChatMessagePartTypeText,
					Text: textPrompt,
				},
				{
					Type: openai.ChatMessagePartTypeImageURL,
					ImageURL: &openai.ChatMessageImageURL{
						URL: imageBase64,
					},
				},
			},
		},
	}

	return g.callVisionLLM(client, model, messages, false) // 不使用JSON格式
}

// convertTextToJSON 将文本描述转换为JSON
func (g *VisionGenerator) convertTextToJSON(textDescription, originalRequest string) (*VisualAnalysisResponse, error) {
	client, model, err := g.createTextClient()
	if err != nil {
		return nil, err
	}

	// 构建转换提示词
	convertPrompt := fmt.Sprintf(`请将以下视觉分析文本描述转换为标准的JSON格式。

原始用户请求：%s

视觉分析描述：
%s

请严格按照以下JSON格式输出，不要添加任何其他内容：
{
  "elements_found": [
    {
      "type": "元素类型",
      "description": "元素描述",
      "coordinates": {"x": 0, "y": 0, "width": 0, "height": 0},
      "confidence": 0.9,
      "text_content": "文本内容",
      "clickable": true
    }
  ],
  "screen_info": {"width": 1920, "height": 1080},
  "recommendations": [
    {
      "type": "建议类型",
      "description": "建议描述",
      "priority": 1
    }
  ]
}

重要要求：
1. 直接输出JSON，不要使用markdown标签包裹
2. 不要输出其他内容，只输出JSON
3. 确保JSON格式完全正确`, originalRequest, textDescription)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你是一个专业的数据转换专家，负责将文本描述转换为结构化的JSON数据。",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: convertPrompt,
		},
	}

	// 使用重试机制调用文本模型
	content, err := g.retryLLMCall(
		func() (string, error) {
			return g.callTextLLM(client, model, messages)
		},
		func(content string) error {
			return g.validateVisualAnalysis(content)
		},
		3,
		"文本转JSON",
	)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var result VisualAnalysisResponse
	schema, err := jsonschema.GenerateSchemaForType(result)
	if err != nil {
		return nil, fmt.Errorf("生成视觉分析schema失败: %v", err)
	}

	if err := schema.Unmarshal(content, &result); err != nil {
		slog.Error("视觉分析JSON解析失败",
			"error", err,
			"raw_content", content,
			"content_length", len(content))
		return nil, fmt.Errorf("JSON schema验证失败: %v", err)
	}

	slog.Info("视觉分析(降级模式)验证成功", "elements_count", len(result.ElementsFound))
	return &result, nil
}

// callVisionLLM 调用视觉LLM
func (g *VisionGenerator) callVisionLLM(client *openai.Client, model string, messages []openai.ChatCompletionMessage, useJSONFormat bool) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	request := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	// 如果支持JSON格式，添加response_format
	if useJSONFormat {
		request.ResponseFormat = &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		}
	}

	resp, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		// 检查是否是response_format不支持的错误
		if useJSONFormat && strings.Contains(err.Error(), "response_format") {
			return "", fmt.Errorf("模型不支持JSON格式: %v", err)
		}
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("LLM返回空响应")
	}

	return resp.Choices[0].Message.Content, nil
}

// callTextLLM 调用文本LLM
func (g *VisionGenerator) callTextLLM(client *openai.Client, model string, messages []openai.ChatCompletionMessage) (string, error) {
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

// validateVisualAnalysis 验证视觉分析结果
func (g *VisionGenerator) validateVisualAnalysis(content string) error {
	var tempResult VisualAnalysisResponse
	schema, err := jsonschema.GenerateSchemaForType(tempResult)
	if err != nil {
		return fmt.Errorf("生成schema失败: %v", err)
	}

	if err := schema.Unmarshal(content, &tempResult); err != nil {
		slog.Error("视觉分析JSON解析失败",
			"error", err,
			"raw_content", content,
			"content_length", len(content))
		return fmt.Errorf("JSON schema验证失败: %v", err)
	}

	// 基本验证：至少要有元素或建议
	if len(tempResult.ElementsFound) == 0 && len(tempResult.Recommendations) == 0 {
		slog.Error("视觉分析验证失败：没有找到元素或建议")
		return fmt.Errorf("必须至少包含元素或操作建议")
	}

	return nil
}
