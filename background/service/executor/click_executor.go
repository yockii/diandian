package executor

import (
	"fmt"
	"log/slog"

	"diandian/background/domain"
	"diandian/background/service"
)

// executeClickStep 执行点击步骤
func (e *EnhancedTaskExecutionEngine) executeClickStep(stepPlan *domain.AutomationStepPlan, screenAnalysis *domain.VisualAnalysisResponse) *StepExecutionResult {
	result := &StepExecutionResult{
		StepType: stepPlan.Type,
		Success:  false,
	}

	// 使用LLM生成点击操作
	clickOp, err := e.llmService.GenerateClickOperation(stepPlan.Context, screenAnalysis)
	if err != nil {
		result.Error = fmt.Sprintf("生成点击操作失败: %v", err)
		return result
	}

	slog.Info("生成点击操作",
		"x", clickOp.X,
		"y", clickOp.Y,
		"button", clickOp.Button)

	// 通过AutomationService执行点击操作
	step := service.AutomationStep{
		Type: "click",
		Parameters: map[string]interface{}{
			"x":      float64(clickOp.X),
			"y":      float64(clickOp.Y),
			"button": clickOp.Button,
		},
	}

	opResult := e.automationService.ExecuteStep(step)
	if !opResult.Success {
		result.Error = fmt.Sprintf("执行点击失败: %s", opResult.Error)
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("成功点击坐标 (%d, %d)", clickOp.X, clickOp.Y)
	return result
}
