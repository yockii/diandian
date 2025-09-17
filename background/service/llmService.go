package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"diandian/background/constant"
	"diandian/background/database"
	"diandian/background/domain"
	"diandian/background/model"
	"diandian/background/service/operation"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

var DefaultLLMService = &LLMService{}

type LLMService struct{}

// 统一的消息处理响应结构
type UnifiedMessageResponse struct {
	ConversationTitle string                  `json:"conversation_title"`
	MessageType       string                  `json:"message_type"`              // "chat" or "automation"
	ChatResponse      string                  `json:"chat_response"`             // 聊天回复内容
	AutomationTask    *AutomationTaskResponse `json:"automation_task,omitempty"` // 自动化任务详情（仅当message_type为automation时）
	Confidence        float64                 `json:"confidence"`                // 0.0-1.0
	Explanation       string                  `json:"explanation"`               // 分类原因
}

// 自动化任务分析响应结构
type AutomationTaskResponse struct {
	TaskName     string   `json:"task_name"`     // 任务名称
	Description  string   `json:"description"`   // 任务描述
	Steps        []string `json:"steps"`         // 执行步骤
	Complexity   string   `json:"complexity"`    // simple, medium, complex
	Risks        []string `json:"risks"`         // 风险提示
	NeedsConfirm bool     `json:"needs_confirm"` // 是否需要用户确认
}

// TextModelConfig 文本模型配置
type TextModelConfig struct {
	BaseURL string
	Token   string
	Model   string
}

// VisionModelConfig 视觉模型配置
type VisionModelConfig struct {
	BaseURL string
	Token   string
	Model   string
}

// GetTextModelConfig 获取文本模型配置 (公开方法用于测试)
func (s *LLMService) GetTextModelConfig() (*TextModelConfig, error) {
	return s.getTextModelConfig()
}

// 获取文本模型配置
func (s *LLMService) getTextModelConfig() (*TextModelConfig, error) {
	var settings []*model.Setting
	err := database.DB.Where("key IN ?", []string{
		model.SettingKeyLlmTextBaseUrl,
		model.SettingKeyLlmTextToken,
		model.SettingKeyLlmTextModel,
	}).Find(&settings).Error

	if err != nil {
		return nil, fmt.Errorf("获取文本模型配置失败: %v", err)
	}

	config := &TextModelConfig{}
	for _, setting := range settings {
		if setting.Value == nil {
			continue
		}
		switch setting.Key {
		case model.SettingKeyLlmTextBaseUrl:
			config.BaseURL = *setting.Value
		case model.SettingKeyLlmTextToken:
			config.Token = *setting.Value
		case model.SettingKeyLlmTextModel:
			config.Model = *setting.Value
		}
	}

	if config.Token == "" || config.Model == "" {
		return nil, fmt.Errorf("文本模型配置不完整")
	}

	return config, nil
}

// GetVisionModelConfig 获取视觉模型配置 (公开方法用于测试)
func (s *LLMService) GetVisionModelConfig() (*VisionModelConfig, error) {
	return s.getVisionModelConfig()
}

// 获取视觉模型配置
func (s *LLMService) getVisionModelConfig() (*VisionModelConfig, error) {
	var settings []*model.Setting
	err := database.DB.Where("key IN ?", []string{
		model.SettingKeyLlmVlBaseUrl,
		model.SettingKeyLlmVlToken,
		model.SettingKeyLlmVlModel,
	}).Find(&settings).Error

	if err != nil {
		return nil, fmt.Errorf("获取视觉模型配置失败: %v", err)
	}

	config := &VisionModelConfig{}
	for _, setting := range settings {
		if setting.Value == nil {
			continue
		}
		switch setting.Key {
		case model.SettingKeyLlmVlBaseUrl:
			config.BaseURL = *setting.Value
		case model.SettingKeyLlmVlToken:
			config.Token = *setting.Value
		case model.SettingKeyLlmVlModel:
			config.Model = *setting.Value
		}
	}

	if config.Token == "" || config.Model == "" {
		return nil, fmt.Errorf("视觉模型配置不完整")
	}

	return config, nil
}

// CreateTextClient 创建文本模型客户端 (公开方法用于测试)
func (s *LLMService) CreateTextClient() (*openai.Client, string, error) {
	return s.createTextClient()
}

// 创建文本模型客户端
func (s *LLMService) createTextClient() (*openai.Client, string, error) {
	config, err := s.getTextModelConfig()
	if err != nil {
		return nil, "", fmt.Errorf("获取文本模型配置失败: %v", err)
	}

	clientConfig := openai.DefaultConfig(config.Token)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	client := openai.NewClientWithConfig(clientConfig)
	return client, config.Model, nil
}

