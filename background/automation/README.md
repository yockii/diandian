# Automation Package

这是一个功能完整的桌面自动化操作包，提供了鼠标、键盘、应用程序启动、文件操作、屏幕操作和系统操作等功能。

## 功能特性

### 1. 鼠标操作
- 点击（左键、右键、中键）
- 双击
- 拖拽
- 移动
- 滚动
- 获取鼠标位置

### 2. 键盘操作
- 文本输入
- 按键操作
- 组合键（Ctrl+C、Ctrl+V等）
- 特殊键（回车、ESC、方向键等）

### 3. 应用程序启动
- 启动常见应用程序（记事本、计算器、浏览器等）
- 通过路径启动自定义应用程序
- 打开URL
- 获取已安装应用程序列表

### 4. 文件操作
- 创建文件和目录
- 移动、复制、删除文件
- 重命名文件
- 读取文件内容
- 列出目录内容

### 5. 屏幕操作
- 截屏（全屏或指定区域）
- 图像识别
- 获取屏幕尺寸
- 保存截屏到文件

### 6. 系统操作
- 剪贴板操作
- 窗口管理（激活、关闭、最小化、最大化）
- 获取系统信息

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "diandian/background/automation"
)

func main() {
    // 创建自动化引擎
    engine := automation.NewEngine()
    
    // 初始化
    result := engine.Initialize()
    if !result.Success {
        fmt.Printf("初始化失败: %s\n", result.Error)
        return
    }
    defer engine.Cleanup()
    
    // 鼠标操作
    engine.Click(100, 100, automation.LeftButton)
    
    // 键盘操作
    engine.Type("Hello, World!")
    
    // 应用程序操作
    engine.Launch("notepad")
    
    // 文件操作
    engine.CreateFile("test.txt", []byte("Hello"))
    
    // 截屏
    imageData, result := engine.Screenshot()
    if result.Success {
        fmt.Printf("截屏成功，大小: %d 字节\n", len(imageData))
    }
}
```

### 便捷函数

```go
// 快速操作，无需手动管理引擎生命周期
automation.QuickClick(100, 100)
automation.QuickType("Hello")
automation.QuickLaunch("calculator")

// 复合操作
automation.OpenNotepadAndType("这是测试文本")
automation.TakeScreenshotAndSave("screenshot.png")
automation.CopyTextToClipboard("复制的文本")
```

## 详细示例

### 鼠标操作示例

```go
engine := automation.NewEngine()
engine.Initialize()
defer engine.Cleanup()

// 点击
engine.Click(100, 100, automation.LeftButton)
engine.RightClick(200, 200)
engine.DoubleClick(300, 300)

// 拖拽
engine.Drag(100, 100, 200, 200)

// 滚动
engine.Scroll(400, 400, "up", 3)

// 获取鼠标位置
pos, result := engine.GetPosition()
if result.Success {
    fmt.Printf("鼠标位置: (%d, %d)\n", pos.X, pos.Y)
}
```

### 键盘操作示例

```go
engine := automation.NewEngine()
engine.Initialize()
defer engine.Cleanup()

// 输入文本
engine.Type("Hello, Automation!")

// 组合键
engine.Hotkey([]automation.KeyModifier{automation.ModCtrl}, "a") // Ctrl+A
engine.Copy()  // Ctrl+C
engine.Paste() // Ctrl+V

// 特殊键
engine.KeyPress("enter")
engine.KeyPress("escape")
engine.KeyPress("tab")
```

### 应用程序操作示例

```go
engine := automation.NewEngine()
engine.Initialize()
defer engine.Cleanup()

// 启动预定义应用程序
engine.Launch("notepad")
engine.Launch("calculator")
engine.Launch("browser")

// 启动自定义应用程序
engine.LaunchWithPath("C:\\MyApp\\app.exe", "--arg1", "--arg2")

// 打开URL
engine.LaunchURL("https://www.example.com")

// 获取已安装应用程序
apps, result := engine.GetInstalledApps()
if result.Success {
    for _, app := range apps {
        fmt.Printf("应用程序: %s (%s)\n", app.DisplayName, app.Path)
    }
}
```

### 文件操作示例

```go
engine := automation.NewEngine()
engine.Initialize()
defer engine.Cleanup()

