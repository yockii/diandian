package executor

import (
	"context"
	"fmt"
	"time"

	"diandian/background/domain"
	"diandian/background/service"
)

// 使用domain包中的类型
type TaskExecutionResult = domain.TaskExecutionResult
type StepExecutionResult = domain.StepExecutionResult

// EnhancedTaskExecutionEngine 增强的任务执行引擎
type EnhancedTaskExecutionEngine struct {
	automationService *service.AutomationService
	llmService        *service.LLMService
	visionService     *service.EnhancedVisionService
}

// NewEnhancedTaskExecutionEngine 创建增强的任务执行引擎
func NewEnhancedTaskExecutionEngine(automationService *service.AutomationService) *EnhancedTaskExecutionEngine {
	llmService := &service.LLMService{}
	return &EnhancedTaskExecutionEngine{
		automationService: automationService,
		llmService:        llmService,
		visionService:     service.NewEnhancedVisionService(llmService, automationService),
	}
}

// ExecuteTaskDecomposition 执行任务分解结果
func (e *EnhancedTaskExecutionEngine) ExecuteTaskDecomposition(ctx context.Context, taskID uint, decomposition *domain.AutomationTaskDecomposition) *TaskExecutionResult {
	startTime := time.Now()

	result := &TaskExecutionResult{
		TaskID:    taskID,
		Success:   false,
		Steps:     make([]*StepExecutionResult, 0),
		StartTime: startTime,
	}

	// 执行每个步骤
	for i, stepPlan := range decomposition.Steps {
		stepResult := e.executeStepPlan(ctx, &stepPlan)
		stepResult.StepIndex = i + 1
		result.Steps = append(result.Steps, stepResult)

		// 如果步骤失败且不是可选步骤，停止执行
		if !stepResult.Success {
			if stepPlan.Optional {
				// 可选步骤失败，继续执行
				continue
			} else {
				// 必需步骤失败，停止执行
				result.Message = fmt.Sprintf("步骤 %d 执行失败: %s", i+1, stepResult.Error)
				break
			}
		}
	}

	// 计算执行结果
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	successCount := 0
	for _, step := range result.Steps {
		if step.Success {
			successCount++
		} else {
			result.ErrorCount++
		}
	}

	if len(result.Steps) > 0 {
		result.SuccessRate = float64(successCount) / float64(len(result.Steps))
		result.Success = result.ErrorCount == 0 || (result.SuccessRate >= 0.8 && result.ErrorCount <= 2)
	}

	if result.Success {
		result.Message = fmt.Sprintf("任务执行成功，成功率: %.1f%%", result.SuccessRate*100)
	} else if result.Message == "" {
		result.Message = fmt.Sprintf("任务执行失败，成功率: %.1f%%", result.SuccessRate*100)
	}

	return result
}
