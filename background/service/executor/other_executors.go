package executor

import (
	"fmt"

	"diandian/background/domain"
	"diandian/background/service"
)

// executeLaunchAppStep 执行启动应用步骤
func (e *EnhancedTaskExecutionEngine) executeLaunchAppStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{
		StepType: stepPlan.Type,
		Success:  false,
	}

	appIdentifier := e.extractAppNameFromContext(stepPlan.Context)
	if appIdentifier == "" {
		result.Error = "无法从上下文中提取应用名称"
		return result
	}

	// 使用智能应用启动器
	appLauncher := service.NewAppLauncher()
	err := appLauncher.LaunchApp(appIdentifier)
	if err != nil {
		result.Error = fmt.Sprintf("智能启动应用失败: %s", err.Error())

		// 如果智能启动失败，回退到原始方法
		step := service.AutomationStep{
			Type: "launch",
			Parameters: map[string]interface{}{
				"app": appIdentifier,
			},
		}

		opResult := e.automationService.ExecuteStep(step)
		if !opResult.Success {
			result.Error = fmt.Sprintf("启动应用失败 (智能启动和原始方法都失败): %s", opResult.Error)
			return result
		}
	}

	result.Success = true
	result.Message = fmt.Sprintf("成功启动应用: %s", appIdentifier)
	return result
}

// executeScreenshotStep 执行截屏步骤
func (e *EnhancedTaskExecutionEngine) executeScreenshotStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{
		StepType: stepPlan.Type,
		Success:  false,
	}

	path := e.extractPathFromContext(stepPlan.Context)
	if path == "" {
		path = "screenshot.png" // 默认路径
	}

	step := service.AutomationStep{
		Type: "screenshot",
		Parameters: map[string]interface{}{
			"path": path,
		},
	}

	opResult := e.automationService.ExecuteStep(step)
	if !opResult.Success {
		result.Error = fmt.Sprintf("截屏失败: %s", opResult.Error)
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("成功截屏: %s", path)
	result.ScreenshotPath = path
	return result
}

// executeClipboardStep 执行剪贴板步骤
func (e *EnhancedTaskExecutionEngine) executeClipboardStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{
		StepType: stepPlan.Type,
		Success:  false,
	}

	var operation string
	var text string

	if e.isGetClipboardOperation(stepPlan.Context) {
		operation = "get"
	} else {
		operation = "set"
		text = e.extractTextFromContext(stepPlan.Context)
		if text == "" {
			result.Error = "无法从上下文中提取文本内容"
			return result
		}
	}

	step := service.AutomationStep{
		Type: "clipboard",
		Parameters: map[string]interface{}{
			"operation": operation,
			"text":      text,
		},
	}

	opResult := e.automationService.ExecuteStep(step)
	if !opResult.Success {
		result.Error = fmt.Sprintf("剪贴板操作失败: %s", opResult.Error)
		return result
	}

	result.Success = true
	if operation == "get" {
		result.Message = "获取剪贴板内容成功"
	} else {
		result.Message = fmt.Sprintf("设置剪贴板内容: %s", text)
	}
	return result
}

// executeWaitStep 执行等待步骤
func (e *EnhancedTaskExecutionEngine) executeWaitStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{
		StepType: stepPlan.Type,
		Success:  false,
	}

	duration := e.extractDurationFromContext(stepPlan.Context)
	if duration <= 0 {
		duration = 1000 // 默认等待1秒
	}

	step := service.AutomationStep{
		Type: "wait",
		Parameters: map[string]interface{}{
			"duration": float64(duration),
		},
	}

	opResult := e.automationService.ExecuteStep(step)
	if !opResult.Success {
		result.Error = fmt.Sprintf("等待操作失败: %s", opResult.Error)
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("等待 %d 毫秒", duration)
	return result
}

// executeKeyPressStep 执行按键步骤
func (e *EnhancedTaskExecutionEngine) executeKeyPressStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{
		StepType: stepPlan.Type,
		Success:  false,
	}

	key, modifiers := e.extractKeyFromContext(stepPlan.Context)
	if key == "" {
		result.Error = "无法从上下文中提取按键信息"
		return result
	}

	step := service.AutomationStep{
		Type: "key",
		Parameters: map[string]interface{}{
			"key":       key,
			"modifiers": modifiers,
		},
	}

	opResult := e.automationService.ExecuteStep(step)
	if !opResult.Success {
		result.Error = fmt.Sprintf("按键操作失败: %s", opResult.Error)
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("成功按键: %s", key)
	return result
}
