package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"diandian/background/constant"
	"diandian/background/database"
	"diandian/background/model"

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

// 获取LLM配置
func (s *LLMService) getLLMConfig() (baseURL, token, textModel, vlModel string, err error) {
	var settings []*model.Setting
	err = database.DB.Where("key IN ?", []string{
		model.SettingKeyLlmBaseUrl,
		model.SettingKeyLlmToken,
		model.SettingKeyTextModel,
		model.SettingKeyVlModel,
	}).Find(&settings).Error

	if err != nil {
		return
	}

	for _, setting := range settings {
		if setting.Value == nil {
			continue
		}
		switch setting.Key {
		case model.SettingKeyLlmBaseUrl:
			baseURL = *setting.Value
		case model.SettingKeyLlmToken:
			token = *setting.Value
		case model.SettingKeyTextModel:
			textModel = *setting.Value
		case model.SettingKeyVlModel:
			vlModel = *setting.Value
		}
	}

	return
}

// 创建OpenAI客户端
func (s *LLMService) createClient() (*openai.Client, string, error) {
	baseURL, token, textModel, _, err := s.getLLMConfig()
	if err != nil {
		return nil, "", fmt.Errorf("获取LLM配置失败: %v", err)
	}

	if baseURL == "" || token == "" || textModel == "" {
		return nil, "", fmt.Errorf("LLM配置不完整")
	}

	config := openai.DefaultConfig(token)
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	client := openai.NewClientWithConfig(config)
	return client, textModel, nil
}

// 简单的文本聊天接口
func (s *LLMService) SimpleChat(userMessage string) (string, error) {
	client, model, err := s.createClient()
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
	client, m, err := s.createClient()
	if err != nil {
		slog.Error("创建LLM客户端失败", "error", err)
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

	// 解析JSON响应
	err = schema.Unmarshal(resp.Choices[0].Message.Content, &result)
	if err != nil {
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

// 分析自动化任务
func (s *LLMService) AnalyzeAutomationTask(userMessage string) (*AutomationTaskResponse, error) {
	client, model, err := s.createClient()
	if err != nil {
		return nil, err
	}

	systemPrompt := `你是一个桌面自动化专家，需要分析用户的自动化任务需求。请分析用户的消息，提供详细的任务分析。

请严格按照以下JSON格式返回结果：
{
  "task_name": "任务的简短名称",
  "description": "任务的详细描述",
  "steps": ["步骤1", "步骤2", "步骤3"],
  "complexity": "simple/medium/complex",
  "risks": ["风险1", "风险2"],
  "needs_confirm": true/false
}

complexity说明：
- simple: 简单操作，如打开软件、截图等
- medium: 中等复杂度，如文件整理、简单的数据处理
- complex: 复杂操作，如涉及多个软件的协同操作

needs_confirm说明：
- 涉及文件删除、系统设置修改、发送邮件等操作时设为true
- 简单的查看、截图等操作可设为false

用户任务：` + userMessage

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
			},
			MaxTokens:   1000,
			Temperature: 0.3,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("调用任务分析API失败: %v", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("任务分析API返回空响应")
	}

	// 解析JSON响应
	var result AutomationTaskResponse
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result)
	if err != nil {
		return nil, fmt.Errorf("解析任务分析结果失败: %v", err)
	}

	return &result, nil
}
