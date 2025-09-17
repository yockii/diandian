package engine

import (
	"fmt"
	"time"

	"diandian/background/automation/app"
	"diandian/background/automation/core"
	"diandian/background/automation/file"
	"diandian/background/automation/keyboard"
	"diandian/background/automation/mouse"
	"diandian/background/automation/screen"
	"diandian/background/automation/system"
)

// Engine 自动化引擎实现
type Engine struct {
	mouse    *mouse.Mouse
	keyboard *keyboard.Keyboard
	app      *app.Launcher
	file     *file.Operator
	screen   *screen.Screen
	system   *system.System
	
	initialized bool
}

// NewEngine 创建自动化引擎实例
func NewEngine() *Engine {
	return &Engine{
		mouse:    mouse.NewMouse(),
		keyboard: keyboard.NewKeyboard(),
		app:      app.NewLauncher(),
		file:     file.NewOperator(),
		screen:   screen.NewScreen(),
		system:   system.NewSystem(),
	}
}

// Initialize 初始化引擎
func (e *Engine) Initialize() *core.OperationResult {
	start := time.Now()
	
	if e.initialized {
		result := core.NewSuccessResult("自动化引擎已初始化", nil)
		result.SetDuration(start)
		return result
	}
	
	// 这里可以添加初始化逻辑，比如检查权限、加载配置等
	e.initialized = true
	
	result := core.NewSuccessResult("自动化引擎初始化成功", nil)
	result.SetDuration(start)
	return result
}

// Cleanup 清理资源
func (e *Engine) Cleanup() *core.OperationResult {
	start := time.Now()
	
	e.initialized = false
	
	result := core.NewSuccessResult("自动化引擎清理完成", nil)
	result.SetDuration(start)
	return result
}

// Wait 等待指定时间
func (e *Engine) Wait(duration int) *core.OperationResult {
	start := time.Now()
	
	time.Sleep(time.Duration(duration) * time.Millisecond)
	
	result := core.NewSuccessResult(
		fmt.Sprintf("等待 %d 毫秒", duration),
		map[string]interface{}{
			"duration_ms": duration,
		},
	)
	result.SetDuration(start)
	return result
}

// 鼠标操作方法
func (e *Engine) Click(x, y int, button core.MouseButton) *core.OperationResult {
	return e.mouse.Click(x, y, button)
}

func (e *Engine) DoubleClick(x, y int) *core.OperationResult {
	return e.mouse.DoubleClick(x, y)
}

func (e *Engine) RightClick(x, y int) *core.OperationResult {
	return e.mouse.RightClick(x, y)
}

func (e *Engine) Drag(fromX, fromY, toX, toY int) *core.OperationResult {
	return e.mouse.Drag(fromX, fromY, toX, toY)
}

func (e *Engine) Move(x, y int) *core.OperationResult {
	return e.mouse.Move(x, y)
}

func (e *Engine) GetPosition() (*core.Point, *core.OperationResult) {
	return e.mouse.GetPosition()
}

func (e *Engine) Scroll(x, y int, direction string, clicks int) *core.OperationResult {
	return e.mouse.Scroll(x, y, direction, clicks)
}

// 键盘操作方法
func (e *Engine) Type(text string) *core.OperationResult {
	return e.keyboard.Type(text)
}

func (e *Engine) KeyPress(key string) *core.OperationResult {
	return e.keyboard.KeyPress(key)
}

func (e *Engine) KeyDown(key string) *core.OperationResult {
	return e.keyboard.KeyDown(key)
}

func (e *Engine) KeyUp(key string) *core.OperationResult {
	return e.keyboard.KeyUp(key)
}

func (e *Engine) Hotkey(modifiers []core.KeyModifier, key string) *core.OperationResult {
	return e.keyboard.Hotkey(modifiers, key)
}

func (e *Engine) Copy() *core.OperationResult {
	return e.keyboard.Copy()
}

func (e *Engine) Paste() *core.OperationResult {
	return e.keyboard.Paste()
}

func (e *Engine) SelectAll() *core.OperationResult {
	return e.keyboard.SelectAll()
}

// 应用程序操作方法
func (e *Engine) Launch(appName string) *core.OperationResult {
	return e.app.Launch(appName)
}

func (e *Engine) LaunchWithPath(path string, args ...string) *core.OperationResult {
	return e.app.LaunchWithPath(path, args...)
}

func (e *Engine) LaunchApp(app *core.AppInfo) *core.OperationResult {
	return e.app.LaunchApp(app)
}

