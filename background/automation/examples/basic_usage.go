package examples

import (
	"fmt"
	"log"
	"strings"

	"diandian/background/automation/engine"
)

// BasicUsageExample 基本使用示例
func BasicUsageExample() {
	// 创建自动化引擎
	autoEngine := engine.NewEngine()

	// 初始化引擎
	result := autoEngine.Initialize()
	if !result.Success {
		log.Printf("初始化失败: %s", result.Error)
		return
	}

	fmt.Println("=== 自动化操作示例 ===")

	// 1. 屏幕操作示例
	fmt.Println("\n1. 屏幕操作:")

	// 获取屏幕尺寸
	size, result := autoEngine.GetScreenSize()
	if result.Success {
		fmt.Printf("屏幕尺寸: %dx%d\n", size.Width, size.Height)
	}

	// 获取鼠标位置
	pos, result := autoEngine.GetPosition()
	if result.Success {
		fmt.Printf("当前鼠标位置: (%d, %d)\n", pos.X, pos.Y)
	}

	// 截屏
	_, result = autoEngine.Screenshot()
	if result.Success {
		fmt.Println("截屏成功")
	}

	// 2. 应用程序操作示例
	fmt.Println("\n2. 应用程序操作:")

	// 获取已安装的应用程序
	apps, result := autoEngine.GetInstalledApps()
	if result.Success {
		fmt.Printf("找到 %d 个已安装的应用程序:\n", len(apps))
		for _, app := range apps {
			fmt.Printf("  - %s (%s)\n", app.DisplayName, app.Name)
		}
	}

	// 3. 文件操作示例
	fmt.Println("\n3. 文件操作:")

	// 创建测试文件
	testFile := "test_automation.txt"
	testContent := "这是自动化测试文件\n创建时间: " + result.Timestamp.String()

	result = autoEngine.WriteTextFile(testFile, testContent)
	if result.Success {
		fmt.Printf("创建文件成功: %s\n", testFile)

		// 读取文件
		content, result := autoEngine.ReadTextFile(testFile)
		if result.Success {
			fmt.Printf("文件内容: %s\n", content)
		}

		// 删除测试文件
		result = autoEngine.DeleteFile(testFile)
		if result.Success {
			fmt.Println("删除测试文件成功")
		}
	}

	// 4. 系统操作示例
	fmt.Println("\n4. 系统操作:")

	// 获取剪贴板内容
	clipText, result := autoEngine.GetClipboard()
	if result.Success {
		fmt.Printf("剪贴板内容: %s\n", clipText)
	}

	// 设置剪贴板内容
	testClipText := "自动化测试剪贴板内容"
	result = autoEngine.SetClipboard(testClipText)
	if result.Success {
		fmt.Println("设置剪贴板内容成功")

		// 验证设置
		newClipText, result := autoEngine.GetClipboard()
		if result.Success && newClipText == testClipText {
			fmt.Println("剪贴板内容验证成功")
		}
	}

	// 获取活动窗口
	window, result := autoEngine.GetActiveWindow()
	if result.Success {
		fmt.Printf("活动窗口: %s (PID: %d)\n", window.Title, window.PID)
	}

	// 清理资源
	result = autoEngine.Cleanup()
	if result.Success {
		fmt.Println("\n引擎清理完成")
	}
}

// MouseOperationExample 鼠标操作示例
func MouseOperationExample() {
	autoEngine := engine.NewEngine()
	autoEngine.Initialize()

	fmt.Println("=== 鼠标操作示例 ===")

	// 获取当前鼠标位置
	pos, result := autoEngine.GetPosition()
	if result.Success {
		fmt.Printf("当前鼠标位置: (%d, %d)\n", pos.X, pos.Y)

		// 移动鼠标
		newX, newY := pos.X+100, pos.Y+100
		result = autoEngine.Move(newX, newY)
		if result.Success {
			fmt.Printf("鼠标移动到: (%d, %d)\n", newX, newY)

			// 等待一秒
			autoEngine.Wait(1000)

			// 移动回原位置
			autoEngine.Move(pos.X, pos.Y)
			fmt.Println("鼠标移动回原位置")
		}
	}

	autoEngine.Cleanup()
}

