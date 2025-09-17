package service

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// AppLauncher 智能应用启动器
type AppLauncher struct{}

// NewAppLauncher 创建应用启动器
func NewAppLauncher() *AppLauncher {
	return &AppLauncher{}
}

// LaunchApp 智能启动应用
func (al *AppLauncher) LaunchApp(appIdentifier string) error {
	slog.Info("尝试启动应用", "app", appIdentifier)

	// 检查是否是第三方应用标记
	if strings.HasPrefix(appIdentifier, "THIRD_PARTY:") {
		appName := strings.TrimPrefix(appIdentifier, "THIRD_PARTY:")
		return al.launchThirdPartyApp(appName)
	}

	// 尝试直接启动系统应用
	return al.launchSystemApp(appIdentifier)
}

// launchSystemApp 启动系统内置应用
func (al *AppLauncher) launchSystemApp(appName string) error {
	slog.Info("启动系统应用", "app", appName)

	cmd := exec.Command(appName)
	err := cmd.Start()
	if err != nil {
		slog.Error("系统应用启动失败", "app", appName, "error", err)
		return fmt.Errorf("启动系统应用失败: %v", err)
	}

	slog.Info("系统应用启动成功", "app", appName, "pid", cmd.Process.Pid)
	return nil
}

// launchThirdPartyApp 启动第三方应用
func (al *AppLauncher) launchThirdPartyApp(appName string) error {
	slog.Info("尝试启动第三方应用", "app", appName)

	// 策略1: 尝试注册表查找
	if path, err := al.findAppInRegistry(appName); err == nil && path != "" {
		slog.Info("通过注册表找到应用", "app", appName, "path", path)
		return al.launchByPath(path)
	}

	// 策略2: 尝试常见安装路径
	if path, err := al.findAppInCommonPaths(appName); err == nil && path != "" {
		slog.Info("通过常见路径找到应用", "app", appName, "path", path)
		return al.launchByPath(path)
	}

	// 策略3: 尝试桌面快捷方式
	if path, err := al.findAppInDesktop(appName); err == nil && path != "" {
		slog.Info("通过桌面快捷方式找到应用", "app", appName, "path", path)
		return al.launchByPath(path)
	}

	// 策略4: 尝试开始菜单
	if path, err := al.findAppInStartMenu(appName); err == nil && path != "" {
		slog.Info("通过开始菜单找到应用", "app", appName, "path", path)
		return al.launchByPath(path)
	}

	// 策略5: 使用Windows搜索API（如果可用）
	if path, err := al.findAppByWindowsSearch(appName); err == nil && path != "" {
		slog.Info("通过Windows搜索找到应用", "app", appName, "path", path)
		return al.launchByPath(path)
	}

	// 所有策略都失败，返回错误
	slog.Error("无法找到第三方应用", "app", appName)
	return fmt.Errorf("无法找到应用: %s，建议检查应用是否已安装", appName)
}

// findAppInRegistry 在注册表中查找应用
func (al *AppLauncher) findAppInRegistry(appName string) (string, error) {
	// 查找已安装程序列表
	registryPaths := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	for _, regPath := range registryPaths {
		if path, err := al.searchInRegistryPath(regPath, appName); err == nil && path != "" {
			return path, nil
		}
	}

	return "", fmt.Errorf("未在注册表中找到应用")
}

