package executor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"diandian/background/domain"
	"diandian/background/service"
)

// executeStepPlan 执行单个步骤计划
func (e *EnhancedTaskExecutionEngine) executeStepPlan(ctx context.Context, stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{
		StepType:  stepPlan.Type,
		Success:   false,
		StartTime: time.Now(),
	}

	slog.Info("开始执行步骤",
		"type", stepPlan.Type,
		"description", stepPlan.Description,
		"requires_screen_analysis", stepPlan.RequiresScreenAnalysis)

	// 如果需要屏幕分析，先进行截图和分析
	var screenAnalysis *domain.VisualAnalysisResponse
	if stepPlan.RequiresScreenAnalysis {
		// 通过AutomationService执行截图
		step := service.AutomationStep{
			Type: "screenshot",
			Parameters: map[string]interface{}{
				"path": "temp_screenshot.png",
			},
		}

		opResult := e.automationService.ExecuteStep(step)
		if !opResult.Success {
			result.Error = fmt.Sprintf("截图失败: %s", opResult.Error)
			result.EndTime = time.Now()
			result.Duration = result.EndTime.Sub(result.StartTime)
			return result
		}

		// 保存截图路径
		result.ScreenshotPath = "temp_screenshot.png"

		// 进行视觉分析 - 这里需要读取截图文件
		// 暂时跳过视觉分析，因为需要文件读取逻辑
		slog.Info("跳过视觉分析，直接执行步骤")
	}

	// 根据步骤类型执行相应操作
	switch stepPlan.Type {
	case "click":
		result = e.executeClickStep(stepPlan, screenAnalysis)
	case "type":
		result = e.executeTypeStep(stepPlan)
	case "launch_app":
		result = e.executeLaunchAppStep(stepPlan)
	case "file":
		result = e.executeFileStep(stepPlan)
	case "screenshot":
		result = e.executeScreenshotStep(stepPlan)
	case "clipboard":
		result = e.executeClipboardStep(stepPlan)
	case "wait":
		result = e.executeWaitStep(stepPlan)
	case "key_press":
		result = e.executeKeyPressStep(stepPlan)
	default:
		result.Error = fmt.Sprintf("不支持的步骤类型: %s", stepPlan.Type)
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	if result.Success {
		slog.Info("步骤执行成功",
			"type", stepPlan.Type,
			"duration", result.Duration,
			"message", result.Message)
	} else {
		slog.Error("步骤执行失败",
			"type", stepPlan.Type,
			"duration", result.Duration,
			"error", result.Error)
	}

	return result
}
