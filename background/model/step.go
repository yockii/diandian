package model

type Step struct {
	Status  int    `json:"status"`  // 1-普通文本回复，2-后台执行任务，3-最小化执行模拟操作
	Content string `json:"content"` // 展示给用户的消息内容
}
