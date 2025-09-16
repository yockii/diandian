package model

type Task struct {
	Base
	ConversationID uint64 `json:"conversation_id,string,omitempty" gorm:"index"`
	Name           string `json:"name" gorm:"size:200"`
	Description    string `json:"description"`
	Status         string `json:"status" gorm:"size:50"`      // pending, running, completed, failed, cancelled
	Progress       int    `json:"progress" gorm:"default:0"`  // 0-100
	Result         string `json:"result" gorm:"type:text"`    // 任务执行结果
	ErrorMsg       string `json:"error_msg" gorm:"type:text"` // 错误信息
}

// 任务状态常量
const (
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"
	TaskStatusCancelled = "cancelled"
)
