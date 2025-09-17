package mouse

import (
	"fmt"
	"time"

	"diandian/background/automation/core"
	"github.com/go-vgo/robotgo"
)

// Mouse 鼠标操作实现
type Mouse struct{}

// NewMouse 创建鼠标操作实例
func NewMouse() *Mouse {
	return &Mouse{}
}

// Click 点击指定位置
func (m *Mouse) Click(x, y int, button core.MouseButton) *core.OperationResult {
	start := time.Now()
	
	// 移动到指定位置
	robotgo.Move(x, y)
	time.Sleep(50 * time.Millisecond) // 短暂延迟确保移动完成
	
	// 执行点击
	var buttonStr string
	switch button {
	case core.LeftButton:
		buttonStr = "left"
	case core.RightButton:
		buttonStr = "right"
	case core.MiddleButton:
		buttonStr = "center"
	default:
		buttonStr = "left"
	}
	
	robotgo.Click(buttonStr)
	
	result := core.NewSuccessResult(
		fmt.Sprintf("点击位置 (%d, %d) 使用 %s 按键", x, y, button),
		map[string]interface{}{
			"x":      x,
			"y":      y,
			"button": button,
		},
	)
	result.SetDuration(start)
	return result
}

// DoubleClick 双击指定位置
func (m *Mouse) DoubleClick(x, y int) *core.OperationResult {
	start := time.Now()
	
	robotgo.Move(x, y)
	time.Sleep(50 * time.Millisecond)
	robotgo.Click("left", true) // true 表示双击
	
	result := core.NewSuccessResult(
		fmt.Sprintf("双击位置 (%d, %d)", x, y),
		map[string]interface{}{
			"x": x,
			"y": y,
		},
	)
	result.SetDuration(start)
	return result
}

// RightClick 右键点击指定位置
func (m *Mouse) RightClick(x, y int) *core.OperationResult {
	start := time.Now()
	
	robotgo.Move(x, y)
	time.Sleep(50 * time.Millisecond)
	robotgo.Click("right")
	
	result := core.NewSuccessResult(
		fmt.Sprintf("右键点击位置 (%d, %d)", x, y),
		map[string]interface{}{
			"x": x,
			"y": y,
		},
	)
	result.SetDuration(start)
	return result
}

// Drag 拖拽操作
func (m *Mouse) Drag(fromX, fromY, toX, toY int) *core.OperationResult {
	start := time.Now()
	
	// 移动到起始位置
	robotgo.Move(fromX, fromY)
	time.Sleep(50 * time.Millisecond)
	
	// 按下鼠标左键
	robotgo.Toggle("left", "down")
	time.Sleep(50 * time.Millisecond)
	
	// 拖拽到目标位置
	robotgo.Move(toX, toY)
	time.Sleep(100 * time.Millisecond)
	
	// 释放鼠标左键
	robotgo.Toggle("left", "up")
	
	result := core.NewSuccessResult(
		fmt.Sprintf("拖拽从 (%d, %d) 到 (%d, %d)", fromX, fromY, toX, toY),
		map[string]interface{}{
			"from_x": fromX,
			"from_y": fromY,
			"to_x":   toX,
			"to_y":   toY,
		},
	)
	result.SetDuration(start)
	return result
}

// Move 移动鼠标到指定位置
func (m *Mouse) Move(x, y int) *core.OperationResult {
	start := time.Now()
	
	robotgo.Move(x, y)
	
	result := core.NewSuccessResult(
		fmt.Sprintf("移动鼠标到 (%d, %d)", x, y),
		map[string]interface{}{
			"x": x,
			"y": y,
		},
	)
	result.SetDuration(start)
	return result
}

// GetPosition 获取当前鼠标位置
func (m *Mouse) GetPosition() (*core.Point, *core.OperationResult) {
	start := time.Now()
	
	x, y := robotgo.GetMousePos()
	point := &core.Point{X: x, Y: y}
	
	result := core.NewSuccessResult(
		fmt.Sprintf("获取鼠标位置 (%d, %d)", x, y),
		point,
	)
	result.SetDuration(start)
	return point, result
}

// Scroll 滚动操作
func (m *Mouse) Scroll(x, y int, direction string, clicks int) *core.OperationResult {
	start := time.Now()
	
	// 移动到指定位置
	robotgo.Move(x, y)
	time.Sleep(50 * time.Millisecond)
	
	// 执行滚动
	if direction == "up" {
		robotgo.Scroll(0, clicks)
	} else if direction == "down" {
		robotgo.Scroll(0, -clicks)
	} else if direction == "left" {
		robotgo.Scroll(-clicks, 0)
	} else if direction == "right" {
		robotgo.Scroll(clicks, 0)
	} else {
		return core.NewErrorResult("无效的滚动方向", fmt.Errorf("direction must be up, down, left, or right"))
	}
	
	result := core.NewSuccessResult(
		fmt.Sprintf("在位置 (%d, %d) 向 %s 滚动 %d 次", x, y, direction, clicks),
		map[string]interface{}{
			"x":         x,
			"y":         y,
			"direction": direction,
			"clicks":    clicks,
		},
	)
	result.SetDuration(start)
	return result
}

// SmoothMove 平滑移动鼠标
func (m *Mouse) SmoothMove(x, y int, duration time.Duration) *core.OperationResult {
	start := time.Now()
	
	currentX, currentY := robotgo.GetMousePos()
	steps := int(duration.Milliseconds() / 10) // 每10ms一步
	if steps < 1 {
		steps = 1
	}
	
	deltaX := float64(x-currentX) / float64(steps)
	deltaY := float64(y-currentY) / float64(steps)
	
	for i := 0; i < steps; i++ {
		newX := currentX + int(deltaX*float64(i))
		newY := currentY + int(deltaY*float64(i))
		robotgo.Move(newX, newY)
		time.Sleep(10 * time.Millisecond)
	}
	
	// 确保到达目标位置
	robotgo.Move(x, y)
	
	result := core.NewSuccessResult(
		fmt.Sprintf("平滑移动鼠标到 (%d, %d)", x, y),
		map[string]interface{}{
			"x":        x,
			"y":        y,
			"duration": duration.String(),
		},
	)
	result.SetDuration(start)
	return result
}
