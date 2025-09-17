package hybrid

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"diandian/background/automation/core"
)

// 平台特定的实现

// Windows平台实现
func (p *PureGoEngine) clickWindows(x, y int) *core.OperationResult {
	start := time.Now()
	
	// 使用PowerShell调用Windows API
	cmd := exec.Command("powershell", "-Command", 
		fmt.Sprintf(`
Add-Type -AssemblyName System.Windows.Forms
[System.Windows.Forms.Cursor]::Position = New-Object System.Drawing.Point(%d, %d)
Add-Type -TypeDefinition '
using System;
using System.Runtime.InteropServices;
public class Mouse {
    [DllImport("user32.dll")]
    public static extern void mouse_event(uint dwFlags, uint dx, uint dy, uint dwData, IntPtr dwExtraInfo);
    public const uint MOUSEEVENTF_LEFTDOWN = 0x02;
    public const uint MOUSEEVENTF_LEFTUP = 0x04;
}
'
[Mouse]::mouse_event(0x02, 0, 0, 0, [IntPtr]::Zero)
[Mouse]::mouse_event(0x04, 0, 0, 0, [IntPtr]::Zero)
`, x, y))

	err := cmd.Run()
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("failed to click at (%d, %d)", x, y),
			err,
		)
		result.SetDuration(start)
		return result
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("clicked at (%d, %d)", x, y),
		map[string]interface{}{
			"x": x,
			"y": y,
		},
	)
	result.SetDuration(start)
	return result
}

func (p *PureGoEngine) screenshotWindows() *core.OperationResult {
	start := time.Now()
	
	// 使用PowerShell截屏
	cmd := exec.Command("powershell", "-Command", `
Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing
$bounds = [System.Windows.Forms.Screen]::PrimaryScreen.Bounds
$bitmap = New-Object System.Drawing.Bitmap $bounds.Width, $bounds.Height
$graphics = [System.Drawing.Graphics]::FromImage($bitmap)
$graphics.CopyFromScreen($bounds.Location, [System.Drawing.Point]::Empty, $bounds.Size)
$ms = New-Object System.IO.MemoryStream
$bitmap.Save($ms, [System.Drawing.Imaging.ImageFormat]::Png)
[System.Convert]::ToBase64String($ms.ToArray())
`)

	output, err := cmd.Output()
	if err != nil {
		result := core.NewErrorResult("failed to take screenshot", err)
		result.SetDuration(start)
		return result
	}

	result := core.NewSuccessResult(
		"screenshot taken",
		map[string]interface{}{
			"format": "png",
			"data":   string(output),
		},
	)
	result.SetDuration(start)
	return result
}

// Linux平台实现
func (p *PureGoEngine) clickLinux(x, y int) *core.OperationResult {
	start := time.Now()
	
	// 使用xdotool
	cmd := exec.Command("xdotool", "mousemove", fmt.Sprintf("%d", x), fmt.Sprintf("%d", y), "click", "1")
	err := cmd.Run()
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("failed to click at (%d, %d)", x, y),
			err,
		)
		result.SetDuration(start)
		return result
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("clicked at (%d, %d)", x, y),
		map[string]interface{}{
			"x": x,
			"y": y,
		},
	)
	result.SetDuration(start)
	return result
}

func (p *PureGoEngine) screenshotLinux() *core.OperationResult {
	start := time.Now()
	
	// 使用scrot或gnome-screenshot
	var cmd *exec.Cmd
	
	// 尝试scrot
	if _, err := exec.LookPath("scrot"); err == nil {
		cmd = exec.Command("scrot", "-z", "-")
	} else if _, err := exec.LookPath("gnome-screenshot"); err == nil {
		cmd = exec.Command("gnome-screenshot", "-f", "/dev/stdout")
	} else {
		result := core.NewErrorResult(
			"no screenshot tool available (scrot or gnome-screenshot required)",
			fmt.Errorf("missing screenshot tool"),
		)
		result.SetDuration(start)
		return result
	}

	output, err := cmd.Output()
	if err != nil {
		result := core.NewErrorResult("failed to take screenshot", err)
		result.SetDuration(start)
		return result
	}

	result := core.NewSuccessResult(
		"screenshot taken",
		map[string]interface{}{
			"format": "png",
			"data":   output,
		},
	)
	result.SetDuration(start)
	return result
}

// macOS平台实现
func (p *PureGoEngine) clickMacOS(x, y int) *core.OperationResult {
	start := time.Now()
	
	// 使用osascript调用AppleScript
	script := fmt.Sprintf(`
tell application "System Events"
	click at {%d, %d}
end tell
`, x, y)

	cmd := exec.Command("osascript", "-e", script)
	err := cmd.Run()
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("failed to click at (%d, %d)", x, y),
			err,
		)
		result.SetDuration(start)
		return result
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("clicked at (%d, %d)", x, y),
		map[string]interface{}{
			"x": x,
			"y": y,
		},
	)
	result.SetDuration(start)
	return result
}

func (p *PureGoEngine) screenshotMacOS() *core.OperationResult {
	start := time.Now()
	
	// 使用screencapture
	cmd := exec.Command("screencapture", "-t", "png", "-")
	output, err := cmd.Output()
	if err != nil {
		result := core.NewErrorResult("failed to take screenshot", err)
		result.SetDuration(start)
		return result
	}

	result := core.NewSuccessResult(
		"screenshot taken",
		map[string]interface{}{
			"format": "png",
			"data":   output,
		},
	)
	result.SetDuration(start)
	return result
}

// 检查平台特定工具是否可用
func (p *PureGoEngine) checkPlatformTools() bool {
	switch runtime.GOOS {
	case "windows":
		// Windows PowerShell通常都有
		return true
	case "linux":
		// 检查xdotool是否可用
		_, err := exec.LookPath("xdotool")
		return err == nil
	case "darwin":
		// macOS osascript通常都有
		return true
	default:
		return false
	}
}
