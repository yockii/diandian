package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"diandian/background/app"
	"diandian/background/constant"
	"diandian/background/database"
	"diandian/background/domain"
	"diandian/background/model"

	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type MessageService struct{}

// 处理新消息
func (s *MessageService) NewMessage(msg *model.Message) {
	go s.processMessageAsync(msg)
}

// 异步处理消息
func (s *MessageService) processMessageAsync(msg *model.Message) {
	// 会话处理
	msg.Role = model.MessageRoleUser
	// 保存消息
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if msg.ConversationID == 0 {
			conversation := &model.Conversation{}
			if err := tx.Create(conversation).Error; err != nil {
				return err
			}
			msg.ConversationID = conversation.ID
		}
		return tx.Create(msg).Error
	})
	if err != nil {
		slog.Error("保存用户消息失败", "error", err)
		s.sendErrorMessage("保存消息失败，请稍后重试")
		return
	}

	assistantMsg, response, err := DefaultLLMService.ProcessMessage(msg.ConversationID)

	if err != nil {
		slog.Error("处理用户消息失败", "error", err, "conversation_id", msg.ConversationID, "message_content", msg.Content)

		// 根据错误类型提供更具体的错误信息
		var errorMsg string
		if strings.Contains(err.Error(), "API") {
			errorMsg = "AI服务暂时不可用，请稍后重试"
		} else if strings.Contains(err.Error(), "解析") {
			errorMsg = "AI响应格式异常，请重新发送消息"
		} else if strings.Contains(err.Error(), "网络") || strings.Contains(err.Error(), "timeout") {
			errorMsg = "网络连接异常，请检查网络后重试"
		} else {
			errorMsg = "系统错误，请稍后重试"
		}

		s.sendErrorMessage(errorMsg)
		return
	}

	// 更新会话标题
	if response.ConversationTitle != "" {
		database.DB.Model(&model.Conversation{}).Where("id = ?", msg.ConversationID).Update("name", response.ConversationTitle)
	}

	if response.MessageType == "automation" {
		// 先发送聊天回复
		s.sendMessage(assistantMsg)
		// 处理自动化任务
		s.handleAutomationTask(response, msg.ConversationID)
	} else {
		s.sendMessage(assistantMsg)
	}
}

// 处理自动化任务
func (s *MessageService) handleAutomationTask(response *UnifiedMessageResponse, conversationID uint64) {
	if response.AutomationTask == nil {
		return
	}

	task := new(model.Task)

	taskAnalysis := response.AutomationTask

	// 更新任务信息
	task.ConversationID = conversationID
	task.Progress = 50
	task.Name = taskAnalysis.TaskName
	task.Description = taskAnalysis.Description
	task.Status = model.TaskStatusPending
	database.DB.Save(task)
	s.sendTaskUpdate(task)

	// 如果需要用户确认，发送确认请求
	if !taskAnalysis.NeedsConfirm {
		// 直接执行简单任务
		s.executeAutomationTask(task, taskAnalysis)
	}
}

// 更新任务状态，并发送更新通知
func (s *MessageService) updateTaskStatus(task *model.Task, status, result string) {
	task.Status = status
	task.Result = result

	switch status {
	case model.TaskStatusCompleted:
		task.Progress = 100
	case model.TaskStatusFailed:
		task.ErrorMsg = result
	}

	database.DB.Save(task)
	s.sendTaskUpdate(task)
}

// 发送任务更新通知
func (s *MessageService) sendTaskUpdate(task *model.Task) {
	app.EmitEvent(constant.EventTaskStatusChanged, task)
}

// 发送消息
func (s *MessageService) sendMessage(msg *model.Message) {
	app.EmitEvent(constant.EventMessageResponsed, msg)
}

// 发送错误消息
func (s *MessageService) sendErrorMessage(content string) {
	app.EmitEvent(constant.EventOperateFailed, &model.Step{
		StepType: model.StepTypeMessage,
		Status:   model.StepStatusFailed,
		Content:  "❌ " + content,
	})
}

// 执行自动化任务
func (s *MessageService) executeAutomationTask(task *model.Task, analysis *AutomationTaskResponse) {
	app.EmitEvent(constant.EventNotify, "任务开始执行...")

	// 发送任务执行开始事件，触发窗口切换
	app.EmitEvent(constant.EventTaskExecutionStarted, task)

	// 更新任务状态
	task.Status = model.TaskStatusRunning
	task.Progress = 70
	database.DB.Save(task)
	s.sendTaskUpdate(task)

	// 启动自动化任务执行（异步）
	go s.runAutomationTask(task, analysis)
}