// CreateVisionClient 创建视觉模型客户端 (公开方法用于测试)
func (s *LLMService) CreateVisionClient() (*openai.Client, string, error) {
	return s.createVisionClient()
}

// 创建视觉模型客户端
func (s *LLMService) createVisionClient() (*openai.Client, string, error) {
	config, err := s.getVisionModelConfig()
	if err != nil {
		return nil, "", fmt.Errorf("获取视觉模型配置失败: %v", err)
	}

	clientConfig := openai.DefaultConfig(config.Token)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	client := openai.NewClientWithConfig(clientConfig)
	return client, config.Model, nil
}

// cleanMarkdownCodeBlock 清理markdown代码块标记和其他格式标记
func cleanMarkdownCodeBlock(content string) string {
	content = strings.TrimSpace(content)

	// 移除开头的各种markdown代码块标记
	patterns := []string{
		"```json",
		"```JSON",
		"```javascript",
		"```js",
		"```",
		"``",
		"`",
	}

	for _, pattern := range patterns {
		if strings.HasPrefix(content, pattern) {
			content = content[len(pattern):]
			break
		}
	}

	// 移除结尾的各种markdown标记
	endPatterns := []string{
		"```",
		"``",
		"`",
	}

	for _, pattern := range endPatterns {
		if strings.HasSuffix(content, pattern) {
			content = content[:len(content)-len(pattern)]
			break
		}
	}

	// 移除可能的语言标识符行
	lines := strings.Split(content, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		// 如果第一行只包含语言标识符，移除它
		if firstLine == "json" || firstLine == "JSON" || firstLine == "javascript" || firstLine == "js" {
			lines = lines[1:]
			content = strings.Join(lines, "\n")
		}
	}

	// 移除多余的空白字符
	content = strings.TrimSpace(content)

	// 移除可能的BOM标记
	if strings.HasPrefix(content, "\ufeff") {
		content = content[3:]
	}

	return content
}

// retryLLMCall 重试LLM调用的通用方法
func (s *LLMService) retryLLMCall(
	callFunc func() (string, error),
	validateFunc func(content string) error,
	maxRetries int,
	operation string,
) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// 调用LLM
		content, err := callFunc()
		if err != nil {
			lastErr = fmt.Errorf("LLM调用失败: %v", err)
			if attempt < maxRetries {
				fmt.Printf("🔄 %s第%d次尝试失败，重试中... 错误: %v\n", operation, attempt, err)
				continue
			}
			break
		}

		// 清理内容
		cleanedContent := cleanMarkdownCodeBlock(content)

		// 验证内容
		if validateFunc != nil {
			if err := validateFunc(cleanedContent); err != nil {
				lastErr = fmt.Errorf("内容验证失败: %v", err)
				if attempt < maxRetries {
					fmt.Printf("🔄 %s第%d次尝试验证失败，重试中... 错误: %v\n", operation, attempt, err)
					continue
				}
				break
			}
		}

		// 成功
		if attempt > 1 {
			fmt.Printf("✅ %s在第%d次尝试后成功\n", operation, attempt)
		}
		return cleanedContent, nil
	}

	return "", fmt.Errorf("%s在%d次尝试后仍然失败，最后错误: %v", operation, maxRetries, lastErr)
}

// 简单的文本聊天接口
func (s *LLMService) SimpleChat(userMessage string) (string, error) {
	client, model, err := s.createTextClient()
	if err != nil {
		return "", err
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
			MaxTokens:   2000,
			Temperature: 0.7,
		},
	)

	if err != nil {
		return "", fmt.Errorf("调用LLM失败: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("API返回空响应")
	}

	return resp.Choices[0].Message.Content, nil
}

