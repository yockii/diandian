package executor

import (
	"fmt"
	"log/slog"

	"diandian/background/domain"
	"diandian/background/service"
)

// executeFileStep 执行文件操作步骤
func (e *EnhancedTaskExecutionEngine) executeFileStep(stepPlan *domain.AutomationStepPlan) *StepExecutionResult {
	result := &StepExecutionResult{Success: false}

	// 使用LLM生成文件操作
	fileOp, err := e.llmService.GenerateFileOperation(stepPlan.Context)
	if err != nil {
		result.Error = fmt.Sprintf("生成文件操作失败: %v", err)
		return result
	}

	slog.Info("生成文件操作",
		"operation", fileOp.Operation,
		"source_path", fileOp.SourcePath,
		"target_path", fileOp.TargetPath)

	// 通过AutomationService执行文件操作
	step := service.AutomationStep{
		Type: "file",
		Parameters: map[string]interface{}{
			"operation":   fileOp.Operation,
			"source_path": fileOp.SourcePath,
			"target_path": fileOp.TargetPath,
			"content":     fileOp.Content,
		},
	}

	opResult := e.automationService.ExecuteStep(step)
	if !opResult.Success {
		result.Error = fmt.Sprintf("执行文件操作失败: %s", opResult.Error)
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("成功执行文件操作: %s", fileOp.Operation)
	return result
}
