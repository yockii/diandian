package executor

import (
	"fmt"
	"log/slog"

	"diandian/background/domain"
	"diandian/background/service"
)

// executeTypeStep 执行输入步骤
func (e *EnhancedTaskExecutionEngine) executeTypeStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{
		StepType: stepPlan.Type,
		Success:  false,
	}

	// 使用LLM生成输入操作
	typeOp, err := e.llmService.GenerateTypeOperation(stepPlan.Context)
	if err != nil {
		result.Error = fmt.Sprintf("生成输入操作失败: %v", err)
		return result
	}

	slog.Info("生成输入操作", "text", typeOp.Text)

	// 通过AutomationService执行输入操作
	step := service.AutomationStep{
		Type: "type",
		Parameters: map[string]interface{}{
			"text": typeOp.Text,
		},
	}

	opResult := e.automationService.ExecuteStep(step)
	if !opResult.Success {
		result.Error = fmt.Sprintf("执行输入失败: %s", opResult.Error)
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("成功输入文本: %s", typeOp.Text)
	return result
}
