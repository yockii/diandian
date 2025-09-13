package service

import (
	"changeme/background/app"
	"changeme/background/model"
	"time"
)

type TaskService struct{}

func (s *TaskService) NewTask(msg string) {
	time.Sleep(5 * time.Second)
	app.EmitEvent("new-step", &model.Step{
		Status:  1,
		Content: "我已收到: " + msg,
	})
}
