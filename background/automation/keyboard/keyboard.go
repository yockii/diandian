package keyboard

import (
	"fmt"
	"strings"
	"time"

	"diandian/background/automation/core"

	"github.com/go-vgo/robotgo"
)

// Keyboard 键盘操作实现
type Keyboard struct{}

// NewKeyboard 创建键盘操作实例
func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

// Type 输入文本
func (k *Keyboard) Type(text string) *core.OperationResult {
	start := time.Now()

	robotgo.TypeStr(text)

	result := core.NewSuccessResult(
		fmt.Sprintf("输入文本: %s", text),
		map[string]interface{}{
			"text":   text,
			"length": len(text),
		},
	)
	result.SetDuration(start)
	return result
}

// TypeSlow 慢速输入文本（模拟人工输入）
func (k *Keyboard) TypeSlow(text string, delay time.Duration) *core.OperationResult {
	start := time.Now()

	for _, char := range text {
		robotgo.TypeStr(string(char))
		time.Sleep(delay)
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("慢速输入文本: %s", text),
		map[string]interface{}{
			"text":  text,
			"delay": delay.String(),
		},
	)
	result.SetDuration(start)
	return result
}

// KeyPress 按下并释放按键
func (k *Keyboard) KeyPress(key string) *core.OperationResult {
	start := time.Now()

	robotgo.KeyTap(key)

	result := core.NewSuccessResult(
		fmt.Sprintf("按键: %s", key),
		map[string]interface{}{
			"key": key,
		},
	)
	result.SetDuration(start)
	return result
}

// KeyDown 按下按键
func (k *Keyboard) KeyDown(key string) *core.OperationResult {
	start := time.Now()

	robotgo.KeyToggle(key, "down")

	result := core.NewSuccessResult(
		fmt.Sprintf("按下按键: %s", key),
		map[string]interface{}{
			"key":    key,
			"action": "down",
		},
	)
	result.SetDuration(start)
	return result
}

// KeyUp 释放按键
func (k *Keyboard) KeyUp(key string) *core.OperationResult {
	start := time.Now()

	robotgo.KeyToggle(key, "up")

	result := core.NewSuccessResult(
		fmt.Sprintf("释放按键: %s", key),
		map[string]interface{}{
			"key":    key,
			"action": "up",
		},
	)
	result.SetDuration(start)
	return result
}

// Hotkey 组合键操作
func (k *Keyboard) Hotkey(modifiers []core.KeyModifier, key string) *core.OperationResult {
	start := time.Now()

	// 转换修饰键
	var modifierStrs []string
	for _, mod := range modifiers {
		switch mod {
		case core.ModCtrl:
			modifierStrs = append(modifierStrs, "ctrl")
		case core.ModAlt:
			modifierStrs = append(modifierStrs, "alt")
		case core.ModShift:
			modifierStrs = append(modifierStrs, "shift")
		case core.ModWin:
			modifierStrs = append(modifierStrs, "cmd") // robotgo中使用cmd表示Windows键
		}
	}

	// 执行组合键
	if len(modifierStrs) > 0 {
		// 转换为interface{}切片
		modifiers := make([]interface{}, len(modifierStrs))
		for i, mod := range modifierStrs {
			modifiers[i] = mod
		}
		robotgo.KeyTap(key, modifiers...)
	} else {
		robotgo.KeyTap(key)
	}

	hotkeyStr := strings.Join(modifierStrs, "+")
	if hotkeyStr != "" {
		hotkeyStr += "+" + key
	} else {
		hotkeyStr = key
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("组合键: %s", hotkeyStr),
		map[string]interface{}{
			"modifiers": modifierStrs,
			"key":       key,
			"hotkey":    hotkeyStr,
		},
	)
	result.SetDuration(start)
	return result
}

// Copy 复制操作 (Ctrl+C)
func (k *Keyboard) Copy() *core.OperationResult {
	return k.Hotkey([]core.KeyModifier{core.ModCtrl}, "c")
}

// Paste 粘贴操作 (Ctrl+V)
func (k *Keyboard) Paste() *core.OperationResult {
	return k.Hotkey([]core.KeyModifier{core.ModCtrl}, "v")
}

// SelectAll 全选操作 (Ctrl+A)
func (k *Keyboard) SelectAll() *core.OperationResult {
	return k.Hotkey([]core.KeyModifier{core.ModCtrl}, "a")
}

// Cut 剪切操作 (Ctrl+X)
func (k *Keyboard) Cut() *core.OperationResult {
	return k.Hotkey([]core.KeyModifier{core.ModCtrl}, "x")
}

// Undo 撤销操作 (Ctrl+Z)
func (k *Keyboard) Undo() *core.OperationResult {
	return k.Hotkey([]core.KeyModifier{core.ModCtrl}, "z")
}

// Redo 重做操作 (Ctrl+Y)
func (k *Keyboard) Redo() *core.OperationResult {
	return k.Hotkey([]core.KeyModifier{core.ModCtrl}, "y")
}

// Save 保存操作 (Ctrl+S)
func (k *Keyboard) Save() *core.OperationResult {
	return k.Hotkey([]core.KeyModifier{core.ModCtrl}, "s")
}

// Find 查找操作 (Ctrl+F)
func (k *Keyboard) Find() *core.OperationResult {
	return k.Hotkey([]core.KeyModifier{core.ModCtrl}, "f")
}

// AltTab 切换窗口 (Alt+Tab)
func (k *Keyboard) AltTab() *core.OperationResult {
	return k.Hotkey([]core.KeyModifier{core.ModAlt}, "tab")
}

// Enter 回车键
func (k *Keyboard) Enter() *core.OperationResult {
	return k.KeyPress("enter")
}

// Escape ESC键
func (k *Keyboard) Escape() *core.OperationResult {
	return k.KeyPress("escape")
}

// Tab Tab键
func (k *Keyboard) Tab() *core.OperationResult {
	return k.KeyPress("tab")
}

// Space 空格键
func (k *Keyboard) Space() *core.OperationResult {
	return k.KeyPress("space")
}

// Backspace 退格键
func (k *Keyboard) Backspace() *core.OperationResult {
	return k.KeyPress("backspace")
}

// Delete Delete键
func (k *Keyboard) Delete() *core.OperationResult {
	return k.KeyPress("delete")
}

// ArrowUp 上箭头键
func (k *Keyboard) ArrowUp() *core.OperationResult {
	return k.KeyPress("up")
}

// ArrowDown 下箭头键
func (k *Keyboard) ArrowDown() *core.OperationResult {
	return k.KeyPress("down")
}

// ArrowLeft 左箭头键
func (k *Keyboard) ArrowLeft() *core.OperationResult {
	return k.KeyPress("left")
}

// ArrowRight 右箭头键
func (k *Keyboard) ArrowRight() *core.OperationResult {
	return k.KeyPress("right")
}

// Home Home键
func (k *Keyboard) Home() *core.OperationResult {
	return k.KeyPress("home")
}

// End End键
func (k *Keyboard) End() *core.OperationResult {
	return k.KeyPress("end")
}

// PageUp Page Up键
func (k *Keyboard) PageUp() *core.OperationResult {
	return k.KeyPress("pageup")
}

// PageDown Page Down键
func (k *Keyboard) PageDown() *core.OperationResult {
	return k.KeyPress("pagedown")
}
