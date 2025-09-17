package executor

import (
	"regexp"
	"strconv"
	"strings"
)

// 辅助方法：从上下文中提取信息

// ExtractAppNameFromContext 公开的应用名称提取方法（用于测试）
func (e *EnhancedTaskExecutionEngine) ExtractAppNameFromContext(context string) string {
	return e.extractAppNameFromContext(context)
}

func (e *EnhancedTaskExecutionEngine) extractAppNameFromContext(context string) string {
	context = strings.ToLower(context)

	// 系统内置应用映射表（可直接通过命令启动）
	systemAppMappings := map[string]string{
		"记事本":        "notepad",
		"notepad":    "notepad",
		"计算器":        "calc",
		"calculator": "calc",
		"画图":         "mspaint",
		"paint":      "mspaint",
		"cmd":        "cmd",
		"命令提示符":      "cmd",
		"powershell": "powershell",
		"任务管理器":      "taskmgr",
		"控制面板":       "control",
		"注册表编辑器":     "regedit",
		"系统配置":       "msconfig",
	}

	// 检查系统内置应用
	for keyword, appName := range systemAppMappings {
		if strings.Contains(context, keyword) {
			return appName
		}
	}

	// 常见第三方应用的可能名称（需要通过其他方式启动）
	thirdPartyApps := []string{
		"wps", "wps文字", "wps writer",
		"word", "winword", "microsoft word",
		"excel", "microsoft excel",
		"powerpoint", "ppt", "microsoft powerpoint",
		"chrome", "谷歌浏览器", "google chrome",
		"firefox", "火狐浏览器",
		"edge", "microsoft edge",
		"微信", "wechat",
		"qq", "腾讯qq",
		"钉钉", "dingtalk",
		"vscode", "visual studio code",
		"photoshop", "ps",
		"autocad", "cad",
	}

	// 检查是否是第三方应用
	for _, app := range thirdPartyApps {
		if strings.Contains(context, app) {
			// 返回特殊标记，表示需要通过其他方式启动
			return "THIRD_PARTY:" + app
		}
	}

	// 尝试提取可能的应用名称
	patterns := []string{
		`(?:启动|打开|运行)\s*([a-zA-Z\p{Han}]+)`,
		`([a-zA-Z\p{Han}]+)(?:应用|程序|软件)`,
		`([a-zA-Z]+)\.exe`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(context)
		if len(matches) > 1 {
			appName := strings.ToLower(matches[1])
			// 如果提取的名称不在系统应用中，标记为第三方
			if _, exists := systemAppMappings[appName]; !exists {
				return "THIRD_PARTY:" + appName
			}
			return appName
		}
	}

	// 默认返回记事本
	return "notepad"
}

func (e *EnhancedTaskExecutionEngine) extractPathFromContext(context string) string {
	// 尝试从上下文中提取文件路径

	// 匹配常见的文件路径模式
	pathPatterns := []string{
		`([a-zA-Z]:\\[^<>:"|?*\n\r]+)`,        // Windows绝对路径
		`([a-zA-Z]:/[^<>:"|?*\n\r]+)`,         // Windows路径（正斜杠）
		`([./][^<>:"|?*\n\r]+\.[a-zA-Z0-9]+)`, // 相对路径
		`([a-zA-Z0-9_.-]+\.[a-zA-Z0-9]+)`,     // 简单文件名
	}

	for _, pattern := range pathPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(context)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	// 如果没有找到路径，根据上下文生成默认路径
	if strings.Contains(strings.ToLower(context), "截图") || strings.Contains(strings.ToLower(context), "screenshot") {
		return "screenshot.png"
	}

	// 默认路径
	return "output.txt"
}

func (e *EnhancedTaskExecutionEngine) isGetClipboardOperation(context string) bool {
	// 判断是否是获取剪贴板操作
	return strings.Contains(strings.ToLower(context), "获取") || strings.Contains(strings.ToLower(context), "get")
}

func (e *EnhancedTaskExecutionEngine) extractTextFromContext(context string) string {
	// 从上下文中提取文本内容

	// 匹配引号中的文本
	quotedPatterns := []string{
		`"([^"]+)"`, // 双引号
		`'([^']+)'`, // 单引号
		`"([^"]+)"`, // 中文双引号
		`'([^']+)'`, // 中文单引号
	}

	for _, pattern := range quotedPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(context)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	// 匹配 "输入XXX" 或 "写入XXX" 模式
	inputPatterns := []string{
		`(?:输入|写入|键入)[:：]\s*(.+)`,
		`(?:输入|写入|键入)\s*(.+)`,
		`(?:文字|文本|内容)[:：]\s*(.+)`,
	}

	for _, pattern := range inputPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(context)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	// 如果没有找到特定模式，返回整个上下文（去除常见的指令词）
	cleanText := context
	removeWords := []string{"输入", "写入", "键入", "文字", "文本", "内容", "：", ":"}
	for _, word := range removeWords {
		cleanText = strings.ReplaceAll(cleanText, word, "")
	}
	cleanText = strings.TrimSpace(cleanText)

	if cleanText != "" {
		return cleanText
	}

	// 默认文本
	return "测试文本"
}