// 统一处理用户消息：同时进行聊天回复和任务判断
func (s *LLMService) ProcessMessage(conversationID uint64) (*model.Message, *UnifiedMessageResponse, error) {
	client, m, err := s.createTextClient()
	if err != nil {
		slog.Error("创建文本模型客户端失败", "error", err)
		return nil, nil, err
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: constant.PromptAnalyzeUserMessage,
		},
	}

	// 获取所有历史消息
	var msgs []*model.Message
	err = database.DB.Where("conversation_id = ?", conversationID).Order("created_at asc").Find(&msgs).Error
	if err != nil {
		return nil, nil, err
	}

	// 构造对话消息
	for _, msg := range msgs {
		role := openai.ChatMessageRoleUser
		if msg.Role == model.MessageRoleAssistant {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	var result UnifiedMessageResponse
	schema, err := jsonschema.GenerateSchemaForType(result)
	if err != nil {
		slog.Error("生成大模型schema规范失败", "error", err)
		return nil, nil, err
	}

	slog.Debug("准备调用大模型消息处理API")

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    m,
			Messages: messages,
			// MaxTokens:   1500,
			// Temperature: 0.3,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:   "UnifiedMessageResponse",
					Schema: schema,
					Strict: true,
				},
			},
		},
	)

	if err != nil {
		slog.Error("调用消息处理API失败", "error", err)
		return nil, nil, fmt.Errorf("调用消息处理API失败: %v", err)
	}

	if len(resp.Choices) == 0 {
		slog.Error("消息处理API返回空响应")
		return nil, nil, fmt.Errorf("消息处理API返回空响应")
	}

	slog.Debug("消息处理API返回", "content", resp.Choices[0].Message.Content)

	// 清理markdown标记并解析JSON响应
	cleanedContent := cleanMarkdownCodeBlock(resp.Choices[0].Message.Content)
	slog.Debug("清理后的内容", "cleaned_content", cleanedContent)

	err = schema.Unmarshal(cleanedContent, &result)
	if err != nil {
		slog.Error("解析消息处理结果失败",
			"error", err,
			"raw_content", resp.Choices[0].Message.Content,
			"cleaned_content", cleanedContent)
		return nil, nil, fmt.Errorf("解析消息处理结果失败: %v", err)
	}

	// 记录返回结果
	msg := &model.Message{
		ConversationID: conversationID,
		Role:           model.MessageRoleAssistant,
		Content:        resp.Choices[0].Message.Content,
	}
	database.DB.Create(msg)

	return msg, &result, nil
}

// 自动化任务分解响应结构（第一阶段：高级分解）
type AutomationTaskDecomposition struct {
	TaskType        string               `json:"task_type"`        // simple, composite, complex
	Description     string               `json:"description"`      // 任务描述
	Steps           []AutomationStepPlan `json:"steps"`            // 执行步骤计划
	ExpectedOutcome string               `json:"expected_outcome"` // 预期结果
	RiskLevel       string               `json:"risk_level"`       // low, medium, high
	EstimatedTime   int                  `json:"estimated_time"`   // 预估时间(秒)
}

// 自动化步骤计划（高级步骤，不包含具体参数）
type AutomationStepPlan struct {
	StepType               string `json:"step_type"`                // "click", "type", "launch_app", "file", "screenshot", "clipboard", "wait", "key_press"
	Description            string `json:"description"`              // 步骤描述
	RequiresScreenAnalysis bool   `json:"requires_screen_analysis"` // 是否需要屏幕分析
	Context                string `json:"context"`                  // 上下文信息，用于第二阶段生成具体操作
	Priority               int    `json:"priority"`                 // 优先级 1-10
	Optional               bool   `json:"optional"`                 // 是否可选
}

// ===== 第二阶段：具体操作定义 =====

// 点击操作
type ClickOperation struct {
	X      int    `json:"x"`      // X坐标
	Y      int    `json:"y"`      // Y坐标
	Button string `json:"button"` // "left", "right", "middle"
}

// 输入操作
type TypeOperation struct {
	Text string `json:"text"` // 要输入的文本
}

// 启动应用操作
type LaunchAppOperation struct {
	AppName string `json:"app_name"` // 应用名称，必须是预定义的
}

// 文件操作
type FileOperation struct {
	Operation  string `json:"operation"`             // "create", "delete", "move", "copy"
	SourcePath string `json:"source_path"`           // 源路径
	TargetPath string `json:"target_path,omitempty"` // 目标路径（移动/复制时需要）
	Content    string `json:"content,omitempty"`     // 文件内容（创建时需要）
}

// 截屏操作
type ScreenshotOperation struct {
	Path string `json:"path"` // 保存路径
}

// 剪贴板操作
type ClipboardOperation struct {
	Operation string `json:"operation"`      // "get", "set"
	Text      string `json:"text,omitempty"` // 设置时的文本内容
}

// 等待操作
type WaitOperation struct {
	Duration int `json:"duration"` // 等待时间（毫秒）
}

// 按键操作
type KeyPressOperation struct {
	Key       string   `json:"key"`                 // 主键
	Modifiers []string `json:"modifiers,omitempty"` // 修饰键 ["ctrl", "alt", "shift"]
}

