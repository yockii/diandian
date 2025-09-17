package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"diandian/background/automation"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// AutomationService 自动化服务
type AutomationService struct {
	app    *application.App
	engine automation.AutomationEngine

	// 当前执行状态
	isRunning     bool
	currentTaskID uint

	// 事件通道
	eventChan chan AutomationEvent
}

// AutomationEvent 自动化事件
type AutomationEvent struct {
	Type    string      `json:"type"`
	TaskID  uint        `json:"task_id"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// AutomationRequest 自动化请求
type AutomationRequest struct {
	TaskID      uint                   `json:"task_id"`
	Instruction string                 `json:"instruction"`
	Steps       []AutomationStep       `json:"steps"`
	Options     map[string]interface{} `json:"options"`
}

// AutomationStep 自动化步骤
type AutomationStep struct {
	Type        string                 `json:"type"`        // "click", "type", "launch", "file", "screenshot", etc.
	Description string                 `json:"description"` // 步骤描述
	Parameters  map[string]interface{} `json:"parameters"`  // 步骤参数
	Expected    string                 `json:"expected"`    // 预期结果
}

// AutomationResponse 自动化响应
type AutomationResponse struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message"`
	TaskID    uint                   `json:"task_id"`
	StepIndex int                    `json:"step_index"`
	Data      map[string]interface{} `json:"data"`
	Error     string                 `json:"error,omitempty"`
}

// NewAutomationService 创建自动化服务
func NewAutomationService(app *application.App) *AutomationService {
	return &AutomationService{
		app:       app,
		engine:    automation.NewEngine(),
		eventChan: make(chan AutomationEvent, 100),
	}
}

// Initialize 初始化自动化服务
func (s *AutomationService) Initialize() error {
	result := s.engine.Initialize()
	if !result.Success {
		return fmt.Errorf("初始化自动化引擎失败: %s", result.Error)
	}

	// 启动事件处理协程
	go s.handleEvents()

	log.Println("自动化服务初始化成功")
	return nil
}

// Cleanup 清理资源
func (s *AutomationService) Cleanup() {
	if s.engine != nil {
		s.engine.Cleanup()
		s.engine = nil
	}
	if s.eventChan != nil {
		close(s.eventChan)
		s.eventChan = nil
	}
	s.isRunning = false
}

// ExecuteAutomationTask 执行自动化任务
func (s *AutomationService) ExecuteAutomationTask(ctx context.Context, request AutomationRequest) *AutomationResponse {
	if s.isRunning {
		return &AutomationResponse{
			Success: false,
			Message: "已有任务正在执行",
			TaskID:  request.TaskID,
			Error:   "automation_busy",
		}
	}

	s.isRunning = true
	s.currentTaskID = request.TaskID
	defer func() {
		s.isRunning = false
		s.currentTaskID = 0
	}()

	// 发送任务开始事件
	s.sendEvent(AutomationEvent{
		Type:    "task_started",
		TaskID:  request.TaskID,
		Message: "自动化任务开始执行",
		Data: map[string]interface{}{
			"instruction": request.Instruction,
			"step_count":  len(request.Steps),
		},
	})

	// 执行步骤
	for i, step := range request.Steps {
		select {
		case <-ctx.Done():
			return &AutomationResponse{
				Success:   false,
				Message:   "任务被取消",
				TaskID:    request.TaskID,
				StepIndex: i,
				Error:     "cancelled",
			}
		default:
		}

		// 发送步骤开始事件
		s.sendEvent(AutomationEvent{
			Type:    "step_started",
			TaskID:  request.TaskID,
			Message: fmt.Sprintf("执行步骤 %d: %s", i+1, step.Description),
			Data: map[string]interface{}{
				"step_index": i,
				"step_type":  step.Type,
			},
		})

		// 执行步骤
		result := s.executeStep(step)
		if !result.Success {
			// 发送步骤失败事件
			s.sendEvent(AutomationEvent{
				Type:    "step_failed",
				TaskID:  request.TaskID,
				Message: fmt.Sprintf("步骤 %d 执行失败: %s", i+1, result.Error),
				Data: map[string]interface{}{
					"step_index": i,
					"error":      result.Error,
				},
			})

			return &AutomationResponse{
				Success:   false,
				Message:   fmt.Sprintf("步骤 %d 执行失败", i+1),
				TaskID:    request.TaskID,
				StepIndex: i,
				Error:     result.Error,
			}
		}

		// 发送步骤完成事件
		s.sendEvent(AutomationEvent{
			Type:    "step_completed",
			TaskID:  request.TaskID,
			Message: fmt.Sprintf("步骤 %d 执行成功", i+1),
			Data: map[string]interface{}{
				"step_index": i,
				"result":     result.Data,
			},
		})

		// 步骤间延迟
		time.Sleep(500 * time.Millisecond)
	}

	// 发送任务完成事件
	s.sendEvent(AutomationEvent{
		Type:    "task_completed",
		TaskID:  request.TaskID,
		Message: "自动化任务执行完成",
		Data: map[string]interface{}{
			"step_count": len(request.Steps),
		},
	})

	return &AutomationResponse{
		Success:   true,
		Message:   "自动化任务执行完成",
		TaskID:    request.TaskID,
		StepIndex: len(request.Steps),
		Data: map[string]interface{}{
			"completed_steps": len(request.Steps),
		},
	}
}