// searchInRegistryPath 在指定注册表路径中搜索
func (al *AppLauncher) searchInRegistryPath(regPath, appName string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return "", err
	}
	defer key.Close()

	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return "", err
	}

	for _, subkey := range subkeys {
		subkeyPath := regPath + `\` + subkey
		if path, err := al.checkRegistryEntry(subkeyPath, appName); err == nil && path != "" {
			return path, nil
		}
	}

	return "", fmt.Errorf("未找到")
}

// checkRegistryEntry 检查注册表条目
func (al *AppLauncher) checkRegistryEntry(keyPath, appName string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()

	// 检查显示名称
	if displayName, _, err := key.GetStringValue("DisplayName"); err == nil {
		if strings.Contains(strings.ToLower(displayName), strings.ToLower(appName)) {
			// 尝试获取安装路径
			if installLocation, _, err := key.GetStringValue("InstallLocation"); err == nil {
				return al.findExecutableInPath(installLocation, appName)
			}
			// 尝试获取卸载字符串中的路径
			if uninstallString, _, err := key.GetStringValue("UninstallString"); err == nil {
				if exePath := al.extractPathFromUninstallString(uninstallString); exePath != "" {
					return exePath, nil
				}
			}
		}
	}

	return "", fmt.Errorf("未找到匹配项")
}

// findAppInCommonPaths 在常见安装路径中查找应用
func (al *AppLauncher) findAppInCommonPaths(appName string) (string, error) {
	commonPaths := []string{
		`C:\Program Files`,
		`C:\Program Files (x86)`,
		`C:\Users\` + os.Getenv("USERNAME") + `\AppData\Local`,
		`C:\Users\` + os.Getenv("USERNAME") + `\AppData\Roaming`,
	}

	for _, basePath := range commonPaths {
		if path, err := al.searchInDirectory(basePath, appName, 2); err == nil && path != "" {
			return path, nil
		}
	}

	return "", fmt.Errorf("未在常见路径中找到应用")
}

// findAppInDesktop 在桌面快捷方式中查找应用
func (al *AppLauncher) findAppInDesktop(appName string) (string, error) {
	desktopPaths := []string{
		filepath.Join(os.Getenv("USERPROFILE"), "Desktop"),
		filepath.Join(os.Getenv("PUBLIC"), "Desktop"),
	}

	for _, desktopPath := range desktopPaths {
		if path, err := al.searchShortcutsInDirectory(desktopPath, appName); err == nil && path != "" {
			return path, nil
		}
	}

	return "", fmt.Errorf("未在桌面找到应用快捷方式")
}

// findAppInStartMenu 在开始菜单中查找应用
func (al *AppLauncher) findAppInStartMenu(appName string) (string, error) {
	startMenuPaths := []string{
		filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs"),
		filepath.Join(os.Getenv("ALLUSERSPROFILE"), "Microsoft", "Windows", "Start Menu", "Programs"),
	}

	for _, startMenuPath := range startMenuPaths {
		if path, err := al.searchShortcutsInDirectory(startMenuPath, appName); err == nil && path != "" {
			return path, nil
		}
	}

	return "", fmt.Errorf("未在开始菜单找到应用")
}

// findAppByWindowsSearch 使用Windows搜索API查找应用
func (al *AppLauncher) findAppByWindowsSearch(appName string) (string, error) {
	// 这里可以实现Windows搜索API调用
	// 由于复杂性，暂时返回未实现
	return "", fmt.Errorf("Windows搜索API暂未实现")
}

// searchInDirectory 在目录中搜索应用
func (al *AppLauncher) searchInDirectory(dirPath, appName string, maxDepth int) (string, error) {
	if maxDepth <= 0 {
		return "", fmt.Errorf("达到最大搜索深度")
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			// 检查目录名是否匹配
			if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(appName)) {
				if exePath, err := al.findExecutableInPath(fullPath, appName); err == nil && exePath != "" {
					return exePath, nil
				}
			}
			// 递归搜索子目录
			if path, err := al.searchInDirectory(fullPath, appName, maxDepth-1); err == nil && path != "" {
				return path, nil
			}
		} else {
			// 检查是否是可执行文件
			if strings.HasSuffix(strings.ToLower(entry.Name()), ".exe") {
				if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(appName)) {
					return fullPath, nil
				}
			}
		}
	}

	return "", fmt.Errorf("未找到")
}

// searchShortcutsInDirectory 在目录中搜索快捷方式
func (al *AppLauncher) searchShortcutsInDirectory(dirPath, appName string) (string, error) {
	return al.searchInDirectoryRecursive(dirPath, appName, ".lnk", 3)
}

// searchInDirectoryRecursive 递归搜索目录
func (al *AppLauncher) searchInDirectoryRecursive(dirPath, appName, extension string, maxDepth int) (string, error) {
	if maxDepth <= 0 {
		return "", fmt.Errorf("达到最大搜索深度")
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			if path, err := al.searchInDirectoryRecursive(fullPath, appName, extension, maxDepth-1); err == nil && path != "" {
				return path, nil
			}
		} else {
			if strings.HasSuffix(strings.ToLower(entry.Name()), extension) {
				if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(appName)) {
					if extension == ".lnk" {
						// 解析快捷方式获取目标路径
						if targetPath, err := al.resolveShortcut(fullPath); err == nil && targetPath != "" {
							return targetPath, nil
						}
					}
					return fullPath, nil
				}
			}
		}
	}

	return "", fmt.Errorf("未找到")
}

// findExecutableInPath 在指定路径中查找可执行文件
func (al *AppLauncher) findExecutableInPath(dirPath, appName string) (string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".exe") {
			if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(appName)) {
				return filepath.Join(dirPath, entry.Name()), nil
			}
		}
	}

	return "", fmt.Errorf("未找到可执行文件")
}

// extractPathFromUninstallString 从卸载字符串中提取路径
func (al *AppLauncher) extractPathFromUninstallString(uninstallString string) string {
	// 简单的路径提取逻辑
	if strings.Contains(uninstallString, ".exe") {
		parts := strings.Split(uninstallString, ".exe")
		if len(parts) > 0 {
			exePath := parts[0] + ".exe"
			exePath = strings.Trim(exePath, `"`)
			if _, err := os.Stat(exePath); err == nil {
				return exePath
			}
		}
	}
	return ""
}

// resolveShortcut 解析快捷方式获取目标路径
func (al *AppLauncher) resolveShortcut(shortcutPath string) (string, error) {
	// 这里需要实现Windows快捷方式解析
	// 由于复杂性，暂时返回快捷方式路径本身
	return shortcutPath, nil
}

// launchByPath 通过路径启动应用
func (al *AppLauncher) launchByPath(path string) error {
	slog.Info("通过路径启动应用", "path", path)

	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", path)
	}

	// 如果是快捷方式，使用系统默认方式打开
	if strings.HasSuffix(strings.ToLower(path), ".lnk") {
		cmd := exec.Command("cmd", "/c", "start", "", path)
		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("启动快捷方式失败: %v", err)
		}
		return nil
	}

	// 直接启动可执行文件
	cmd := exec.Command(path)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("启动应用失败: %v", err)
	}

	slog.Info("应用启动成功", "path", path, "pid", cmd.Process.Pid)
	return nil
}
