package app

import (
	"changeme/background/database"
	"changeme/background/util"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

const (
	snapShow = 10 // 窗口贴边显示出来的像素
)

type WindowManager struct {
	app                *application.App
	winMap             map[string]*application.WebviewWindow
	floatingStickySide int // 浮动窗口贴边位置：0-无贴边，1-左，2-右，3-上

	initializeSuccess bool // 初始化是否成功
}

func (wm *WindowManager) Run() error {
	wm.GetWindow(WindowMain)
	return wm.app.Run()
}

func (wm *WindowManager) initializeMain() {
	win := wm.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "点点小助理",
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
		DisableResize:    true,
		Windows:          application.WindowsWindow{
			// 	ExStyle: w32.WS_EX_TOOLWINDOW,
		},
	})

	wm.winMap[WindowMain] = win
}

func (wm *WindowManager) initializeFloating() {
	win := wm.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "点点飘啊飘",
		Frameless:      true, // 无边框窗口
		Width:          80,
		Height:         80,
		BackgroundType: application.BackgroundTypeTransparent,
		URL:            "/floating",
		AlwaysOnTop:    true,
		DisableResize:  true,
		Windows: application.WindowsWindow{
			ExStyle:         w32.WS_EX_LAYERED | w32.WS_EX_TOOLWINDOW | w32.WS_EX_TOPMOST,
			HiddenOnTaskbar: true,
		},
	})

	win.RegisterHook(events.Common.WindowDidMove, func(event *application.WindowEvent) {
		rect := win.Bounds()
		// 判断窗口是否靠近屏幕边缘
		screenWidth := w32.GetSystemMetrics(w32.SM_CXSCREEN)
		edgeThreshold := 40 // 靠近边缘的阈值，单位为像素

		if rect.X <= edgeThreshold {
			// 左侧
			wm.floatingStickySide = 1
		} else if rect.X+rect.Width >= screenWidth-edgeThreshold {
			// 右侧
			wm.floatingStickySide = 2
		} else if rect.Y <= edgeThreshold {
			// 上方
			wm.floatingStickySide = 3
		} else {
			wm.floatingStickySide = 0
		}
	})

	// 鼠标进入时显示窗口，这里只能用前端传入的自定义事件完成
	wm.app.Event.On("mouse-enter-floating", func(event *application.CustomEvent) {
		switch wm.floatingStickySide {
		case 1:
			rect := win.Bounds()
			win.SetPosition(0, rect.Y)
		case 2:
			rect := win.Bounds()
			screenWidth := w32.GetSystemMetrics(w32.SM_CXSCREEN)
			win.SetPosition(screenWidth-rect.Width, rect.Y)
		case 3:
			rect := win.Bounds()
			win.SetPosition(rect.X, 0)
		}
	})

	// 鼠标离开时隐藏窗口
	wm.app.Event.On("mouse-leave-floating", func(event *application.CustomEvent) {
		switch wm.floatingStickySide {
		case 1:
			rect := win.Bounds()
			win.SetPosition(0-rect.Width+snapShow, rect.Y)
		case 2:
			rect := win.Bounds()
			screenWidth := w32.GetSystemMetrics(w32.SM_CXSCREEN)
			win.SetPosition(screenWidth-snapShow, rect.Y)
		case 3:
			rect := win.Bounds()
			win.SetPosition(rect.X, 0-rect.Height+snapShow)
		}
	})

	wm.winMap[WindowFloating] = win
}

func (wm *WindowManager) initializeSettings() {
	win := wm.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "点点配置项",
		Frameless: true, // 无边框窗口
		Width:     600,
		Height:    500,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/settings",
		DisableResize:    true,
		Windows:          application.WindowsWindow{
			// 	ExStyle: w32.WS_EX_TOOLWINDOW,
		},
	})

	wm.winMap[WindowSettings] = win
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
			wm.initializeSettings()
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

func (wm *WindowManager) ShowSettings() {
	settingsWindow := wm.GetWindow(WindowSettings)
	if settingsWindow.IsVisible() {
		settingsWindow.Focus()
		return
	}
	settingsWindow.Show()
}

func (wm *WindowManager) HideSettings() {
	wm.GetWindow(WindowSettings).Hide()
}

// 右键菜单
func (wm *WindowManager) buildContextMenu() {
	contextMenu := application.NewContextMenu("floating-context-menu")

	click2ShowMain := contextMenu.Add("显示主界面")
	click2ShowMain.OnClick(func(ctx *application.Context) {
		wm.ShowMain()
	})
	click2Close := contextMenu.Add("关闭")
	click2Close.OnClick(func(ctx *application.Context) {
		wm.app.Quit()
	})

}

// 发送事件
func (wm *WindowManager) EmitEvent(name string, data any) {
	wm.app.Event.EmitEvent(&application.CustomEvent{
		Name: name,
		Data: data,
	})
}

// 程序启动时调用
func (wm *WindowManager) OnAppStart() {
	if err := util.InitializeSnowflake(); err != nil {
		return
	}
	if err := database.Initialize(); err != nil {
		return
	}
	wm.initializeSuccess = true
}

func (wm *WindowManager) IsInitializeSuccess() bool {
	return wm.initializeSuccess
}
