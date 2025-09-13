package app

import "github.com/wailsapp/wails/v3/pkg/application"

const (
	WindowMain     = "main"
	WindowFloating = "floating"
	WindowSettings = "settings"
)

var DefaultManager = &WindowManager{
	winMap: make(map[string]*application.WebviewWindow),
}

func Initialize(app *application.App) {
	DefaultManager.app = app
	DefaultManager.buildContextMenu()
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

func EmitEvent(name string, data any) {
	DefaultManager.EmitEvent(name, data)
}