// executeStep 执行单个步骤
func (s *AutomationService) executeStep(step AutomationStep) *automation.OperationResult {
	switch step.Type {
	case "click":
		return s.executeClickStep(step)
	case "type":
		return s.executeTypeStep(step)
	case "launch":
		return s.executeLaunchStep(step)
	case "file":
		return s.executeFileStep(step)
	case "screenshot":
		return s.executeScreenshotStep(step)
	case "wait":
		return s.executeWaitStep(step)
	case "key":
		return s.executeKeyStep(step)
	case "clipboard":
		return s.executeClipboardStep(step)
	default:
		return automation.NewErrorResult(
			fmt.Sprintf("不支持的步骤类型: %s", step.Type),
			fmt.Errorf("unsupported step type"),
		)
	}
}

// executeClickStep 执行点击步骤
func (s *AutomationService) executeClickStep(step AutomationStep) *automation.OperationResult {
	x, ok1 := step.Parameters["x"].(float64)
	y, ok2 := step.Parameters["y"].(float64)
	if !ok1 || !ok2 {
		return automation.NewErrorResult("点击步骤缺少坐标参数", fmt.Errorf("missing coordinates"))
	}

	button := automation.LeftButton
	if buttonStr, ok := step.Parameters["button"].(string); ok {
		switch buttonStr {
		case "right":
			button = automation.RightButton
		case "middle":
			button = automation.MiddleButton
		}
	}

	return s.engine.Click(int(x), int(y), button)
}

// executeTypeStep 执行输入步骤
func (s *AutomationService) executeTypeStep(step AutomationStep) *automation.OperationResult {
	text, ok := step.Parameters["text"].(string)
	if !ok {
		return automation.NewErrorResult("输入步骤缺少文本参数", fmt.Errorf("missing text parameter"))
	}

	return s.engine.Type(text)
}

// executeLaunchStep 执行启动应用步骤
func (s *AutomationService) executeLaunchStep(step AutomationStep) *automation.OperationResult {
	app, ok := step.Parameters["app"].(string)
	if !ok {
		return automation.NewErrorResult("启动步骤缺少应用参数", fmt.Errorf("missing app parameter"))
	}

	return s.engine.Launch(app)
}

// executeFileStep 执行文件操作步骤
func (s *AutomationService) executeFileStep(step AutomationStep) *automation.OperationResult {
	operation, ok := step.Parameters["operation"].(string)
	if !ok {
		return automation.NewErrorResult("文件步骤缺少操作参数", fmt.Errorf("missing operation parameter"))
	}

	switch operation {
	case "create":
		path, ok1 := step.Parameters["path"].(string)
		content, ok2 := step.Parameters["content"].(string)
		if !ok1 || !ok2 {
			return automation.NewErrorResult("创建文件缺少路径或内容参数", fmt.Errorf("missing path or content"))
		}
		return s.engine.CreateFile(path, []byte(content))

	case "delete":
		path, ok := step.Parameters["path"].(string)
		if !ok {
			return automation.NewErrorResult("删除文件缺少路径参数", fmt.Errorf("missing path parameter"))
		}
		return s.engine.DeleteFile(path)

	default:
		return automation.NewErrorResult(
			fmt.Sprintf("不支持的文件操作: %s", operation),
			fmt.Errorf("unsupported file operation"),
		)
	}
}

