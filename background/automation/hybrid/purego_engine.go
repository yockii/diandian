package hybrid

import (
	"fmt"
	"runtime"
	"time"

	"diandian/background/automation/core"

	"github.com/micmonay/keybd_event"
)

// PureGoEngine 纯Go实现的自动化引擎
type PureGoEngine struct {
	keybd *keybd_event.KeyBonding
}

// NewPureGoEngine 创建纯Go引擎
func NewPureGoEngine() (*PureGoEngine, error) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return nil, fmt.Errorf("failed to create keyboard binding: %v", err)
	}

	// Linux需要等待2秒
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	return &PureGoEngine{
		keybd: &kb,
	}, nil
}

// IsAvailable 检查引擎是否可用
func (p *PureGoEngine) IsAvailable() bool {
	return p.keybd != nil
}

// Click 点击操作 - 纯Go版本有限制
func (p *PureGoEngine) Click(x, y int) *core.OperationResult {
	start := time.Now()

	// 纯Go版本暂时不支持鼠标点击
	// 这里可以根据平台使用不同的实现
	switch runtime.GOOS {
	case "windows":
		// Windows下可以使用go-hook或系统调用
		return p.clickWindows(x, y)
	case "linux":
		// Linux下可以使用xdotool命令或X11调用
		return p.clickLinux(x, y)
	case "darwin":
		// macOS下可以使用CGEvent
		return p.clickMacOS(x, y)
	default:
		result := core.NewErrorResult(
			fmt.Sprintf("click not supported on %s", runtime.GOOS),
			fmt.Errorf("unsupported platform"),
		)
		result.SetDuration(start)
		return result
	}
}

// Type 输入文本
func (p *PureGoEngine) Type(text string) *core.OperationResult {
	start := time.Now()

	// 清除之前的按键设置
	p.keybd.Clear()

	// 逐字符输入
	for _, char := range text {
		// 将字符转换为按键码
		vk := p.charToVK(char)
		if vk == 0 {
			continue // 跳过不支持的字符
		}

		p.keybd.SetKeys(vk)

		// 处理大写字母
		if char >= 'A' && char <= 'Z' {
			p.keybd.HasSHIFT(true)
		} else {
			p.keybd.HasSHIFT(false)
		}

		err := p.keybd.Launching()
		if err != nil {
			result := core.NewErrorResult(
				fmt.Sprintf("failed to type character '%c'", char),
				err,
			)
			result.SetDuration(start)
			return result
		}

		// 短暂延迟，模拟真实输入
		time.Sleep(10 * time.Millisecond)
		p.keybd.Clear()
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("typed text: %s", text),
		map[string]interface{}{
			"text":   text,
			"length": len(text),
		},
	)
	result.SetDuration(start)
	return result
}

// KeyPress 按键操作
func (p *PureGoEngine) KeyPress(key string) *core.OperationResult {
	start := time.Now()

	vk := p.keyNameToVK(key)
	if vk == 0 {
		result := core.NewErrorResult(
			fmt.Sprintf("unsupported key: %s", key),
			fmt.Errorf("unknown key"),
		)
		result.SetDuration(start)
		return result
	}

	p.keybd.Clear()
	p.keybd.SetKeys(vk)

	err := p.keybd.Launching()
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("failed to press key: %s", key),
			err,
		)
		result.SetDuration(start)
		return result
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("pressed key: %s", key),
		map[string]interface{}{
			"key": key,
			"vk":  vk,
		},
	)
	result.SetDuration(start)
	return result
}

// Screenshot 截屏 - 使用系统调用
func (p *PureGoEngine) Screenshot() *core.OperationResult {
	start := time.Now()

	switch runtime.GOOS {
	case "windows":
		return p.screenshotWindows()
	case "linux":
		return p.screenshotLinux()
	case "darwin":
		return p.screenshotMacOS()
	default:
		result := core.NewErrorResult(
			fmt.Sprintf("screenshot not supported on %s", runtime.GOOS),
			fmt.Errorf("unsupported platform"),
		)
		result.SetDuration(start)
		return result
	}
}

// charToVK 将字符转换为虚拟键码
func (p *PureGoEngine) charToVK(char rune) int {
	switch {
	case char >= 'a' && char <= 'z':
		return keybd_event.VK_A + int(char-'a')
	case char >= 'A' && char <= 'Z':
		return keybd_event.VK_A + int(char-'A')
	case char >= '0' && char <= '9':
		return keybd_event.VK_0 + int(char-'0')
	case char == ' ':
		return keybd_event.VK_SPACE
	case char == '\n':
		return keybd_event.VK_ENTER
	case char == '\t':
		return keybd_event.VK_TAB
	default:
		return 0 // 不支持的字符
	}
}

// keyNameToVK 将按键名称转换为虚拟键码
func (p *PureGoEngine) keyNameToVK(keyName string) int {
	keyMap := map[string]int{
		"enter":     keybd_event.VK_ENTER,
		"space":     keybd_event.VK_SPACE,
		"tab":       keybd_event.VK_TAB,
		"escape":    keybd_event.VK_ESC,
		"backspace": keybd_event.VK_BACKSPACE,
		"delete":    keybd_event.VK_DELETE,
		"up":        keybd_event.VK_UP,
		"down":      keybd_event.VK_DOWN,
		"left":      keybd_event.VK_LEFT,
		"right":     keybd_event.VK_RIGHT,
		"f1":        keybd_event.VK_F1,
		"f2":        keybd_event.VK_F2,
		"f3":        keybd_event.VK_F3,
		"f4":        keybd_event.VK_F4,
		"f5":        keybd_event.VK_F5,
		"f6":        keybd_event.VK_F6,
		"f7":        keybd_event.VK_F7,
		"f8":        keybd_event.VK_F8,
		"f9":        keybd_event.VK_F9,
		"f10":       keybd_event.VK_F10,
		"f11":       keybd_event.VK_F11,
		"f12":       keybd_event.VK_F12,
	}

	if vk, exists := keyMap[keyName]; exists {
		return vk
	}

	// 尝试单字符按键
	if len(keyName) == 1 {
		return p.charToVK(rune(keyName[0]))
	}

	return 0
}
