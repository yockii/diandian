package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"diandian/background/automation"
	"diandian/background/automation/core"
	"diandian/background/domain"
)

// TaskExecutionResult 任务执行结果
type TaskExecutionResult struct {
	TaskID         uint                   `json:"task_id"`
	Success        bool                   `json:"success"`
	Message        string                 `json:"message"`
	CompletedSteps int                    `json:"completed_steps"`
	TotalSteps     int                    `json:"total_steps"`
	Data           map[string]interface{} `json:"data"`
	Error          string                 `json:"error,omitempty"`
	Duration       time.Duration          `json:"duration"`
	StartTime      time.Time              `json:"start_time"`
}

// StepExecutionResult 步骤执行结果
type StepExecutionResult struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
	Error   string                 `json:"error,omitempty"`
}

// EnhancedTaskExecutionEngine 增强的任务执行引擎，支持两阶段架构
type EnhancedTaskExecutionEngine struct {
	automationService *AutomationService
	llmService        *LLMService
	engine            automation.AutomationEngine
}

// NewEnhancedTaskExecutionEngine 创建增强的任务执行引擎
func NewEnhancedTaskExecutionEngine(automationService *AutomationService) *EnhancedTaskExecutionEngine {
	return &EnhancedTaskExecutionEngine{
		automationService: automationService,
		llmService:        &LLMService{},
		engine:            automationService.engine,
	}
}

// ExecuteTaskDecomposition 执行任务分解结果
func (e *EnhancedTaskExecutionEngine) ExecuteTaskDecomposition(ctx context.Context, taskID uint, decomposition *domain.AutomationTaskDecomposition) *TaskExecutionResult {
	startTime := time.Now()

	result := &TaskExecutionResult{
		TaskID:         taskID,
		Success:        false,
		TotalSteps:     len(decomposition.Steps),
		CompletedSteps: 0,
		StartTime:      startTime,
	}

	slog.Info("开始执行任务分解", "task_id", taskID, "step_count", len(decomposition.Steps))

	// 发送任务开始事件
	e.automationService.sendEvent(AutomationEvent{
		Type:    "task_started",
		TaskID:  taskID,
		Message: "增强任务执行开始",
		Data: map[string]interface{}{
			"task_type":  decomposition.TaskType,
			"step_count": len(decomposition.Steps),
			"risk_level": decomposition.RiskLevel,
		},
	})

	// 逐步执行
	for i, stepPlan := range decomposition.Steps {
		select {
		case <-ctx.Done():
			result.Message = "任务被取消"
			result.Error = "context cancelled"
			return result
		default:
		}

		slog.Info("执行步骤", "step", i+1, "type", stepPlan.Type, "description", stepPlan.Description)

		// 发送步骤开始事件
		e.automationService.sendEvent(AutomationEvent{
			Type:    "step_started",
			TaskID:  taskID,
			Message: fmt.Sprintf("执行步骤 %d: %s", i+1, stepPlan.Description),
			Data: map[string]interface{}{
				"step_index":               i,
				"step_type":                stepPlan.Type,
				"requires_screen_analysis": stepPlan.RequiresScreenAnalysis,
			},
		})

		// 执行单个步骤
		stepResult := e.executeStepPlan(ctx, &stepPlan)

		if !stepResult.Success {
			if stepPlan.Optional {
				slog.Warn("可选步骤失败，继续执行", "step", i+1, "error", stepResult.Error)
				// 发送步骤跳过事件
				e.automationService.sendEvent(AutomationEvent{
					Type:    "step_skipped",
					TaskID:  taskID,
					Message: fmt.Sprintf("步骤 %d 失败但为可选步骤，已跳过", i+1),
					Data: map[string]interface{}{
						"step_index": i,
						"error":      stepResult.Error,
					},
				})
			} else {
				result.Message = fmt.Sprintf("步骤 %d 执行失败: %s", i+1, stepResult.Error)
				result.Error = stepResult.Error
				result.CompletedSteps = i
				result.Duration = time.Since(startTime)

				// 发送步骤失败事件
				e.automationService.sendEvent(AutomationEvent{
					Type:    "step_failed",
					TaskID:  taskID,
					Message: fmt.Sprintf("步骤 %d 执行失败", i+1),
					Data: map[string]interface{}{
						"step_index": i,
						"error":      stepResult.Error,
					},
				})

				return result
			}
		} else {
			// 发送步骤完成事件
			e.automationService.sendEvent(AutomationEvent{
				Type:    "step_completed",
				TaskID:  taskID,
				Message: fmt.Sprintf("步骤 %d 执行成功", i+1),
				Data: map[string]interface{}{
					"step_index": i,
					"result":     stepResult.Data,
				},
			})
		}

		result.CompletedSteps = i + 1
	}

	// 任务完成
	result.Success = true
	result.Message = "增强任务执行完成"
	result.Duration = time.Since(startTime)

	// 发送任务完成事件
	e.automationService.sendEvent(AutomationEvent{
		Type:    "task_completed",
		TaskID:  taskID,
		Message: "增强任务执行完成",
		Data: map[string]interface{}{
			"step_count":      len(decomposition.Steps),
			"completed_steps": result.CompletedSteps,
			"duration_ms":     result.Duration.Milliseconds(),
		},
	})

	slog.Info("任务执行完成", "task_id", taskID, "duration", result.Duration)
	return result
}

