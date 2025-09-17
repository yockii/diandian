package screen

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log/slog"
	"os"
	"time"

	"diandian/background/automation/core"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
)

// Screen 屏幕操作实现
type Screen struct{}

// NewScreen 创建屏幕操作实例
func NewScreen() *Screen {
	return &Screen{}
}

// Screenshot 截取屏幕（智能多屏幕支持）
func (s *Screen) Screenshot() ([]byte, *core.OperationResult) {
	start := time.Now()

	// 尝试智能截图：优先截取活动窗口所在的屏幕
	imageData, result := s.SmartScreenshot()
	if result.Success {
		result.SetDuration(start)
		return imageData, result
	}

	// 如果智能截图失败，回退到主屏幕截图
	return s.ScreenshotPrimary()
}

// ScreenshotPrimary 截取主屏幕
func (s *Screen) ScreenshotPrimary() ([]byte, *core.OperationResult) {
	start := time.Now()

	// 获取主屏幕尺寸
	width, height := robotgo.GetScreenSize()

	// 截取主屏幕
	img, err := screenshot.CaptureRect(image.Rect(0, 0, width, height))
	if err != nil {
		result := core.NewErrorResult("主屏幕截屏失败", err)
		result.SetDuration(start)
		return nil, result
	}

	// 将图像转换为PNG字节数组
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		result := core.NewErrorResult("图像编码失败", err)
		result.SetDuration(start)
		return nil, result
	}

	imageData := buf.Bytes()

	result := core.NewSuccessResult(
		fmt.Sprintf("主屏幕截屏成功 (%dx%d)", width, height),
		map[string]interface{}{
			"width":  width,
			"height": height,
			"size":   len(imageData),
			"screen": "primary",
		},
	)
	result.SetDuration(start)
	return imageData, result
}

// SmartScreenshot 智能截图：尝试截取活动窗口所在的屏幕
func (s *Screen) SmartScreenshot() ([]byte, *core.OperationResult) {
	start := time.Now()

	// 获取显示器数量
	numDisplays := screenshot.NumActiveDisplays()
	if numDisplays == 0 {
		result := core.NewErrorResult("无法获取显示器信息", fmt.Errorf("no active displays found"))
		result.SetDuration(start)
		return nil, result
	}

	// 获取所有显示器边界
	var displays []image.Rectangle
	for i := 0; i < numDisplays; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		displays = append(displays, bounds)
	}

	// 如果只有一个显示器，直接截取
	if len(displays) == 1 {
		return s.screenshotDisplay(displays[0], 0)
	}

	// 多显示器环境：尝试找到活动窗口所在的屏幕
	activeDisplayIndex := s.findActiveWindowDisplay(displays)
	if activeDisplayIndex >= 0 {
		return s.screenshotDisplay(displays[activeDisplayIndex], activeDisplayIndex)
	}

	// 如果找不到活动窗口，截取最大的显示器
	largestDisplayIndex := s.findLargestDisplay(displays)
	return s.screenshotDisplay(displays[largestDisplayIndex], largestDisplayIndex)
}

// screenshotDisplay 截取指定显示器
func (s *Screen) screenshotDisplay(display image.Rectangle, displayIndex int) ([]byte, *core.OperationResult) {
	start := time.Now()

	// 截取指定显示器
	img, err := screenshot.CaptureRect(display)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("显示器截屏失败 (显示器 %d)", displayIndex),
			err,
		)
		result.SetDuration(start)
		return nil, result
	}

	// 将图像转换为PNG字节数组
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		result := core.NewErrorResult("图像编码失败", err)
		result.SetDuration(start)
		return nil, result
	}

	imageData := buf.Bytes()
	width := display.Dx()
	height := display.Dy()

	result := core.NewSuccessResult(
		fmt.Sprintf("显示器截屏成功 (%dx%d, 显示器 %d)", width, height, displayIndex),
		map[string]interface{}{
			"width":         width,
			"height":        height,
			"size":          len(imageData),
			"display_index": displayIndex,
			"display_rect":  display,
		},
	)
	result.SetDuration(start)
	return imageData, result
}