// 运行自动化任务（后台执行）
func (s *MessageService) runAutomationTask(task *model.Task, analysis *AutomationTaskResponse) {
	// 模拟执行过程
	for i, step := range analysis.Steps {
		slog.Debug("执行步骤", "step", step)

		// 更新进度
		progress := 70 + (30 * (i + 1) / len(analysis.Steps))
		task.Progress = progress
		database.DB.Save(task)
		s.sendTaskUpdate(task)
	}

	// 完成任务
	s.updateTaskStatus(task, model.TaskStatusCompleted, "任务执行完成")

	// 发送任务执行完成事件，触发窗口恢复
	app.EmitEvent(constant.EventTaskExecutionCompleted, task)
}

// 执行新的自动化任务（使用增强的执行引擎）
func (s *MessageService) executeAutomationTaskEnhanced(task *model.Task, decomposition *domain.AutomationTaskDecomposition) {
	app.EmitEvent(constant.EventNotify, "增强任务开始执行...")

	// 发送任务执行开始事件，触发窗口切换
	app.EmitEvent(constant.EventTaskExecutionStarted, task)

	// 更新任务状态
	task.Status = model.TaskStatusRunning
	task.Progress = 70
	database.DB.Save(task)
	s.sendTaskUpdate(task)

	// 启动增强任务执行（异步）
	go s.runAutomationTaskEnhanced(task, decomposition)
}

// 运行增强的自动化任务（后台执行）
func (s *MessageService) runAutomationTaskEnhanced(task *model.Task, decomposition *domain.AutomationTaskDecomposition) {
	// 创建自动化服务
	automationService := NewAutomationService(app.GetApp())
	err := automationService.Initialize()
	if err != nil {
		slog.Error("初始化自动化服务失败", "error", err)
		s.updateTaskStatus(task, model.TaskStatusFailed, "初始化自动化服务失败")
		return
	}
	defer automationService.Cleanup()

	// 创建增强的任务执行引擎
	enhancedEngine := NewEnhancedTaskExecutionEngine(automationService)

	// 执行任务
	ctx := context.Background()
	result := enhancedEngine.ExecuteTaskDecomposition(ctx, uint(task.ID), decomposition)

	if result.Success {
		s.updateTaskStatus(task, model.TaskStatusCompleted, "增强任务执行完成")
		app.EmitEvent(constant.EventNotify, "✅ 增强自动化任务执行完成")
	} else {
		s.updateTaskStatus(task, model.TaskStatusFailed, fmt.Sprintf("增强任务执行失败: %s", result.Error))
		app.EmitEvent(constant.EventNotify, fmt.Sprintf("❌ 增强任务执行失败: %s", result.Error))
	}

	// 发送任务执行完成事件，触发窗口恢复
	app.EmitEvent(constant.EventTaskExecutionCompleted, task)
}

// 发送执行状态更新（发送到浮动窗口）
// func (s *MessageService) sendExecutionUpdate(message string) {
// 	app.EmitEvent("automation-execution-update", map[string]interface{}{
// 		"message":   message,
// 		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
// 	})
// }

// 确认执行自动化任务
func (s *MessageService) ConfirmAutomationTask(t *model.Task, confirmed bool) error {
	// 查找任务
	var task model.Task
	if err := database.DB.First(&task, t.ID).Error; err != nil {
		return fmt.Errorf("任务不存在")
	}

	if confirmed {
		// 重新分析任务（从数据库获取原始内容）
		llmService := &LLMService{}

		// 构建对话历史
		var conversationHistory []openai.ChatCompletionMessage
		// TODO: 从数据库获取对话历史

		conversationHistory = append(conversationHistory, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: task.Description,
		})

		taskDecomposition, err := llmService.DecomposeAutomationTask(conversationHistory)
		if err != nil {
			s.updateTaskStatus(&task, model.TaskStatusFailed, "重新分析任务失败")
			return err
		}
		s.executeAutomationTaskEnhanced(&task, taskDecomposition)
	} else {
		s.updateTaskStatus(&task, model.TaskStatusCancelled, "用户取消执行")
	}

	return nil
}
