package core

import (
	"image"
	"time"
)

// OperationResult 操作结果
type OperationResult struct {
	Success   bool          `json:"success"`
	Message   string        `json:"message"`
	Data      interface{}   `json:"data,omitempty"`
	Error     string        `json:"error,omitempty"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

// Point 坐标点
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Size 尺寸
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Rect 矩形区域
type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// MouseButton 鼠标按键
type MouseButton string

const (
	LeftButton   MouseButton = "left"
	RightButton  MouseButton = "right"
	MiddleButton MouseButton = "middle"
)

// KeyModifier 键盘修饰键
type KeyModifier string

const (
	ModCtrl  KeyModifier = "ctrl"
	ModAlt   KeyModifier = "alt"
	ModShift KeyModifier = "shift"
	ModWin   KeyModifier = "win"
)

// FileOperation 文件操作类型
type FileOperation string

const (
	FileCreate FileOperation = "create"
	FileMove   FileOperation = "move"
	FileCopy   FileOperation = "copy"
	FileDelete FileOperation = "delete"
	FileRename FileOperation = "rename"
)

// AppInfo 应用程序信息
type AppInfo struct {
	Name        string            `json:"name"`
	DisplayName string            `json:"display_name"`
	Path        string            `json:"path"`
	Args        []string          `json:"args,omitempty"`
	WorkDir     string            `json:"work_dir,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
}

// WindowInfo 窗口信息
type WindowInfo struct {
	Title  string  `json:"title"`
	Class  string  `json:"class"`
	PID    int     `json:"pid"`
	Handle uintptr `json:"handle"`
	Rect   Rect    `json:"rect"`
}

// ScreenInfo 屏幕信息
type ScreenInfo struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	DPI    int `json:"dpi"`
}

// DisplayCapture 显示器截图信息
type DisplayCapture struct {
	Index     int             `json:"index"`     // 显示器索引
	Bounds    image.Rectangle `json:"bounds"`    // 显示器边界
	ImageData []byte          `json:"-"`         // 图像数据（不序列化）
	Width     int             `json:"width"`     // 宽度
	Height    int             `json:"height"`    // 高度
	IsActive  bool            `json:"is_active"` // 是否是活动显示器
}

// NewResult 创建操作结果
func NewResult(success bool, message string) *OperationResult {
	return &OperationResult{
		Success:   success,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewSuccessResult 创建成功结果
func NewSuccessResult(message string, data interface{}) *OperationResult {
	return &OperationResult{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewErrorResult 创建错误结果
func NewErrorResult(message string, err error) *OperationResult {
	result := &OperationResult{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	}
	if err != nil {
		result.Error = err.Error()
	}
	return result
}

// SetDuration 设置执行时间
func (r *OperationResult) SetDuration(start time.Time) {
	r.Duration = time.Since(start)
}