// KeyboardOperationExample 键盘操作示例
func KeyboardOperationExample() {
	autoEngine := engine.NewEngine()
	autoEngine.Initialize()

	fmt.Println("=== 键盘操作示例 ===")

	// 注意：这些操作会影响当前活动窗口，请谨慎使用
	fmt.Println("将在3秒后开始键盘操作演示...")
	autoEngine.Wait(3000)

	// 输入文本
	result := autoEngine.Type("Hello, Automation!")
	if result.Success {
		fmt.Println("文本输入成功")
	}

	// 全选
	autoEngine.Wait(500)
	result = autoEngine.SelectAll()
	if result.Success {
		fmt.Println("全选操作成功")
	}

	// 复制
	autoEngine.Wait(500)
	result = autoEngine.Copy()
	if result.Success {
		fmt.Println("复制操作成功")
	}

	// 按Delete键删除
	autoEngine.Wait(500)
	result = autoEngine.KeyPress("delete")
	if result.Success {
		fmt.Println("删除操作成功")
	}

	// 粘贴
	autoEngine.Wait(500)
	result = autoEngine.Paste()
	if result.Success {
		fmt.Println("粘贴操作成功")
	}

	autoEngine.Cleanup()
}

// AppLaunchExample 应用程序启动示例
func AppLaunchExample() {
	autoEngine := engine.NewEngine()
	autoEngine.Initialize()

	fmt.Println("=== 应用程序启动示例 ===")

	// 启动记事本
	result := autoEngine.Launch("notepad")
	if result.Success {
		fmt.Println("记事本启动成功")

		// 等待应用程序启动
		autoEngine.Wait(2000)

		// 在记事本中输入文本
		autoEngine.Type("这是通过自动化程序输入的文本！\n")
		autoEngine.Type("当前时间: " + result.Timestamp.String())

		fmt.Println("文本输入完成")
	} else {
		fmt.Printf("记事本启动失败: %s\n", result.Error)
	}

	// 启动计算器
	result = autoEngine.Launch("calculator")
	if result.Success {
		fmt.Println("计算器启动成功")
	} else {
		fmt.Printf("计算器启动失败: %s\n", result.Error)
	}

	autoEngine.Cleanup()
}

// FileOperationExample 文件操作示例
func FileOperationExample() {
	autoEngine := engine.NewEngine()
	autoEngine.Initialize()

	fmt.Println("=== 文件操作示例 ===")

	// 创建测试目录
	testDir := "automation_test"
	result := autoEngine.CreateDir(testDir)
	if result.Success {
		fmt.Printf("创建目录成功: %s\n", testDir)

		// 创建多个测试文件
		for i := 1; i <= 3; i++ {
			fileName := fmt.Sprintf("%s/test_file_%d.txt", testDir, i)
			content := fmt.Sprintf("这是测试文件 %d\n内容由自动化程序创建", i)

			result = autoEngine.WriteTextFile(fileName, content)
			if result.Success {
				fmt.Printf("创建文件: %s\n", fileName)
			}
		}

		// 列出目录内容
		files, result := autoEngine.ListDir(testDir)
		if result.Success {
			fmt.Printf("目录 %s 包含 %d 个文件:\n", testDir, len(files))
			for _, file := range files {
				fmt.Printf("  - %s\n", file)
			}
		}

		// 复制文件
		srcFile := testDir + "/test_file_1.txt"
		dstFile := testDir + "/test_file_1_copy.txt"
		result = autoEngine.CopyFile(srcFile, dstFile)
		if result.Success {
			fmt.Printf("文件复制成功: %s -> %s\n", srcFile, dstFile)
		}

		// 重命名文件
		oldName := testDir + "/test_file_2.txt"
		newName := testDir + "/renamed_file.txt"
		result = autoEngine.RenameFile(oldName, newName)
		if result.Success {
			fmt.Printf("文件重命名成功: %s -> %s\n", oldName, newName)
		}

		// 清理测试文件
		fmt.Println("清理测试文件...")
		autoEngine.DeleteDir(testDir)
		fmt.Println("清理完成")
	}

	autoEngine.Cleanup()
}

// RunAllExamples 运行所有示例
func RunAllExamples() {
	fmt.Println("开始运行自动化操作示例...")

	BasicUsageExample()
	fmt.Println("\n" + strings.Repeat("=", 50))

	MouseOperationExample()
	fmt.Println("\n" + strings.Repeat("=", 50))

	// 注意：键盘操作示例可能会影响当前活动窗口
	// KeyboardOperationExample()
	// fmt.Println("\n" + strings.Repeat("=", 50))

	AppLaunchExample()
	fmt.Println("\n" + strings.Repeat("=", 50))

	FileOperationExample()

	fmt.Println("\n所有示例运行完成！")
}
