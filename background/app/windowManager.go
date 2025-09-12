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

type WindowManager struct {
	app    *application.App
	winMap map[string]*application.WebviewWindow
}

func Initialize(app *application.App) {
	DefaultManager.app = app
}

func (wm *WindowManager) Run() error {
	wm.GetWindow(WindowMain)
	return wm.app.Run()
}

func (wm *WindowManager) initializeMain() {
	wm.winMap[WindowMain] = wm.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "Window 1",
		Frameless: true, // 无边框窗口
		Width:     400,
		Height:    800,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})
}

func (wm *WindowManager) initializeFloating() {
	wm.winMap[WindowFloating] = wm.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "Floating Window",
		Frameless:      true, // 无边框窗口
		Width:          60,
		Height:         60,
		BackgroundType: application.BackgroundTypeTransparent,
		URL:            "/floating",
	})
}

func (wm *WindowManager) GetWindow(name string) *application.WebviewWindow {
	switch name {
	case WindowMain:
		if _, ok := wm.winMap[name]; !ok {
			wm.initializeMain()
		}
		return wm.winMap[WindowMain]
	case WindowFloating:
		if _, ok := wm.winMap[name]; !ok {
			wm.initializeFloating()
		}
		return wm.winMap[WindowFloating]
	case WindowSettings:
		if _, ok := wm.winMap[name]; !ok {
			// wm.initializeSettings()
		}
		return wm.winMap[WindowSettings]
	default:
		return nil
	}
}

func (wm *WindowManager) ShowFloating() {
	wm.GetWindow(WindowMain).Hide()
	wm.GetWindow(WindowFloating).Show()
}

func (wm *WindowManager) ShowMain() {
	wm.GetWindow(WindowFloating).Hide()
	wm.GetWindow(WindowMain).Show()
}