// findActiveWindowDisplay 找到活动窗口所在的显示器索引
func (s *Screen) findActiveWindowDisplay(displays []image.Rectangle) int {
	// 获取鼠标位置作为活动区域的参考
	mouseX, mouseY := robotgo.GetMousePos()
	mousePoint := image.Point{X: mouseX, Y: mouseY}

	// 检查鼠标在哪个显示器上
	for i, display := range displays {
		if mousePoint.In(display) {
			return i
		}
	}

	return -1 // 未找到
}

// findLargestDisplay 找到最大的显示器索引
func (s *Screen) findLargestDisplay(displays []image.Rectangle) int {
	if len(displays) == 0 {
		return -1
	}

	largestIndex := 0
	largestArea := displays[0].Dx() * displays[0].Dy()

	for i, display := range displays[1:] {
		area := display.Dx() * display.Dy()
		if area > largestArea {
			largestIndex = i + 1 // 因为从displays[1:]开始，所以要+1
			largestArea = area
		}
	}

	return largestIndex
}

// CaptureAllDisplays 截取所有显示器（内存操作）
func (s *Screen) CaptureAllDisplays() ([]core.DisplayCapture, *core.OperationResult) {
	start := time.Now()

	numDisplays := screenshot.NumActiveDisplays()
	if numDisplays == 0 {
		result := core.NewErrorResult("无法获取显示器信息", fmt.Errorf("no active displays found"))
		result.SetDuration(start)
		return nil, result
	}

	var captures []core.DisplayCapture
	for i := 0; i < numDisplays; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		// 截取显示器
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			// 如果某个显示器截取失败，记录错误但继续其他显示器
			continue
		}

		// 转换为PNG字节数组
		var buf bytes.Buffer
		err = png.Encode(&buf, img)
		if err != nil {
			continue
		}

		capture := core.DisplayCapture{
			Index:     i,
			Bounds:    bounds,
			ImageData: buf.Bytes(),
			Width:     bounds.Dx(),
			Height:    bounds.Dy(),
		}
		captures = append(captures, capture)
	}

	if len(captures) == 0 {
		result := core.NewErrorResult("所有显示器截取失败", fmt.Errorf("failed to capture any display"))
		result.SetDuration(start)
		return nil, result
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("成功截取 %d 个显示器", len(captures)),
		map[string]interface{}{
			"total_displays":    numDisplays,
			"captured_displays": len(captures),
			"captures":          captures,
		},
	)
	result.SetDuration(start)
	return captures, result
}

// CaptureSpecificDisplay 截取指定显示器（内存操作）
func (s *Screen) CaptureSpecificDisplay(displayIndex int) (*core.DisplayCapture, *core.OperationResult) {
	start := time.Now()

	numDisplays := screenshot.NumActiveDisplays()
	if displayIndex < 0 || displayIndex >= numDisplays {
		result := core.NewErrorResult(
			fmt.Sprintf("显示器索引无效: %d (总共 %d 个显示器)", displayIndex, numDisplays),
			fmt.Errorf("invalid display index"),
		)
		result.SetDuration(start)
		return nil, result
	}

	bounds := screenshot.GetDisplayBounds(displayIndex)

	// 截取显示器
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("显示器 %d 截取失败", displayIndex),
			err,
		)
		result.SetDuration(start)
		return nil, result
	}

	// 转换为PNG字节数组
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		result := core.NewErrorResult("图像编码失败", err)
		result.SetDuration(start)
		return nil, result
	}

	capture := &core.DisplayCapture{
		Index:     displayIndex,
		Bounds:    bounds,
		ImageData: buf.Bytes(),
		Width:     bounds.Dx(),
		Height:    bounds.Dy(),
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("成功截取显示器 %d (%dx%d)", displayIndex, capture.Width, capture.Height),
		capture,
	)
	result.SetDuration(start)
	return capture, result
}

