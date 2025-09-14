package model

type Step struct {
	Base
	TaskID  uint64 `json:"task_id"`
	Content string `json:"content"` // 展示给用户的消息内容
	Status  int    `json:"status"`  // 1-普通文本回复，2-后台执行任务，3-最小化执行模拟操作
}
