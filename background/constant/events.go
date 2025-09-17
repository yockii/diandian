package constant

const (
	EventThemeChanged       = "theme-changed"
	EventCanWorkChanged     = "can-work-changed"
	EventStickySideChanged  = "sticky-side-changed"
	EventMouseEnterFloating = "mouse-enter-floating"
	EventMouseLeaveFloating = "mouse-leave-floating"

	EventMessageResponsed  = "message-responsed"
	EventTaskStatusChanged = "task-status-changed"
	EventOperateFailed     = "operate-failed" // model.Step

	EventNotify = "notify"

	// 任务执行相关事件
	EventTaskExecutionStarted   = "task-execution-started"
	EventTaskExecutionCompleted = "task-execution-completed"
)
