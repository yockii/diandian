package main

import (
	"changeme/background/app"
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	a := application.New(application.Options{
		Name:        "aipc",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Initialize(a)
	err := app.DefaultManager.Run()

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
