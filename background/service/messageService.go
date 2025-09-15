package service

import (
	"changeme/background/app"
	"changeme/background/model"
	"time"
)

type MessageService struct{}

func (s *MessageService) NewMessage(msg *model.Message) {
	time.Sleep(5 * time.Second)
	app.EmitEvent("new-msg", &model.Step{
		Status:  1,
		Content: "我已收到: " + msg.Content,
	})
}