// CaptureToMemory 截取屏幕到内存（不保存文件）
func (s *Screen) CaptureToMemory() (*core.OperationResult, []byte, string) {
	start := time.Now()

	// 使用智能截图选择最佳显示器
	captures, result := s.CaptureAllDisplays()
	if !result.Success {
		return result, nil, ""
	}

	if len(captures) == 0 {
		result := core.NewErrorResult("没有可用的显示器", fmt.Errorf("no displays captured"))
		result.SetDuration(start)
		return result, nil, ""
	}

	// 选择最大的显示器或活动显示器
	var selectedCapture *core.DisplayCapture
	maxArea := 0

	for i := range captures {
		area := captures[i].Width * captures[i].Height
		if area > maxArea {
			maxArea = area
			selectedCapture = &captures[i]
		}
	}

	if selectedCapture == nil {
		result := core.NewErrorResult("无法选择显示器", fmt.Errorf("no suitable display found"))
		result.SetDuration(start)
		return result, nil, ""
	}

	// 转换为base64
	imageBase64 := base64.StdEncoding.EncodeToString(selectedCapture.ImageData)

	result = core.NewSuccessResult(
		fmt.Sprintf("成功截取显示器 %d 到内存 (%dx%d)",
			selectedCapture.Index, selectedCapture.Width, selectedCapture.Height),
		map[string]interface{}{
			"display_index": selectedCapture.Index,
			"width":         selectedCapture.Width,
			"height":        selectedCapture.Height,
			"data_size":     len(selectedCapture.ImageData),
		},
	)
	result.SetDuration(start)

	slog.Info("内存截图完成",
		"display", selectedCapture.Index,
		"size", fmt.Sprintf("%dx%d", selectedCapture.Width, selectedCapture.Height),
		"data_size", len(selectedCapture.ImageData))

	return result, selectedCapture.ImageData, imageBase64
}

// GetDisplayInfo 获取所有显示器信息
func (s *Screen) GetDisplayInfo() ([]map[string]interface{}, *core.OperationResult) {
	start := time.Now()

	numDisplays := screenshot.NumActiveDisplays()
	if numDisplays == 0 {
		result := core.NewErrorResult("无法获取显示器信息", fmt.Errorf("no active displays found"))
		result.SetDuration(start)
		return nil, result
	}

	var displayInfos []map[string]interface{}
	for i := 0; i < numDisplays; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		displayInfo := map[string]interface{}{
			"index":  i,
			"x":      bounds.Min.X,
			"y":      bounds.Min.Y,
			"width":  bounds.Dx(),
			"height": bounds.Dy(),
			"bounds": bounds,
		}
		displayInfos = append(displayInfos, displayInfo)
	}

	// 获取鼠标位置，标记活动显示器
	mouseX, mouseY := robotgo.GetMousePos()
	mousePoint := image.Point{X: mouseX, Y: mouseY}

	for i, info := range displayInfos {
		bounds := info["bounds"].(image.Rectangle)
		if mousePoint.In(bounds) {
			displayInfos[i]["is_active"] = true
		} else {
			displayInfos[i]["is_active"] = false
		}
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("获取到 %d 个显示器信息", numDisplays),
		map[string]interface{}{
			"count":    numDisplays,
			"displays": displayInfos,
			"mouse_x":  mouseX,
			"mouse_y":  mouseY,
		},
	)
	result.SetDuration(start)
	return displayInfos, result
}

// ScreenshotArea 截取指定区域
func (s *Screen) ScreenshotArea(rect core.Rect) ([]byte, *core.OperationResult) {
	start := time.Now()

	// 截取指定区域
	img, err := screenshot.CaptureRect(image.Rect(rect.X, rect.Y, rect.X+rect.Width, rect.Y+rect.Height))
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("截取区域失败: (%d,%d,%d,%d)", rect.X, rect.Y, rect.Width, rect.Height),
			err,
		)
		result.SetDuration(start)
		return nil, result
	}

	// 将图像转换为PNG字节数组
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		result := core.NewErrorResult("图像编码失败", err)
		result.SetDuration(start)
		return nil, result
	}

	imageData := buf.Bytes()

	result := core.NewSuccessResult(
		fmt.Sprintf("截取区域成功 (%dx%d)", rect.Width, rect.Height),
		map[string]interface{}{
			"x":      rect.X,
			"y":      rect.Y,
			"width":  rect.Width,
			"height": rect.Height,
			"size":   len(imageData),
		},
	)
	result.SetDuration(start)
	return imageData, result
}

// GetScreenSize 获取屏幕尺寸
func (s *Screen) GetScreenSize() (*core.Size, *core.OperationResult) {
	start := time.Now()

	width, height := robotgo.GetScreenSize()
	size := &core.Size{
		Width:  width,
		Height: height,
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("获取屏幕尺寸: %dx%d", width, height),
		size,
	)
	result.SetDuration(start)
	return size, result
}

