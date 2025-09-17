package app

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"diandian/background/automation/core"
)

// Launcher 应用程序启动器
type Launcher struct {
	predefinedApps map[string]*core.AppInfo
}

// NewLauncher 创建应用程序启动器
func NewLauncher() *Launcher {
	launcher := &Launcher{
		predefinedApps: make(map[string]*core.AppInfo),
	}
	launcher.initPredefinedApps()
	return launcher
}

// initPredefinedApps 初始化预定义应用程序
func (l *Launcher) initPredefinedApps() {
	if runtime.GOOS == "windows" {
		l.predefinedApps = map[string]*core.AppInfo{
			"notepad": {
				Name:        "notepad",
				DisplayName: "记事本",
				Path:        "notepad.exe",
			},
			"calculator": {
				Name:        "calculator",
				DisplayName: "计算器",
				Path:        "calc.exe",
			},
			"browser": {
				Name:        "browser",
				DisplayName: "浏览器",
				Path:        "msedge.exe", // 默认使用Edge
			},
			"chrome": {
				Name:        "chrome",
				DisplayName: "Chrome浏览器",
				Path:        `C:\Program Files\Google\Chrome\Application\chrome.exe`,
			},
			"firefox": {
				Name:        "firefox",
				DisplayName: "Firefox浏览器",
				Path:        `C:\Program Files\Mozilla Firefox\firefox.exe`,
			},
			"explorer": {
				Name:        "explorer",
				DisplayName: "文件资源管理器",
				Path:        "explorer.exe",
			},
			"cmd": {
				Name:        "cmd",
				DisplayName: "命令提示符",
				Path:        "cmd.exe",
			},
			"powershell": {
				Name:        "powershell",
				DisplayName: "PowerShell",
				Path:        "powershell.exe",
			},
			"word": {
				Name:        "word",
				DisplayName: "Microsoft Word",
				Path:        `C:\Program Files\Microsoft Office\root\Office16\WINWORD.EXE`,
			},
			"excel": {
				Name:        "excel",
				DisplayName: "Microsoft Excel",
				Path:        `C:\Program Files\Microsoft Office\root\Office16\EXCEL.EXE`,
			},
			"wechat": {
				Name:        "wechat",
				DisplayName: "微信",
				Path:        `C:\Program Files\Tencent\WeChat\WeChat.exe`,
			},
			"qq": {
				Name:        "qq",
				DisplayName: "QQ",
				Path:        `C:\Program Files\Tencent\QQ\Bin\QQScLauncher.exe`,
			},
			"dingtalk": {
				Name:        "dingtalk",
				DisplayName: "钉钉",
				Path:        `C:\Users\%USERNAME%\AppData\Local\DingTalk\DingtalkLauncher.exe`,
			},
		}
	} else if runtime.GOOS == "darwin" {
		// macOS应用程序
		l.predefinedApps = map[string]*core.AppInfo{
			"textedit": {
				Name:        "textedit",
				DisplayName: "文本编辑",
				Path:        "/Applications/TextEdit.app",
			},
			"calculator": {
				Name:        "calculator",
				DisplayName: "计算器",
				Path:        "/Applications/Calculator.app",
			},
			"safari": {
				Name:        "safari",
				DisplayName: "Safari浏览器",
				Path:        "/Applications/Safari.app",
			},
			"chrome": {
				Name:        "chrome",
				DisplayName: "Chrome浏览器",
				Path:        "/Applications/Google Chrome.app",
			},
			"finder": {
				Name:        "finder",
				DisplayName: "访达",
				Path:        "/System/Library/CoreServices/Finder.app",
			},
		}
	}
}

// Launch 启动应用程序
func (l *Launcher) Launch(appName string) *core.OperationResult {
	start := time.Now()

	// 查找预定义应用程序
	if app, exists := l.predefinedApps[strings.ToLower(appName)]; exists {
		return l.LaunchApp(app)
	}

	// 尝试直接启动
	result := l.LaunchWithPath(appName)
	result.SetDuration(start)
	return result
}