// executeStepPlan 执行单个步骤计划
func (e *EnhancedTaskExecutionEngine) executeStepPlan(ctx context.Context, stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{
		Success: false,
	}

	// 如果需要屏幕分析，先进行截屏和分析
	var screenAnalysis *domain.VisualAnalysisResponse
	if stepPlan.RequiresScreenAnalysis {
		imageData, screenshotResult := e.engine.Screenshot()
		if !screenshotResult.Success {
			result.Error = fmt.Sprintf("截屏失败: %s", screenshotResult.Error)
			return result
		}

		// 调用视觉分析
		analysis, err := e.llmService.AnalyzeScreenshot(imageData, stepPlan.Context)
		if err != nil {
			slog.Warn("视觉分析失败，使用默认策略", "error", err)
			// 不返回错误，继续执行，但没有屏幕分析结果
		} else {
			screenAnalysis = analysis
		}
	}

	// 根据步骤类型生成具体操作并执行
	switch stepPlan.Type {
	case "click":
		return e.executeClickStep(stepPlan, screenAnalysis)
	case "type":
		return e.executeTypeStep(stepPlan)
	case "launch_app":
		return e.executeLaunchAppStep(stepPlan)
	case "file":
		return e.executeFileStep(stepPlan)
	case "screenshot":
		return e.executeScreenshotStep(stepPlan)
	case "clipboard":
		return e.executeClipboardStep(stepPlan)
	case "wait":
		return e.executeWaitStep(stepPlan)
	case "key_press":
		return e.executeKeyPressStep(stepPlan)
	default:
		result.Error = fmt.Sprintf("不支持的步骤类型: %s", stepPlan.Type)
		return result
	}
}

// executeClickStep 执行点击步骤
func (e *EnhancedTaskExecutionEngine) executeClickStep(stepPlan *domain.AutomationStepPlan, screenAnalysis *domain.VisualAnalysisResponse) *StepExecutionResult {
	result := &StepExecutionResult{Success: false}

	// 生成具体的点击操作
	clickOp, err := e.llmService.GenerateClickOperation(stepPlan.Context, screenAnalysis)
	if err != nil {
		result.Error = fmt.Sprintf("生成点击操作失败: %v", err)
		return result
	}

	// 执行点击操作
	opResult := e.engine.Click(clickOp.X, clickOp.Y, core.MouseButton(clickOp.Button))
	if !opResult.Success {
		result.Error = opResult.Error
		return result
	}

	result.Success = true
	result.Data = map[string]interface{}{
		"x":      clickOp.X,
		"y":      clickOp.Y,
		"button": clickOp.Button,
	}
	return result
}

// executeTypeStep 执行输入步骤
func (e *EnhancedTaskExecutionEngine) executeTypeStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{Success: false}

	// 生成具体的输入操作
	typeOp, err := e.llmService.GenerateTypeOperation(stepPlan.Context)
	if err != nil {
		result.Error = fmt.Sprintf("生成输入操作失败: %v", err)
		return result
	}

	// 执行输入操作
	opResult := e.engine.Type(typeOp.Text)
	if !opResult.Success {
		result.Error = opResult.Error
		return result
	}

	result.Success = true
	result.Data = map[string]interface{}{
		"text":   typeOp.Text,
		"length": len(typeOp.Text),
	}
	return result
}

// executeLaunchAppStep 执行启动应用步骤
func (e *EnhancedTaskExecutionEngine) executeLaunchAppStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{Success: false}

	// 从上下文中提取应用名称
	// 这里可以使用简单的字符串匹配或者调用LLM来解析
	appName := e.extractAppNameFromContext(stepPlan.Context)
	if appName == "" {
		result.Error = "无法从上下文中提取应用名称"
		return result
	}

	// 执行启动应用操作
	opResult := e.engine.LaunchApp(&core.AppInfo{Name: appName})
	if !opResult.Success {
		result.Error = opResult.Error
		return result
	}

	result.Success = true
	result.Data = map[string]interface{}{
		"app_name": appName,
	}
	return result
}

// executeFileStep 执行文件操作步骤
func (e *EnhancedTaskExecutionEngine) executeFileStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{Success: false}

	// 生成具体的文件操作
	fileOp, err := e.llmService.GenerateFileOperation(stepPlan.Context)
	if err != nil {
		result.Error = fmt.Sprintf("生成文件操作失败: %v", err)
		return result
	}

	// 根据操作类型执行相应的文件操作
	var opResult *core.OperationResult
	switch fileOp.Operation {
	case "create":
		opResult = e.engine.CreateFile(fileOp.SourcePath, []byte(fileOp.Content))
	case "delete":
		opResult = e.engine.DeleteFile(fileOp.SourcePath)
	case "move":
		opResult = e.engine.MoveFile(fileOp.SourcePath, fileOp.TargetPath)
	case "copy":
		opResult = e.engine.CopyFile(fileOp.SourcePath, fileOp.TargetPath)
	default:
		result.Error = fmt.Sprintf("不支持的文件操作: %s", fileOp.Operation)
		return result
	}

	if !opResult.Success {
		result.Error = opResult.Error
		return result
	}

	result.Success = true
	result.Data = map[string]interface{}{
		"operation":      fileOp.Operation,
		"source_path":    fileOp.SourcePath,
		"target_path":    fileOp.TargetPath,
		"content_length": len(fileOp.Content),
	}
	return result
}