// 创建目录
engine.CreateDir("test_folder")

// 创建文件
engine.CreateFile("test_folder/test.txt", []byte("Hello, World!"))

// 复制文件
engine.CopyFile("test_folder/test.txt", "test_folder/test_copy.txt")

// 移动文件
engine.MoveFile("test_folder/test_copy.txt", "test_folder/moved.txt")

// 读取文件
content, result := engine.ReadTextFile("test_folder/test.txt")
if result.Success {
    fmt.Printf("文件内容: %s\n", content)
}

// 列出目录
files, result := engine.ListDir("test_folder")
if result.Success {
    fmt.Printf("目录包含 %d 个文件\n", len(files))
}

// 删除文件和目录
engine.DeleteFile("test_folder/test.txt")
engine.DeleteDir("test_folder")
```

### 屏幕操作示例

```go
engine := automation.NewEngine()
engine.Initialize()
defer engine.Cleanup()

// 全屏截图
imageData, result := engine.Screenshot()
if result.Success {
    fmt.Printf("截屏成功，大小: %d 字节\n", len(imageData))
}

// 区域截图
rect := automation.Rect{X: 100, Y: 100, Width: 400, Height: 300}
imageData, result = engine.ScreenshotArea(rect)

// 保存截图到文件
engine.SaveScreenshotToFile("screenshot.png")

// 获取屏幕尺寸
size, result := engine.GetScreenSize()
if result.Success {
    fmt.Printf("屏幕尺寸: %dx%d\n", size.Width, size.Height)
}
```

### 系统操作示例

```go
engine := automation.NewEngine()
engine.Initialize()
defer engine.Cleanup()

// 剪贴板操作
engine.SetClipboard("Hello, Clipboard!")
text, result := engine.GetClipboard()
if result.Success {
    fmt.Printf("剪贴板内容: %s\n", text)
}

// 窗口操作
windows, result := engine.GetWindows()
if result.Success {
    fmt.Printf("找到 %d 个窗口\n", len(windows))
    for _, window := range windows {
        fmt.Printf("窗口: %s\n", window.Title)
    }
}

// 获取活动窗口
activeWindow, result := engine.GetActiveWindow()
if result.Success {
    fmt.Printf("活动窗口: %s\n", activeWindow.Title)
}
```

## 预定义应用程序

包中预定义了以下常见应用程序：

### Windows
- `notepad` - 记事本
- `calculator` - 计算器
- `browser` - 默认浏览器
- `chrome` - Chrome浏览器
- `firefox` - Firefox浏览器
- `explorer` - 文件资源管理器
- `cmd` - 命令提示符
- `powershell` - PowerShell
- `word` - Microsoft Word
- `excel` - Microsoft Excel
- `wechat` - 微信
- `qq` - QQ
- `dingtalk` - 钉钉

### macOS
- `textedit` - 文本编辑
- `calculator` - 计算器
- `safari` - Safari浏览器
- `chrome` - Chrome浏览器
- `finder` - 访达

## 错误处理

所有操作都返回 `OperationResult` 结构，包含以下信息：

```go
type OperationResult struct {
    Success   bool        // 操作是否成功
    Message   string      // 操作描述信息
    Data      interface{} // 返回的数据（可选）
    Error     string      // 错误信息（失败时）
    Duration  time.Duration // 操作耗时
    Timestamp time.Time   // 操作时间戳
}
```

使用示例：

```go
result := engine.Click(100, 100, automation.LeftButton)
if !result.Success {
    fmt.Printf("操作失败: %s\n", result.Error)
    return
}
fmt.Printf("操作成功: %s (耗时: %v)\n", result.Message, result.Duration)
```

## 注意事项

1. **权限要求**: 某些操作可能需要管理员权限
2. **安全性**: 自动化操作具有潜在风险，请谨慎使用
3. **性能**: 频繁的操作可能影响系统性能
4. **兼容性**: 主要支持Windows，部分功能支持macOS和Linux
5. **测试**: 建议在测试环境中验证操作序列

## 依赖项

- `github.com/go-vgo/robotgo` - 跨平台GUI自动化
- `github.com/kbinani/screenshot` - 屏幕截图

## 许可证

本项目采用MIT许可证。