func (e *Engine) GetInstalledApps() ([]*core.AppInfo, *core.OperationResult) {
	return e.app.GetInstalledApps()
}

func (e *Engine) FindApp(name string) (*core.AppInfo, *core.OperationResult) {
	return e.app.FindApp(name)
}

// 文件操作方法
func (e *Engine) CreateFile(path string, content []byte) *core.OperationResult {
	return e.file.CreateFile(path, content)
}

func (e *Engine) CreateDir(path string) *core.OperationResult {
	return e.file.CreateDir(path)
}

func (e *Engine) MoveFile(src, dst string) *core.OperationResult {
	return e.file.MoveFile(src, dst)
}

func (e *Engine) CopyFile(src, dst string) *core.OperationResult {
	return e.file.CopyFile(src, dst)
}

func (e *Engine) DeleteFile(path string) *core.OperationResult {
	return e.file.DeleteFile(path)
}

func (e *Engine) DeleteDir(path string) *core.OperationResult {
	return e.file.DeleteDir(path)
}

func (e *Engine) RenameFile(oldPath, newPath string) *core.OperationResult {
	return e.file.RenameFile(oldPath, newPath)
}

func (e *Engine) FileExists(path string) (bool, *core.OperationResult) {
	return e.file.FileExists(path)
}

func (e *Engine) GetFileInfo(path string) (interface{}, *core.OperationResult) {
	return e.file.GetFileInfo(path)
}

func (e *Engine) ListDir(path string) ([]string, *core.OperationResult) {
	return e.file.ListDir(path)
}

// 屏幕操作方法
func (e *Engine) Screenshot() ([]byte, *core.OperationResult) {
	return e.screen.Screenshot()
}

func (e *Engine) ScreenshotArea(rect core.Rect) ([]byte, *core.OperationResult) {
	return e.screen.ScreenshotArea(rect)
}

func (e *Engine) GetScreenSize() (*core.Size, *core.OperationResult) {
	return e.screen.GetScreenSize()
}

func (e *Engine) FindImage(templatePath string) (*core.Point, *core.OperationResult) {
	return e.screen.FindImage(templatePath)
}

func (e *Engine) FindText(text string) (*core.Point, *core.OperationResult) {
	return e.screen.FindText(text)
}

// 系统操作方法
func (e *Engine) GetClipboard() (string, *core.OperationResult) {
	return e.system.GetClipboard()
}

func (e *Engine) SetClipboard(text string) *core.OperationResult {
	return e.system.SetClipboard(text)
}

func (e *Engine) GetActiveWindow() (*core.WindowInfo, *core.OperationResult) {
	return e.system.GetActiveWindow()
}

func (e *Engine) GetWindows() ([]*core.WindowInfo, *core.OperationResult) {
	return e.system.GetWindows()
}

func (e *Engine) ActivateWindow(handle uintptr) *core.OperationResult {
	return e.system.ActivateWindow(handle)
}

func (e *Engine) CloseWindow(handle uintptr) *core.OperationResult {
	return e.system.CloseWindow(handle)
}

func (e *Engine) MinimizeWindow(handle uintptr) *core.OperationResult {
	return e.system.MinimizeWindow(handle)
}

func (e *Engine) MaximizeWindow(handle uintptr) *core.OperationResult {
	return e.system.MaximizeWindow(handle)
}

// 高级操作方法

// ClickAndWait 点击并等待
func (e *Engine) ClickAndWait(x, y int, waitMs int) *core.OperationResult {
	result := e.Click(x, y, core.LeftButton)
	if !result.Success {
		return result
	}
	
	return e.Wait(waitMs)
}

// TypeAndEnter 输入文本并按回车
func (e *Engine) TypeAndEnter(text string) *core.OperationResult {
	result := e.Type(text)
	if !result.Success {
		return result
	}
	
	return e.KeyPress("enter")
}

// SaveScreenshotToFile 保存截屏到文件
func (e *Engine) SaveScreenshotToFile(filePath string) *core.OperationResult {
	return e.screen.SaveScreenshot(filePath)
}

// LaunchURL 打开URL
func (e *Engine) LaunchURL(url string) *core.OperationResult {
	return e.app.LaunchURL(url)
}

// WriteTextFile 写入文本文件
func (e *Engine) WriteTextFile(path, content string) *core.OperationResult {
	return e.file.WriteTextFile(path, content)
}

// ReadTextFile 读取文本文件
func (e *Engine) ReadTextFile(path string) (string, *core.OperationResult) {
	return e.file.ReadTextFile(path)
}
