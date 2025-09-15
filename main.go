package main

import (
	"changeme/background/app"
	"changeme/background/service"
	"changeme/background/util"
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
		Name:        "aipc",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
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

	// app.Window.NewWithOptions(application.WebviewWindowOptions{
	// 	Title:     "Window 1",
	// 	Frameless: true, // 无边框窗口
	// 	Width:     400,
	// 	Height:    800,
	// 	Mac: application.MacWindow{
	// 		InvisibleTitleBarHeight: 50,
	// 		Backdrop:                application.MacBackdropTranslucent,
	// 		TitleBar:                application.MacTitleBarHiddenInset,
	// 	},
	// 	BackgroundColour: application.NewRGB(27, 38, 54),
	// 	// URL:              "/",
	// })

	// // go func() {
	// // 	for {
	// // 		now := time.Now().Format(time.RFC1123)
	// // 		app.Event.Emit("time", now)
	// // 		time.Sleep(time.Second)
	// // 	}
	// // }()

	// // Run the application. This blocks until the application has been exited.
	// err := app.Run()

	// // If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