// FindImage 在屏幕上查找图像
func (s *Screen) FindImage(templatePath string) (*core.Point, *core.OperationResult) {
	start := time.Now()

	// 注意：robotgo的图像识别API在不同版本中可能不同
	// 这里返回未实现错误，实际使用时需要根据具体版本调整
	result := core.NewErrorResult(
		"图像识别功能需要根据robotgo版本调整",
		fmt.Errorf("not implemented"),
	)
	result.SetDuration(start)
	return nil, result
}

// FindText 在屏幕上查找文本（OCR）
func (s *Screen) FindText(text string) (*core.Point, *core.OperationResult) {
	start := time.Now()

	// 注意：这里需要集成OCR功能，目前返回未实现错误
	result := core.NewErrorResult(
		"OCR文本识别功能尚未实现",
		fmt.Errorf("OCR not implemented"),
	)
	result.SetDuration(start)
	return nil, result
}

// SaveScreenshot 保存截屏到文件
func (s *Screen) SaveScreenshot(filePath string) *core.OperationResult {
	start := time.Now()

	// 截取屏幕
	imageData, result := s.Screenshot()
	if !result.Success {
		return result
	}

	// 保存到文件
	err := os.WriteFile(filePath, imageData, 0644)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("保存截屏失败: %s", filePath),
			err,
		)
		result.SetDuration(start)
		return result
	}

	result = core.NewSuccessResult(
		fmt.Sprintf("截屏已保存: %s", filePath),
		map[string]interface{}{
			"path": filePath,
			"size": len(imageData),
		},
	)
	result.SetDuration(start)
	return result
}

// SaveScreenshotArea 保存指定区域截屏到文件
func (s *Screen) SaveScreenshotArea(rect core.Rect, filePath string) *core.OperationResult {
	start := time.Now()

	// 截取指定区域
	imageData, result := s.ScreenshotArea(rect)
	if !result.Success {
		return result
	}

	// 保存到文件
	err := os.WriteFile(filePath, imageData, 0644)
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("保存区域截屏失败: %s", filePath),
			err,
		)
		result.SetDuration(start)
		return result
	}

	result = core.NewSuccessResult(
		fmt.Sprintf("区域截屏已保存: %s", filePath),
		map[string]interface{}{
			"path":   filePath,
			"x":      rect.X,
			"y":      rect.Y,
			"width":  rect.Width,
			"height": rect.Height,
			"size":   len(imageData),
		},
	)
	result.SetDuration(start)
	return result
}

// GetPixelColor 获取指定位置的像素颜色
func (s *Screen) GetPixelColor(x, y int) (string, *core.OperationResult) {
	start := time.Now()

	color := robotgo.GetPixelColor(x, y)

	result := core.NewSuccessResult(
		fmt.Sprintf("获取像素颜色: (%d, %d) = %s", x, y, color),
		map[string]interface{}{
			"x":     x,
			"y":     y,
			"color": color,
		},
	)
	result.SetDuration(start)
	return color, result
}

// WaitForImage 等待图像出现
func (s *Screen) WaitForImage(templatePath string, timeout time.Duration) (*core.Point, *core.OperationResult) {
	start := time.Now()

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		point, result := s.FindImage(templatePath)
		if result.Success {
			result.Message = fmt.Sprintf("等待图像出现成功: %s", templatePath)
			result.SetDuration(start)
			return point, result
		}

		time.Sleep(500 * time.Millisecond) // 每500ms检查一次
	}

	result := core.NewErrorResult(
		fmt.Sprintf("等待图像超时: %s", templatePath),
		fmt.Errorf("timeout waiting for image"),
	)
	result.SetDuration(start)
	return nil, result
}

// GetScreenInfo 获取屏幕信息
func (s *Screen) GetScreenInfo() (*core.ScreenInfo, *core.OperationResult) {
	start := time.Now()

	width, height := robotgo.GetScreenSize()

	// 获取DPI（这里使用默认值，实际可能需要系统API）
	dpi := 96 // Windows默认DPI

	screenInfo := &core.ScreenInfo{
		Width:  width,
		Height: height,
		DPI:    dpi,
	}

	result := core.NewSuccessResult(
		fmt.Sprintf("获取屏幕信息: %dx%d, DPI: %d", width, height, dpi),
		screenInfo,
	)
	result.SetDuration(start)
	return screenInfo, result
}
