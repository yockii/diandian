package domain

import (
	"time"
)

// ===== 任务分解相关结构体 =====

// AutomationTaskDecomposition 自动化任务分解响应结构（第一阶段：高级分解）
type AutomationTaskDecomposition struct {
	TaskType        string               `json:"task_type"`        // simple, medium, complex
	Description     string               `json:"description"`      // 任务描述
	Steps           []AutomationStepPlan `json:"steps"`            // 执行步骤计划
	ExpectedOutcome string               `json:"expected_outcome"` // 预期结果
	RiskLevel       string               `json:"risk_level"`       // low, medium, high
	EstimatedTime   int                  `json:"estimated_time"`   // 预估时间(秒)
}

// AutomationStepPlan 自动化步骤计划（高级步骤，不包含具体参数）
type AutomationStepPlan struct {
	Type                   string `json:"type"`                     // click, type, launch_app, file, screenshot, clipboard, wait, key_press
	Description            string `json:"description"`              // 步骤描述
	RequiresScreenAnalysis bool   `json:"requires_screen_analysis"` // 是否需要屏幕分析
	Context                string `json:"context"`                  // 上下文信息，用于第二阶段生成具体操作
	Priority               int    `json:"priority"`                 // 优先级 1-10
	Optional               bool   `json:"optional"`                 // 是否可选
}

// ===== 具体操作结构体 =====

// ClickOperation 点击操作
type ClickOperation struct {
	X      int    `json:"x"`      // X坐标
	Y      int    `json:"y"`      // Y坐标
	Button string `json:"button"` // "left", "right", "middle"
}

// TypeOperation 输入操作
type TypeOperation struct {
	Text string `json:"text"` // 要输入的文本
}

// LaunchAppOperation 启动应用操作
type LaunchAppOperation struct {
	AppName string `json:"app_name"` // 应用名称，必须是预定义的
}

// FileOperation 文件操作
type FileOperation struct {
	Operation  string `json:"operation"`             // "create", "delete", "move", "copy"
	SourcePath string `json:"source_path"`           // 源路径
	TargetPath string `json:"target_path,omitempty"` // 目标路径（移动/复制时需要）
	Content    string `json:"content,omitempty"`     // 文件内容（创建时需要）
}

// ScreenshotOperation 截屏操作
type ScreenshotOperation struct {
	Path string `json:"path"` // 保存路径
}

// ClipboardOperation 剪贴板操作
type ClipboardOperation struct {
	Operation string `json:"operation"`      // "get", "set"
	Text      string `json:"text,omitempty"` // 设置时的文本内容
}

// WaitOperation 等待操作
type WaitOperation struct {
	Duration int `json:"duration"` // 等待时间（毫秒）
}

// KeyPressOperation 按键操作
type KeyPressOperation struct {
	Key       string   `json:"key"`                 // 主键
	Modifiers []string `json:"modifiers,omitempty"` // 修饰键 ["ctrl", "alt", "shift"]
}

// ===== 视觉分析相关结构体 =====

// VisualAnalysisResponse 视觉分析响应结构
type VisualAnalysisResponse struct {
	ElementsFound   []VisualElement        `json:"elements_found"`  // 找到的元素
	ScreenInfo      ScreenInfo             `json:"screen_info"`     // 屏幕信息
	Recommendations []ActionRecommendation `json:"recommendations"` // 操作建议
}

// VisualElement 视觉元素
type VisualElement struct {
	Type        string      `json:"type"`         // 元素类型
	Description string      `json:"description"`  // 元素描述
	Coordinates Coordinates `json:"coordinates"`  // 坐标信息
	Confidence  float64     `json:"confidence"`   // 置信度
	TextContent string      `json:"text_content"` // 文本内容
	Clickable   bool        `json:"clickable"`    // 是否可点击
}

// Coordinates 坐标信息
type Coordinates struct {
	X      int `json:"x"`      // x坐标
	Y      int `json:"y"`      // y坐标
	Width  int `json:"width"`  // 宽度
	Height int `json:"height"` // 高度
}

// ScreenInfo 屏幕信息
type ScreenInfo struct {
	Resolution         string `json:"resolution"`          // 分辨率
	ActiveWindow       string `json:"active_window"`       // 活动窗口
	OverallDescription string `json:"overall_description"` // 整体描述
}

// ActionRecommendation 操作建议
type ActionRecommendation struct {
	Action string `json:"action"` // 建议操作
	Target string `json:"target"` // 操作目标
	Reason string `json:"reason"` // 建议原因
}

// ===== 执行结果相关结构体 =====

// TaskExecutionResult 任务执行结果
type TaskExecutionResult struct {
	TaskID      uint                   `json:"task_id"`
	Success     bool                   `json:"success"`
	Message     string                 `json:"message"`
	Steps       []*StepExecutionResult `json:"steps"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     time.Time              `json:"end_time"`
	Duration    time.Duration          `json:"duration"`
	ErrorCount  int                    `json:"error_count"`
	SuccessRate float64                `json:"success_rate"`
}

// StepExecutionResult 步骤执行结果
type StepExecutionResult struct {
	StepIndex      int           `json:"step_index"`
	StepType       string        `json:"step_type"`
	Success        bool          `json:"success"`
	Message        string        `json:"message"`
	Error          string        `json:"error,omitempty"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	Duration       time.Duration `json:"duration"`
	RetryCount     int           `json:"retry_count"`
	ScreenshotPath string        `json:"screenshot_path,omitempty"` // 执行前的截图路径
}

// ===== 统一的操作响应结构 =====

// OperationResponse 统一的操作响应结构
type OperationResponse struct {
	OperationType string      `json:"operation_type"` // 操作类型
	Operation     interface{} `json:"operation"`      // 具体操作（会被序列化为对应的操作类型）
}