// executeScreenshotStep 执行截屏步骤
func (s *AutomationService) executeScreenshotStep(step AutomationStep) *automation.OperationResult {
	path, ok := step.Parameters["path"].(string)
	if !ok {
		// 如果没有指定路径，生成默认路径
		path = fmt.Sprintf("screenshot_%d.png", time.Now().Unix())
	}

	imageData, result := s.engine.Screenshot()
	if !result.Success {
		return result
	}

	return s.engine.CreateFile(path, imageData)
}

// executeWaitStep 执行等待步骤
func (s *AutomationService) executeWaitStep(step AutomationStep) *automation.OperationResult {
	duration, ok := step.Parameters["duration"].(float64)
	if !ok {
		duration = 1000 // 默认等待1秒
	}

	return s.engine.Wait(int(duration))
}

// executeKeyStep 执行按键步骤
func (s *AutomationService) executeKeyStep(step AutomationStep) *automation.OperationResult {
	key, ok := step.Parameters["key"].(string)
	if !ok {
		return automation.NewErrorResult("按键步骤缺少按键参数", fmt.Errorf("missing key parameter"))
	}

	return s.engine.KeyPress(key)
}

// executeClipboardStep 执行剪贴板步骤
func (s *AutomationService) executeClipboardStep(step AutomationStep) *automation.OperationResult {
	operation, ok := step.Parameters["operation"].(string)
	if !ok {
		return automation.NewErrorResult("剪贴板步骤缺少操作参数", fmt.Errorf("missing operation parameter"))
	}

	switch operation {
	case "get":
		text, result := s.engine.GetClipboard()
		if result.Success {
			result.Data = map[string]interface{}{
				"text": text,
			}
		}
		return result

	case "set":
		text, ok := step.Parameters["text"].(string)
		if !ok {
			return automation.NewErrorResult("设置剪贴板缺少文本参数", fmt.Errorf("missing text parameter"))
		}
		return s.engine.SetClipboard(text)

	default:
		return automation.NewErrorResult(
			fmt.Sprintf("不支持的剪贴板操作: %s", operation),
			fmt.Errorf("unsupported clipboard operation"),
		)
	}
}

// sendEvent 发送事件
func (s *AutomationService) sendEvent(event AutomationEvent) {
	select {
	case s.eventChan <- event:
	default:
		log.Printf("事件通道已满，丢弃事件: %s", event.Type)
	}
}

// handleEvents 处理事件
func (s *AutomationService) handleEvents() {
	for event := range s.eventChan {
		// 将事件发送到前端
		eventData, _ := json.Marshal(event)
		// 注意：Wails v3的事件发送API可能不同，这里先记录日志
		log.Printf("自动化事件: %s", string(eventData))
		// TODO: 实现正确的事件发送方式
		// s.app.Events.Emit("automation-event", string(eventData))
	}
}

// GetStatus 获取自动化服务状态
func (s *AutomationService) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"is_running":       s.isRunning,
		"current_task_id":  s.currentTaskID,
		"engine_available": s.engine != nil,
	}
}

// GetEngine 获取自动化引擎（用于高级功能）
func (s *AutomationService) GetEngine() automation.AutomationEngine {
	return s.engine
}

// StopCurrentTask 停止当前任务
func (s *AutomationService) StopCurrentTask() *AutomationResponse {
	if !s.isRunning {
		return &AutomationResponse{
			Success: false,
			Message: "没有正在执行的任务",
			Error:   "no_running_task",
		}
	}

	// 发送停止事件
	s.sendEvent(AutomationEvent{
		Type:    "task_stopped",
		TaskID:  s.currentTaskID,
		Message: "任务被用户停止",
	})

	return &AutomationResponse{
		Success: true,
		Message: "任务停止请求已发送",
		TaskID:  s.currentTaskID,
	}
}

// ExecuteStep 公开的步骤执行方法
func (s *AutomationService) ExecuteStep(step AutomationStep) *automation.OperationResult {
	return s.executeStep(step)
}