// 统一的操作响应结构
type OperationResponse struct {
	OperationType string      `json:"operation_type"` // 操作类型
	Operation     interface{} `json:"operation"`      // 具体操作（会被序列化为对应的操作类型）
}

// 视觉分析响应结构
type VisualAnalysisResponse struct {
	ElementsFound   []VisualElement        `json:"elements_found"`  // 找到的元素
	ScreenInfo      ScreenInfo             `json:"screen_info"`     // 屏幕信息
	Recommendations []ActionRecommendation `json:"recommendations"` // 操作建议
}

// 视觉元素
type VisualElement struct {
	Type        string      `json:"type"`         // 元素类型
	Description string      `json:"description"`  // 元素描述
	Coordinates Coordinates `json:"coordinates"`  // 坐标信息
	Confidence  float64     `json:"confidence"`   // 置信度
	TextContent string      `json:"text_content"` // 文本内容
	Clickable   bool        `json:"clickable"`    // 是否可点击
}

// 坐标信息
type Coordinates struct {
	X      int `json:"x"`      // x坐标
	Y      int `json:"y"`      // y坐标
	Width  int `json:"width"`  // 宽度
	Height int `json:"height"` // 高度
}

// 屏幕信息
type ScreenInfo struct {
	Resolution         string `json:"resolution"`          // 分辨率
	ActiveWindow       string `json:"active_window"`       // 活动窗口
	OverallDescription string `json:"overall_description"` // 整体描述
}

// 操作建议
type ActionRecommendation struct {
	Action string `json:"action"` // 建议操作
	Target string `json:"target"` // 操作目标
	Reason string `json:"reason"` // 建议原因
}

// 分析自动化任务并分解为具体步骤
func (s *LLMService) DecomposeAutomationTask(conversationHistory []openai.ChatCompletionMessage) (*domain.AutomationTaskDecomposition, error) {
	client, model, err := s.createTextClient()
	if err != nil {
		return nil, err
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: constant.PromptAutomationTaskDecomposition,
		},
	}

	// 添加对话历史
	messages = append(messages, conversationHistory...)

	var result domain.AutomationTaskDecomposition
	schema, err := jsonschema.GenerateSchemaForType(result)
	if err != nil {
		return nil, fmt.Errorf("生成schema失败: %v", err)
	}

	// 定义LLM调用函数
	callFunc := func() (string, error) {
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    model,
				Messages: messages,
				ResponseFormat: &openai.ChatCompletionResponseFormat{
					Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
					JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
						Name:   "AutomationTaskDecomposition",
						Schema: schema,
						Strict: true,
					},
				},
			},
		)

		if err != nil {
			return "", err
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("API返回空响应")
		}

		return resp.Choices[0].Message.Content, nil
	}

	// 定义验证函数
	validateFunc := func(content string) error {
		var tempResult AutomationTaskDecomposition
		if err := schema.Unmarshal(content, &tempResult); err != nil {
			// 记录原始LLM返回内容用于排查
			slog.Error("任务分解JSON解析失败",
				"error", err,
				"raw_content", content,
				"content_length", len(content))
			return fmt.Errorf("JSON schema验证失败: %v", err)
		}

		// 验证必要字段
		if tempResult.TaskType == "" {
			slog.Error("任务分解验证失败：缺少task_type", "content", content)
			return fmt.Errorf("缺少task_type字段")
		}
		if tempResult.Description == "" {
			slog.Error("任务分解验证失败：缺少description", "content", content)
			return fmt.Errorf("缺少description字段")
		}
		if len(tempResult.Steps) == 0 {
			slog.Error("任务分解验证失败：缺少steps", "content", content)
			return fmt.Errorf("缺少steps字段或步骤为空")
		}

		// 验证每个步骤
		for i, step := range tempResult.Steps {
			if step.StepType == "" {
				slog.Error("任务分解验证失败：步骤缺少step_type",
					"step_index", i+1,
					"step", step,
					"content", content)
				return fmt.Errorf("步骤%d缺少step_type字段", i+1)
			}
			if step.Description == "" {
				slog.Error("任务分解验证失败：步骤缺少description",
					"step_index", i+1,
					"step", step,
					"content", content)
				return fmt.Errorf("步骤%d缺少description字段", i+1)
			}
		}

		// 记录成功的解析结果
		slog.Info("任务分解验证成功",
			"task_type", tempResult.TaskType,
			"steps_count", len(tempResult.Steps),
			"risk_level", tempResult.RiskLevel)

		return nil
	}

	// 使用重试机制调用LLM
	content, err := s.retryLLMCall(callFunc, validateFunc, 3, "任务分解")
	if err != nil {
		return nil, fmt.Errorf("任务分解失败: %v", err)
	}

	// 最终解析 - 再次清理确保万无一失
	finalCleanedContent := cleanMarkdownCodeBlock(content)
	slog.Debug("最终清理后的内容", "final_cleaned_content", finalCleanedContent)

	err = schema.Unmarshal(finalCleanedContent, &result)
	if err != nil {
		slog.Error("解析任务分解结果失败",
			"error", err,
			"original_content", content,
			"final_cleaned_content", finalCleanedContent)
		return nil, fmt.Errorf("解析任务分解结果失败: %v", err)
	}

	return &result, nil
}

