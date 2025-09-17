package core

// MouseOperator 鼠标操作接口
type MouseOperator interface {
	// Click 点击指定位置
	Click(x, y int, button MouseButton) *OperationResult
	
	// DoubleClick 双击指定位置
	DoubleClick(x, y int) *OperationResult
	
	// RightClick 右键点击指定位置
	RightClick(x, y int) *OperationResult
	
	// Drag 拖拽操作
	Drag(fromX, fromY, toX, toY int) *OperationResult
	
	// Move 移动鼠标到指定位置
	Move(x, y int) *OperationResult
	
	// GetPosition 获取当前鼠标位置
	GetPosition() (*Point, *OperationResult)
	
	// Scroll 滚动操作
	Scroll(x, y int, direction string, clicks int) *OperationResult
}

// KeyboardOperator 键盘操作接口
type KeyboardOperator interface {
	// Type 输入文本
	Type(text string) *OperationResult
	
	// KeyPress 按下并释放按键
	KeyPress(key string) *OperationResult
	
	// KeyDown 按下按键
	KeyDown(key string) *OperationResult
	
	// KeyUp 释放按键
	KeyUp(key string) *OperationResult
	
	// Hotkey 组合键操作
	Hotkey(modifiers []KeyModifier, key string) *OperationResult
	
	// Copy 复制操作 (Ctrl+C)
	Copy() *OperationResult
	
	// Paste 粘贴操作 (Ctrl+V)
	Paste() *OperationResult
	
	// SelectAll 全选操作 (Ctrl+A)
	SelectAll() *OperationResult
}

// AppLauncher 应用程序启动接口
type AppLauncher interface {
	// Launch 启动应用程序
	Launch(appName string) *OperationResult
	
	// LaunchWithPath 通过路径启动应用程序
	LaunchWithPath(path string, args ...string) *OperationResult
	
	// LaunchApp 启动预定义的应用程序
	LaunchApp(app *AppInfo) *OperationResult
	
	// GetInstalledApps 获取已安装的应用程序列表
	GetInstalledApps() ([]*AppInfo, *OperationResult)
	
	// FindApp 查找应用程序
	FindApp(name string) (*AppInfo, *OperationResult)
}

// FileOperator 文件操作接口
type FileOperator interface {
	// CreateFile 创建文件
	CreateFile(path string, content []byte) *OperationResult
	
	// CreateDir 创建目录
	CreateDir(path string) *OperationResult
	
	// MoveFile 移动文件
	MoveFile(src, dst string) *OperationResult
	
	// CopyFile 复制文件
	CopyFile(src, dst string) *OperationResult
	
	// DeleteFile 删除文件
	DeleteFile(path string) *OperationResult
	
	// DeleteDir 删除目录
	DeleteDir(path string) *OperationResult
	
	// RenameFile 重命名文件
	RenameFile(oldPath, newPath string) *OperationResult
	
	// FileExists 检查文件是否存在
	FileExists(path string) (bool, *OperationResult)
	
	// GetFileInfo 获取文件信息
	GetFileInfo(path string) (interface{}, *OperationResult)
	
	// ListDir 列出目录内容
	ListDir(path string) ([]string, *OperationResult)
}

// ScreenOperator 屏幕操作接口
type ScreenOperator interface {
	// Screenshot 截取屏幕
	Screenshot() ([]byte, *OperationResult)
	
	// ScreenshotArea 截取指定区域
	ScreenshotArea(rect Rect) ([]byte, *OperationResult)
	
	// GetScreenSize 获取屏幕尺寸
	GetScreenSize() (*Size, *OperationResult)
	
	// FindImage 在屏幕上查找图像
	FindImage(templatePath string) (*Point, *OperationResult)
	
	// FindText 在屏幕上查找文本（OCR）
	FindText(text string) (*Point, *OperationResult)
}

// SystemOperator 系统操作接口
type SystemOperator interface {
	// GetClipboard 获取剪贴板内容
	GetClipboard() (string, *OperationResult)
	
	// SetClipboard 设置剪贴板内容
	SetClipboard(text string) *OperationResult
	
	// GetActiveWindow 获取当前活动窗口
	GetActiveWindow() (*WindowInfo, *OperationResult)
	
	// GetWindows 获取所有窗口
	GetWindows() ([]*WindowInfo, *OperationResult)
	
	// ActivateWindow 激活指定窗口
	ActivateWindow(handle uintptr) *OperationResult
	
	// CloseWindow 关闭指定窗口
	CloseWindow(handle uintptr) *OperationResult
	
	// MinimizeWindow 最小化窗口
	MinimizeWindow(handle uintptr) *OperationResult
	
	// MaximizeWindow 最大化窗口
	MaximizeWindow(handle uintptr) *OperationResult
}

// AutomationEngine 自动化引擎接口
type AutomationEngine interface {
	MouseOperator
	KeyboardOperator
	AppLauncher
	FileOperator
	ScreenOperator
	SystemOperator
	
	// Initialize 初始化引擎
	Initialize() *OperationResult
	
	// Cleanup 清理资源
	Cleanup() *OperationResult
	
	// Wait 等待指定时间
	Wait(duration int) *OperationResult
}
