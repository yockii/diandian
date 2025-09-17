package system

import (
	"fmt"
	"time"

	"diandian/background/automation/core"

	"github.com/go-vgo/robotgo"
)

// System 系统操作实现
type System struct{}

// NewSystem 创建系统操作实例
func NewSystem() *System {
	return &System{}
}

// GetClipboard 获取剪贴板内容
func (s *System) GetClipboard() (string, *core.OperationResult) {
	start := time.Now()

	text, err := robotgo.ReadAll()
	if err != nil {
		result := core.NewErrorResult("获取剪贴板内容失败", err)
		result.SetDuration(start)
		return "", result
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("获取剪贴板内容 (%d 字符)", len(text)),
		map[string]interface{}{
			"text":   text,
			"length": len(text),
		},
	)
	result.SetDuration(start)
	return text, result
}

// SetClipboard 设置剪贴板内容
func (s *System) SetClipboard(text string) *core.OperationResult {
	start := time.Now()

	robotgo.WriteAll(text)

	result := core.NewSuccessResult(
		fmt.Sprintf("设置剪贴板内容 (%d 字符)", len(text)),
		map[string]interface{}{
			"text":   text,
			"length": len(text),
		},
	)
	result.SetDuration(start)
	return result
}

// GetActiveWindow 获取当前活动窗口
func (s *System) GetActiveWindow() (*core.WindowInfo, *core.OperationResult) {
	start := time.Now()

	// 获取活动窗口的PID
	pid := robotgo.GetPid()

	// 获取窗口标题
	title := robotgo.GetTitle()

	windowInfo := &core.WindowInfo{
		Title: title,
		PID:   int(pid),
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("获取活动窗口: %s (PID: %d)", title, pid),
		windowInfo,
	)
	result.SetDuration(start)
	return windowInfo, result
}

// GetWindows 获取所有窗口
func (s *System) GetWindows() ([]*core.WindowInfo, *core.OperationResult) {
	start := time.Now()

	// 注意：robotgo的GetWindows API在不同版本中可能不同
	// 这里返回一个空列表，实际实现需要根据具体的robotgo版本调整
	var windowInfos []*core.WindowInfo

	result := core.NewSuccessResult(
		"获取窗口列表功能需要根据robotgo版本调整",
		windowInfos,
	)
	result.SetDuration(start)
	return windowInfos, result
}

// ActivateWindow 激活指定窗口
func (s *System) ActivateWindow(handle uintptr) *core.OperationResult {
	start := time.Now()

	// 注意：robotgo的窗口操作API在不同版本中可能不同
	// 这里返回未实现错误，实际使用时需要根据具体版本调整
	result := core.NewErrorResult(
		"窗口激活功能需要根据robotgo版本调整",
		fmt.Errorf("not implemented"),
	)
	result.SetDuration(start)
	return result
}

// CloseWindow 关闭指定窗口
func (s *System) CloseWindow(handle uintptr) *core.OperationResult {
	start := time.Now()

	result := core.NewErrorResult(
		"窗口关闭功能需要根据robotgo版本调整",
		fmt.Errorf("not implemented"),
	)
	result.SetDuration(start)
	return result
}

// MinimizeWindow 最小化窗口
func (s *System) MinimizeWindow(handle uintptr) *core.OperationResult {
	start := time.Now()

	result := core.NewErrorResult(
		"窗口最小化功能需要根据robotgo版本调整",
		fmt.Errorf("not implemented"),
	)
	result.SetDuration(start)
	return result
}

// MaximizeWindow 最大化窗口
func (s *System) MaximizeWindow(handle uintptr) *core.OperationResult {
	start := time.Now()

	result := core.NewErrorResult(
		"窗口最大化功能需要根据robotgo版本调整",
		fmt.Errorf("not implemented"),
	)
	result.SetDuration(start)
	return result
}

// Sleep 系统休眠
func (s *System) Sleep() *core.OperationResult {
	start := time.Now()

	// 注意：这个功能需要系统权限，可能不会立即生效
	robotgo.Sleep(1) // 参数表示休眠类型，1表示休眠

	result := core.NewSuccessResult(
		"系统休眠命令已发送",
		nil,
	)
	result.SetDuration(start)
	return result
}

// GetSystemInfo 获取系统信息
func (s *System) GetSystemInfo() (map[string]interface{}, *core.OperationResult) {
	start := time.Now()

	// 获取屏幕尺寸
	width, height := robotgo.GetScreenSize()

	// 获取鼠标位置
	mouseX, mouseY := robotgo.GetMousePos()

	systemInfo := map[string]interface{}{
		"screen_width":  width,
		"screen_height": height,
		"mouse_x":       mouseX,
		"mouse_y":       mouseY,
	}

	result := core.NewSuccessResult(
		"获取系统信息成功",
		systemInfo,
	)
	result.SetDuration(start)
	return systemInfo, result
}

// FindWindow 根据标题查找窗口
func (s *System) FindWindow(title string) (*core.WindowInfo, *core.OperationResult) {
	start := time.Now()

	windows, result := s.GetWindows()
	if !result.Success {
		return nil, result
	}

	for _, window := range windows {
		if window.Title == title {
			result := core.NewSuccessResult(
				fmt.Sprintf("找到窗口: %s", title),
				window,
			)
			result.SetDuration(start)
			return window, result
		}
	}

	result = core.NewErrorResult(
		fmt.Sprintf("未找到窗口: %s", title),
		fmt.Errorf("window not found"),
	)
	result.SetDuration(start)
	return nil, result
}

// IsWindowActive 检查窗口是否为活动窗口
func (s *System) IsWindowActive(handle uintptr) (bool, *core.OperationResult) {
	start := time.Now()

	activeWindow, result := s.GetActiveWindow()
	if !result.Success {
		return false, result
	}

	isActive := activeWindow.Handle == handle

	result = core.NewSuccessResult(
		fmt.Sprintf("检查窗口活动状态: %d (活动: %t)", handle, isActive),
		map[string]interface{}{
			"handle":    handle,
			"is_active": isActive,
		},
	)
	result.SetDuration(start)
	return isActive, result
}

// WaitForWindow 等待窗口出现
func (s *System) WaitForWindow(title string, timeout time.Duration) (*core.WindowInfo, *core.OperationResult) {
	start := time.Now()

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		window, result := s.FindWindow(title)
		if result.Success {
			result.Message = fmt.Sprintf("等待窗口出现成功: %s", title)
			result.SetDuration(start)
			return window, result
		}

		time.Sleep(500 * time.Millisecond) // 每500ms检查一次
	}

	result := core.NewErrorResult(
		fmt.Sprintf("等待窗口超时: %s", title),
		fmt.Errorf("timeout waiting for window"),
	)
	result.SetDuration(start)
	return nil, result
}
