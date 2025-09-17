package app

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	WindowMain     = "main"
	WindowFloating = "floating"
	WindowSettings = "settings"
)

var DefaultManager = &WindowManager{
	winMap: make(map[string]*application.WebviewWindow),
}

func OnAppStart() {
	DefaultManager.OnAppStart()
}

func Initialize(app *application.App) {
	DefaultManager.app = app
	DefaultManager.buildContextMenu()
}

func IsInitializeSuccess() bool {
	return DefaultManager.IsInitializeSuccess()
}

func Run() error {
	return DefaultManager.Run()
}

func ShowFloating() {
	DefaultManager.ShowFloating()
}

func ShowMain() {
	DefaultManager.ShowMain()
}

func ShowSettings() {
	DefaultManager.ShowSettings()
}

func HideSettings() {
	DefaultManager.HideSettings()
}

func EmitEvent(name string, data any) {
	DefaultManager.EmitEvent(name, data)
}

func FloatingStickySide() int {
	return DefaultManager.FloatingStickySide()
}

func GetApp() *application.App {
	return DefaultManager.app
}