// executeScreenshotStep 执行截屏步骤
func (e *EnhancedTaskExecutionEngine) executeScreenshotStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{Success: false}

	// 从上下文中提取路径，或使用默认路径
	path := e.extractPathFromContext(stepPlan.Context)
	if path == "" {
		path = fmt.Sprintf("screenshot_%d.png", time.Now().Unix())
	}

	// 执行截屏操作
	imageData, opResult := e.engine.Screenshot()
	if !opResult.Success {
		result.Error = opResult.Error
		return result
	}

	// 保存截屏
	saveResult := e.engine.CreateFile(path, imageData)
	if !saveResult.Success {
		result.Error = fmt.Sprintf("保存截屏失败: %s", saveResult.Error)
		return result
	}

	result.Success = true
	result.Data = map[string]interface{}{
		"path": path,
		"size": len(imageData),
	}
	return result
}

// executeClipboardStep 执行剪贴板步骤
func (e *EnhancedTaskExecutionEngine) executeClipboardStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{Success: false}

	// 从上下文中判断是获取还是设置剪贴板
	if e.isGetClipboardOperation(stepPlan.Context) {
		text, opResult := e.engine.GetClipboard()
		if !opResult.Success {
			result.Error = opResult.Error
			return result
		}
		result.Success = true
		result.Data = map[string]interface{}{
			"operation": "get",
			"text":      text,
			"length":    len(text),
		}
	} else {
		// 设置剪贴板
		text := e.extractTextFromContext(stepPlan.Context)
		opResult := e.engine.SetClipboard(text)
		if !opResult.Success {
			result.Error = opResult.Error
			return result
		}
		result.Success = true
		result.Data = map[string]interface{}{
			"operation": "set",
			"text":      text,
			"length":    len(text),
		}
	}

	return result
}

// executeWaitStep 执行等待步骤
func (e *EnhancedTaskExecutionEngine) executeWaitStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{Success: false}

	// 从上下文中提取等待时间
	duration := e.extractDurationFromContext(stepPlan.Context)
	if duration <= 0 {
		duration = 1000 // 默认1秒
	}

	// 执行等待操作
	opResult := e.engine.Wait(duration)
	if !opResult.Success {
		result.Error = opResult.Error
		return result
	}

	result.Success = true
	result.Data = map[string]interface{}{
		"duration_ms": duration,
	}
	return result
}

// executeKeyPressStep 执行按键步骤
func (e *EnhancedTaskExecutionEngine) executeKeyPressStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{Success: false}

	// 从上下文中提取按键信息
	key, modifiers := e.extractKeyFromContext(stepPlan.Context)
	if key == "" {
		result.Error = "无法从上下文中提取按键信息"
		return result
	}

	// 执行按键操作
	var opResult *core.OperationResult
	if len(modifiers) > 0 {
		// 转换修饰键类型
		coreModifiers := make([]core.KeyModifier, len(modifiers))
		for i, mod := range modifiers {
			coreModifiers[i] = core.KeyModifier(mod)
		}
		opResult = e.engine.Hotkey(coreModifiers, key)
	} else {
		opResult = e.engine.KeyPress(key)
	}

	if !opResult.Success {
		result.Error = opResult.Error
		return result
	}

	result.Success = true
	result.Data = map[string]interface{}{
		"key":       key,
		"modifiers": modifiers,
	}
	return result
}

// 辅助方法：从上下文中提取信息
func (e *EnhancedTaskExecutionEngine) extractAppNameFromContext(context string) string {
	// 简单的字符串匹配，实际项目中可以使用更复杂的解析逻辑
	// 或者调用LLM来解析
	return "notepad" // 示例：默认返回记事本
}

func (e *EnhancedTaskExecutionEngine) extractPathFromContext(context string) string {
	// 从上下文中提取文件路径
	return "" // 返回空字符串使用默认路径
}

func (e *EnhancedTaskExecutionEngine) isGetClipboardOperation(context string) bool {
	// 判断是否是获取剪贴板操作
	return false // 默认为设置操作
}

func (e *EnhancedTaskExecutionEngine) extractTextFromContext(context string) string {
	// 从上下文中提取文本内容
	return context // 简单实现：直接使用上下文作为文本
}

func (e *EnhancedTaskExecutionEngine) extractDurationFromContext(context string) int {
	// 从上下文中提取等待时间
	return 1000 // 默认1秒
}

func (e *EnhancedTaskExecutionEngine) extractKeyFromContext(context string) (string, []string) {
	// 从上下文中提取按键和修饰键
	return "enter", []string{} // 默认回车键
}
