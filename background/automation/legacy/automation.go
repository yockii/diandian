// Package automation 提供桌面自动化操作功能
//
// 这个包包含了完整的桌面自动化操作能力，包括：
// - 鼠标操作：点击、拖拽、移动、滚动等
// - 键盘操作：输入文本、按键、组合键等
// - 应用程序启动：启动常见应用程序、打开URL等
// - 文件操作：创建、移动、复制、删除文件和目录
// - 屏幕操作：截屏、图像识别等
// - 系统操作：剪贴板、窗口管理等
//
// 使用示例：
//
//	import "diandian/background/automation"
//
//	// 创建自动化引擎
//	engine := automation.NewEngine()
//	engine.Initialize()
//
//	// 鼠标操作
//	engine.Click(100, 100, automation.LeftButton)
//	engine.Type("Hello World!")
//
//	// 应用程序操作
//	engine.Launch("notepad")
//
//	// 文件操作
//	engine.CreateFile("test.txt", []byte("Hello"))
//
//	// 清理
//	engine.Cleanup()
package automation

import (
	"diandian/background/automation/core"
	"diandian/background/automation/engine"
)

// 导出核心类型
type (
	// OperationResult 操作结果
	OperationResult = core.OperationResult

	// Point 坐标点
	Point = core.Point

	// Size 尺寸
	Size = core.Size

	// Rect 矩形区域
	Rect = core.Rect

	// MouseButton 鼠标按键
	MouseButton = core.MouseButton

	// KeyModifier 键盘修饰键
	KeyModifier = core.KeyModifier

	// AppInfo 应用程序信息
	AppInfo = core.AppInfo

	// WindowInfo 窗口信息
	WindowInfo = core.WindowInfo

	// ScreenInfo 屏幕信息
	ScreenInfo = core.ScreenInfo

	// AutomationEngine 自动化引擎接口
	AutomationEngine = core.AutomationEngine
)

// 导出常量
const (
	// 鼠标按键
	LeftButton   = core.LeftButton
	RightButton  = core.RightButton
	MiddleButton = core.MiddleButton

	// 键盘修饰键
	ModCtrl  = core.ModCtrl
	ModAlt   = core.ModAlt
	ModShift = core.ModShift
	ModWin   = core.ModWin
)

// NewEngine 创建新的自动化引擎实例
//
// 返回一个完全初始化的自动化引擎，包含所有操作模块。
// 使用前需要调用 Initialize() 方法进行初始化。
//
// 示例：
//
//	engine := automation.NewEngine()
//	result := engine.Initialize()
//	if !result.Success {
//		log.Fatal("初始化失败:", result.Error)
//	}
//	defer engine.Cleanup()
func NewEngine() AutomationEngine {
	return engine.NewEngine()
}

// NewResult 创建操作结果
func NewResult(success bool, message string) *OperationResult {
	return core.NewResult(success, message)
}

// NewSuccessResult 创建成功结果
func NewSuccessResult(message string, data interface{}) *OperationResult {
	return core.NewSuccessResult(message, data)
}

// NewErrorResult 创建错误结果
func NewErrorResult(message string, err error) *OperationResult {
	return core.NewErrorResult(message, err)
}

// 便捷函数

// QuickClick 快速点击指定位置
func QuickClick(x, y int) *OperationResult {
	engine := NewEngine()
	engine.Initialize()
	defer engine.Cleanup()

	return engine.Click(x, y, LeftButton)
}

// QuickType 快速输入文本
func QuickType(text string) *OperationResult {
	engine := NewEngine()
	engine.Initialize()
	defer engine.Cleanup()

	return engine.Type(text)
}

// QuickLaunch 快速启动应用程序
func QuickLaunch(appName string) *OperationResult {
	engine := NewEngine()
	engine.Initialize()
	defer engine.Cleanup()

	return engine.Launch(appName)
}

// QuickScreenshot 快速截屏
func QuickScreenshot() ([]byte, *OperationResult) {
	engine := NewEngine()
	engine.Initialize()
	defer engine.Cleanup()

	return engine.Screenshot()
}

// QuickCreateFile 快速创建文件
func QuickCreateFile(path string, content []byte) *OperationResult {
	engine := NewEngine()
	engine.Initialize()
	defer engine.Cleanup()

	return engine.CreateFile(path, content)
}

// QuickGetClipboard 快速获取剪贴板内容
func QuickGetClipboard() (string, *OperationResult) {
	engine := NewEngine()
	engine.Initialize()
	defer engine.Cleanup()

	return engine.GetClipboard()
}

// QuickSetClipboard 快速设置剪贴板内容
func QuickSetClipboard(text string) *OperationResult {
	engine := NewEngine()
	engine.Initialize()
	defer engine.Cleanup()

	return engine.SetClipboard(text)
}

// 预定义的常用操作序列

// OpenNotepadAndType 打开记事本并输入文本
func OpenNotepadAndType(text string) *OperationResult {
	engine := NewEngine()
	engine.Initialize()
	defer engine.Cleanup()

	// 启动记事本
	result := engine.Launch("notepad")
	if !result.Success {
		return result
	}

	// 等待应用程序启动
	engine.Wait(2000)

	// 输入文本
	return engine.Type(text)
}

// TakeScreenshotAndSave 截屏并保存到文件
func TakeScreenshotAndSave(filePath string) *OperationResult {
	engine := NewEngine()
	engine.Initialize()
	defer engine.Cleanup()

	// 截屏
	imageData, result := engine.Screenshot()
	if !result.Success {
		return result
	}

	// 保存到文件
	return QuickCreateFile(filePath, imageData)
}

// CopyTextToClipboard 复制文本到剪贴板
func CopyTextToClipboard(text string) *OperationResult {
	engine := NewEngine()
	engine.Initialize()
	defer engine.Cleanup()

	// 设置剪贴板内容
	result := engine.SetClipboard(text)
	if !result.Success {
		return result
	}

	// 验证设置是否成功
	clipText, result := engine.GetClipboard()
	if !result.Success {
		return result
	}

	if clipText != text {
		return NewErrorResult("剪贴板内容验证失败", nil)
	}

	return NewSuccessResult("文本已复制到剪贴板", map[string]interface{}{
		"text": text,
	})
}

// CreateTextFileWithContent 创建包含指定内容的文本文件
func CreateTextFileWithContent(filePath, content string) *OperationResult {
	engine := NewEngine()
	engine.Initialize()
	defer engine.Cleanup()

	return engine.CreateFile(filePath, []byte(content))
}

// GetSystemInfo 获取系统信息
func GetSystemInfo() (map[string]interface{}, *OperationResult) {
	engine := NewEngine()
	engine.Initialize()
	defer engine.Cleanup()

	// 获取屏幕尺寸
	screenSize, result := engine.GetScreenSize()
	if !result.Success {
		return nil, result
	}

	// 获取鼠标位置
	mousePos, result := engine.GetPosition()
	if !result.Success {
		return nil, result
	}

	// 获取活动窗口
	activeWindow, result := engine.GetActiveWindow()
	if !result.Success {
		return nil, result
	}

	// 获取剪贴板内容
	clipboardText, _ := engine.GetClipboard()

	systemInfo := map[string]interface{}{
		"screen_size":    screenSize,
		"mouse_position": mousePos,
		"active_window":  activeWindow,
		"clipboard_text": clipboardText,
	}

	return systemInfo, NewSuccessResult("获取系统信息成功", systemInfo)
}
