package domain

// ===== LLM相关结构体 =====

// UserMessageAnalysisResponse 用户消息分析响应结构
type UserMessageAnalysisResponse struct {
	MessageType       string                  `json:"message_type"`              // "chat" or "automation"
	ChatResponse      string                  `json:"chat_response"`             // 聊天回复内容
	AutomationTask    *AutomationTaskResponse `json:"automation_task,omitempty"` // 自动化任务详情（仅当message_type为automation时）
	Confidence        float64                 `json:"confidence"`                // 0.0-1.0
	Explanation       string                  `json:"explanation"`               // 分类原因
}

// AutomationTaskResponse 自动化任务分析响应结构
type AutomationTaskResponse struct {
	TaskType     string   `json:"task_type"`     // simple, medium, complex
	Description  string   `json:"description"`   // 任务描述
	Steps        []string `json:"steps"`         // 执行步骤
	Complexity   string   `json:"complexity"`    // simple, medium, complex
	Risks        []string `json:"risks"`         // 风险提示
	NeedsConfirm bool     `json:"needs_confirm"` // 是否需要用户确认
}

// ===== 模型配置相关结构体 =====

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
