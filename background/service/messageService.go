package service

import (
	"fmt"
	"log/slog"

	"diandian/background/app"
	"diandian/background/constant"
	"diandian/background/database"
	"diandian/background/model"

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
		slog.Error("处理用户消息失败", "error", err)
		s.sendErrorMessage("系统错误，请稍后重试")
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

// // 发送任务分析结果
// func (s *MessageService) sendTaskAnalysis(analysis *AutomationTaskResponse) {
// 	message := fmt.Sprintf(`📋 **任务分析结果**
//
// **任务名称**: %s
// **描述**: %s
// **复杂度**: %s
//
// **执行步骤**:
// `, analysis.TaskName, analysis.Description, analysis.Complexity)
//
// 	for i, step := range analysis.Steps {
// 		message += fmt.Sprintf("%d. %s\n", i+1, step)
// 	}
//
// 	if len(analysis.Risks) > 0 {
// 		message += "\n⚠️ **风险提示**:\n"
// 		for _, risk := range analysis.Risks {
// 			message += fmt.Sprintf("• %s\n", risk)
// 		}
// 	}
//
// 	// s.sendMessage(message)
// }
//
// 发送确认请求
// func (s *MessageService) sendConfirmationRequest(task *model.Task, analysis *AutomationTaskResponse) {
// 	// 更新任务状态为等待确认
// 	task.Status = "waiting_confirm"
// 	task.Progress = 80
// 	database.DB.Save(task)
// 	s.sendTaskUpdate(task)
//
// 	// 发送确认消息
// 	confirmMessage := `🔐 **需要您的确认**
//
// 此任务涉及重要操作，需要您的明确授权才能继续执行。
// 确认后，任务将自动执行，无需进一步干预。
//
// ⚠️ **重要提示**：
// • 确认后界面将切换到浮动模式
// • 任务将在后台自动执行
// • 执行过程中请勿手动操作电脑
// • 您可以随时通过浮动窗口监控进度
//
// 请仔细检查任务详情后点击确认：`
//
// 	s.sendMessage(confirmMessage)
//
// 	// 发送确认按钮事件
// 	app.EmitEvent("automation-confirm-request", map[string]interface{}{
// 		"task_id":  task.ID,
// 		"analysis": analysis,
// 	})
// }

// 执行自动化任务
func (s *MessageService) executeAutomationTask(task *model.Task, analysis *AutomationTaskResponse) {
	app.EmitEvent(constant.EventNotify, "任务开始执行...")

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
	// 发送执行状态更新
	// s.sendExecutionUpdate("🔍 正在分析屏幕...")

	// TODO: 这里将来会调用自动化任务执行引擎
	// 目前先模拟执行过程

	// 模拟执行步骤
	for i, step := range analysis.Steps {
		// s.sendExecutionUpdate(fmt.Sprintf("⚡ 执行步骤 %d/%d: %s", i+1, len(analysis.Steps), step))

		// 模拟执行时间
		// time.Sleep(2 * time.Second)
		slog.Debug("执行步骤", "step", step)

		// 更新进度
		progress := 70 + (30 * (i + 1) / len(analysis.Steps))
		task.Progress = progress
		database.DB.Save(task)
		s.sendTaskUpdate(task)
	}

	// 完成任务
	// s.sendExecutionUpdate("✅ 自动化任务执行完成")
	s.updateTaskStatus(task, model.TaskStatusCompleted, "任务执行完成")
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
		taskAnalysis, err := llmService.AnalyzeAutomationTask(task.Description)
		if err != nil {
			s.updateTaskStatus(&task, model.TaskStatusFailed, "重新分析任务失败")
			return err
		}
		s.executeAutomationTask(&task, taskAnalysis)
	} else {
		s.updateTaskStatus(&task, model.TaskStatusCancelled, "用户取消执行")
	}

	return nil
}
