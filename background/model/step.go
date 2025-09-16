package model

type Step struct {
	Base
	TaskID      uint64 `json:"task_id,string"`
	Content     string `json:"content"`     // 展示给用户的消息内容
	StepType    string `json:"step_type"`   // message, action, screenshot, analysis
	Status      string `json:"status"`      // pending, running, completed, failed
	ActionType  string `json:"action_type"` // click, type, key, scroll, wait
	Coordinates string `json:"coordinates"` // 操作坐标 (x,y)
	ActionData  string `json:"action_data"` // 操作数据（如输入的文本、按键等）
	Screenshot  string `json:"screenshot"`  // 截图文件路径或base64
	Result      string `json:"result"`      // 步骤执行结果
	ErrorMsg    string `json:"error_msg"`   // 错误信息
}

// 步骤类型常量
const (
	StepTypeMessage    = "message"    // 消息显示
	StepTypeAction     = "action"     // 操作执行
	StepTypeScreenshot = "screenshot" // 截图
	StepTypeAnalysis   = "analysis"   // 分析
)

// 步骤状态常量
const (
	StepStatusPending   = "pending"
	StepStatusRunning   = "running"
	StepStatusCompleted = "completed"
	StepStatusFailed    = "failed"
)

// 操作类型常量
const (
	ActionTypeClick  = "click"
	ActionTypeType   = "type"
	ActionTypeKey    = "key"
	ActionTypeScroll = "scroll"
	ActionTypeWait   = "wait"
)
