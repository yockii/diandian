package main

import (
	"diandian/background/app"
	"diandian/background/service"
	"diandian/background/util"
	"embed"
	"log"
	"log/slog"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	util.InitializeLog()

	a := application.New(application.Options{
		Name:        "DianDian",
		Description: "点点小助理，AI聊天，自动操作电脑",
		Services: []application.Service{
			application.NewService(&service.WindowService{}),
			application.NewService(&service.MessageService{}),
			application.NewService(&service.SettingService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	a.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		slog.Info("程序启动中...")
		app.OnAppStart()
		service.InitializeData()
	})

	app.Initialize(a)

	slog.Info("程序初始化完成，运行中...")

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