// 分析屏幕截图
func (s *LLMService) AnalyzeScreenshot(imageData []byte, analysisRequest string) (*domain.VisualAnalysisResponse, error) {
	generator := operation.NewVisionGenerator()

	result, err := generator.Analyze(imageData, analysisRequest)
	if err != nil {
		return nil, err
	}

	// 转换回原类型
	elements := make([]domain.VisualElement, len(result.ElementsFound))
	for i, elem := range result.ElementsFound {
		elements[i] = domain.VisualElement{
			Type:        elem.Type,
			Description: elem.Description,
			Coordinates: domain.Coordinates{
				X:      elem.Coordinates.X,
				Y:      elem.Coordinates.Y,
				Width:  elem.Coordinates.Width,
				Height: elem.Coordinates.Height,
			},
			Confidence:  elem.Confidence,
			TextContent: elem.TextContent,
			Clickable:   elem.Clickable,
		}
	}

	recommendations := make([]domain.ActionRecommendation, len(result.Recommendations))
	for i, rec := range result.Recommendations {
		recommendations[i] = domain.ActionRecommendation{
			Action: rec.Type,
			Target: rec.Description,
			Reason: fmt.Sprintf("Priority: %d", rec.Priority),
		}
	}

	return &domain.VisualAnalysisResponse{
		ElementsFound:   elements,
		ScreenInfo:      domain.ScreenInfo{}, // 保持空的ScreenInfo以兼容现有代码
		Recommendations: recommendations,
	}, nil
}

// ===== 第二阶段：具体操作生成方法 =====

// 生成点击操作
func (s *LLMService) GenerateClickOperation(contextInfo string, screenAnalysis *domain.VisualAnalysisResponse) (*domain.ClickOperation, error) {
	generator := operation.NewClickGenerator()

	// 简化转换：只传递必要信息
	var operationScreenAnalysis *operation.VisualAnalysisResponse
	if screenAnalysis != nil {
		operationElements := make([]operation.VisualElement, len(screenAnalysis.ElementsFound))
		for i, elem := range screenAnalysis.ElementsFound {
			operationElements[i] = operation.VisualElement{
				Type:        elem.Type,
				Description: elem.Description,
				Coordinates: operation.Coordinates{
					X:      elem.Coordinates.X,
					Y:      elem.Coordinates.Y,
					Width:  elem.Coordinates.Width,
					Height: elem.Coordinates.Height,
				},
				Confidence:  elem.Confidence,
				TextContent: elem.TextContent,
				Clickable:   elem.Clickable,
			}
		}

		operationScreenAnalysis = &operation.VisualAnalysisResponse{
			ElementsFound: operationElements,
		}
	}

	result, err := generator.Generate(contextInfo, operationScreenAnalysis)
	if err != nil {
		return nil, err
	}

	// 转换回原类型
	return &domain.ClickOperation{
		X:      result.X,
		Y:      result.Y,
		Button: result.Button,
	}, nil
}

// 生成输入操作
func (s *LLMService) GenerateTypeOperation(contextInfo string) (*domain.TypeOperation, error) {
	generator := operation.NewTypeGenerator()

	result, err := generator.Generate(contextInfo)
	if err != nil {
		return nil, err
	}

	// 转换回原类型
	return &domain.TypeOperation{
		Text: result.Text,
	}, nil
}

// 生成文件操作
func (s *LLMService) GenerateFileOperation(contextInfo string) (*domain.FileOperation, error) {
	generator := operation.NewFileGenerator()

	result, err := generator.Generate(contextInfo)
	if err != nil {
		return nil, err
	}

	// 转换回原类型
	return &domain.FileOperation{
		Operation:  result.Operation,
		SourcePath: result.SourcePath,
		TargetPath: result.TargetPath,
		Content:    result.Content,
	}, nil
}