// LaunchWithPath 通过路径启动应用程序
func (l *Launcher) LaunchWithPath(path string, args ...string) *core.OperationResult {
	start := time.Now()

	// 展开环境变量
	expandedPath := os.ExpandEnv(path)

	// 检查文件是否存在
	if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
		result := core.NewErrorResult(
			fmt.Sprintf("应用程序不存在: %s", expandedPath),
			err,
		)
		result.SetDuration(start)
		return result
	}

	// 创建命令
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(expandedPath, args...)
	} else if runtime.GOOS == "darwin" {
		if strings.HasSuffix(expandedPath, ".app") {
			// macOS应用程序包
			cmd = exec.Command("open", expandedPath)
		} else {
			cmd = exec.Command(expandedPath, args...)
		}
	} else {
		// Linux
		cmd = exec.Command(expandedPath, args...)
	}

	// 启动应用程序
	err := cmd.Start()
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("启动应用程序失败: %s", expandedPath),
			err,
		)
		result.SetDuration(start)
		return result
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("成功启动应用程序: %s", expandedPath),
		map[string]interface{}{
			"path": expandedPath,
			"args": args,
			"pid":  cmd.Process.Pid,
		},
	)
	result.SetDuration(start)
	return result
}

// LaunchApp 启动预定义的应用程序
func (l *Launcher) LaunchApp(app *core.AppInfo) *core.OperationResult {
	start := time.Now()

	// 设置工作目录
	var cmd *exec.Cmd
	expandedPath := os.ExpandEnv(app.Path)

	if runtime.GOOS == "windows" {
		cmd = exec.Command(expandedPath, app.Args...)
	} else if runtime.GOOS == "darwin" {
		if strings.HasSuffix(expandedPath, ".app") {
			cmd = exec.Command("open", expandedPath)
		} else {
			cmd = exec.Command(expandedPath, app.Args...)
		}
	} else {
		cmd = exec.Command(expandedPath, app.Args...)
	}

	// 设置工作目录
	if app.WorkDir != "" {
		cmd.Dir = app.WorkDir
	}

	// 设置环境变量
	if app.Env != nil {
		env := os.Environ()
		for key, value := range app.Env {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = env
	}

	// 启动应用程序
	err := cmd.Start()
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("启动应用程序失败: %s (%s)", app.DisplayName, app.Path),
			err,
		)
		result.SetDuration(start)
		return result
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("成功启动应用程序: %s", app.DisplayName),
		map[string]interface{}{
			"name":         app.Name,
			"display_name": app.DisplayName,
			"path":         app.Path,
			"pid":          cmd.Process.Pid,
		},
	)
	result.SetDuration(start)
	return result
}

// GetInstalledApps 获取已安装的应用程序列表
func (l *Launcher) GetInstalledApps() ([]*core.AppInfo, *core.OperationResult) {
	start := time.Now()

	var apps []*core.AppInfo

	// 返回预定义的应用程序列表
	for _, app := range l.predefinedApps {
		// 检查应用程序是否存在
		expandedPath := os.ExpandEnv(app.Path)
		if _, err := os.Stat(expandedPath); err == nil {
			apps = append(apps, app)
		}
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("找到 %d 个已安装的应用程序", len(apps)),
		apps,
	)
	result.SetDuration(start)
	return apps, result
}

// FindApp 查找应用程序
func (l *Launcher) FindApp(name string) (*core.AppInfo, *core.OperationResult) {
	start := time.Now()

	lowerName := strings.ToLower(name)

	// 在预定义应用程序中查找
	for _, app := range l.predefinedApps {
		if strings.Contains(strings.ToLower(app.Name), lowerName) ||
			strings.Contains(strings.ToLower(app.DisplayName), lowerName) {

			// 检查应用程序是否存在
			expandedPath := os.ExpandEnv(app.Path)
			if _, err := os.Stat(expandedPath); err == nil {
				result := core.NewSuccessResult(
					fmt.Sprintf("找到应用程序: %s", app.DisplayName),
					app,
				)
				result.SetDuration(start)
				return app, result
			}
		}
	}

	result := core.NewErrorResult(
		fmt.Sprintf("未找到应用程序: %s", name),
		fmt.Errorf("application not found"),
	)
	result.SetDuration(start)
	return nil, result
}

// AddApp 添加自定义应用程序
func (l *Launcher) AddApp(app *core.AppInfo) *core.OperationResult {
	start := time.Now()

	l.predefinedApps[app.Name] = app

	result := core.NewSuccessResult(
		fmt.Sprintf("添加应用程序: %s", app.DisplayName),
		app,
	)
	result.SetDuration(start)
	return result
}

// LaunchURL 打开URL
func (l *Launcher) LaunchURL(url string) *core.OperationResult {
	start := time.Now()

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	err := cmd.Start()
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("打开URL失败: %s", url),
			err,
		)
		result.SetDuration(start)
		return result
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("成功打开URL: %s", url),
		map[string]interface{}{
			"url": url,
		},
	)
	result.SetDuration(start)
	return result
}