func (e *EnhancedTaskExecutionEngine) extractDurationFromContext(context string) int {
	// 从上下文中提取等待时间，返回毫秒

	// 匹配数字+时间单位的模式
	timePatterns := []struct {
		pattern    string
		multiplier int
	}{
		{`(\d+(?:\.\d+)?)\s*毫秒`, 1},
		{`(\d+(?:\.\d+)?)\s*ms`, 1},
		{`(\d+(?:\.\d+)?)\s*秒`, 1000},
		{`(\d+(?:\.\d+)?)\s*s`, 1000},
		{`(\d+(?:\.\d+)?)\s*分钟`, 60000},
		{`(\d+(?:\.\d+)?)\s*min`, 60000},
		{`(\d+(?:\.\d+)?)\s*分`, 60000},
	}

	for _, tp := range timePatterns {
		re := regexp.MustCompile(tp.pattern)
		matches := re.FindStringSubmatch(context)
		if len(matches) > 1 {
			if duration, err := strconv.ParseFloat(matches[1], 64); err == nil {
				return int(duration * float64(tp.multiplier))
			}
		}
	}

	// 匹配纯数字（默认为秒）
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(context)
	if len(matches) > 1 {
		if duration, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return int(duration * 1000) // 默认为秒，转换为毫秒
		}
	}

	return 1000 // 默认1秒
}

func (e *EnhancedTaskExecutionEngine) extractKeyFromContext(context string) (string, []string) {
	// 从上下文中提取按键和修饰键
	context = strings.ToLower(context)

	// 组合键映射
	comboMappings := map[string]struct {
		key       string
		modifiers []string
	}{
		"ctrl+c":       {"c", []string{"ctrl"}},
		"ctrl+v":       {"v", []string{"ctrl"}},
		"ctrl+x":       {"x", []string{"ctrl"}},
		"ctrl+z":       {"z", []string{"ctrl"}},
		"ctrl+y":       {"y", []string{"ctrl"}},
		"ctrl+a":       {"a", []string{"ctrl"}},
		"ctrl+s":       {"s", []string{"ctrl"}},
		"ctrl+o":       {"o", []string{"ctrl"}},
		"ctrl+n":       {"n", []string{"ctrl"}},
		"ctrl+f":       {"f", []string{"ctrl"}},
		"ctrl+h":       {"h", []string{"ctrl"}},
		"ctrl+p":       {"p", []string{"ctrl"}},
		"alt+f4":       {"f4", []string{"alt"}},
		"alt+tab":      {"tab", []string{"alt"}},
		"shift+tab":    {"tab", []string{"shift"}},
		"ctrl+shift+n": {"n", []string{"ctrl", "shift"}},
	}

	// 检查组合键
	for combo, keyInfo := range comboMappings {
		if strings.Contains(context, combo) {
			return keyInfo.key, keyInfo.modifiers
		}
	}

	// 单键映射
	singleKeyMappings := map[string]string{
		"enter":     "enter",
		"回车":        "enter",
		"回车键":       "enter",
		"space":     "space",
		"空格":        "space",
		"空格键":       "space",
		"escape":    "escape",
		"esc":       "escape",
		"退出":        "escape",
		"tab":       "tab",
		"制表符":       "tab",
		"delete":    "delete",
		"删除":        "delete",
		"backspace": "backspace",
		"退格":        "backspace",
		"home":      "home",
		"end":       "end",
		"pageup":    "pageup",
		"pagedown":  "pagedown",
		"up":        "up",
		"down":      "down",
		"left":      "left",
		"right":     "right",
		"上":         "up",
		"下":         "down",
		"左":         "left",
		"右":         "right",
		"f1":        "f1",
		"f2":        "f2",
		"f3":        "f3",
		"f4":        "f4",
		"f5":        "f5",
		"f6":        "f6",
		"f7":        "f7",
		"f8":        "f8",
		"f9":        "f9",
		"f10":       "f10",
		"f11":       "f11",
		"f12":       "f12",
	}

	// 检查单键
	for keyword, key := range singleKeyMappings {
		if strings.Contains(context, keyword) {
			return key, []string{}
		}
	}

	// 尝试匹配通用模式 "按XXX键"
	re := regexp.MustCompile(`按\s*([a-zA-Z0-9]+)\s*键?`)
	matches := re.FindStringSubmatch(context)
	if len(matches) > 1 {
		return strings.ToLower(matches[1]), []string{}
	}

	// 默认回车键
	return "enter", []string{}
}
