package service

import (
	"fmt"
	"log/slog"

	"diandian/background/automation/core"
	"diandian/background/domain"
)

// EnhancedVisionService 增强的视觉分析服务
type EnhancedVisionService struct {
	llmService         *LLMService
	automationService  *AutomationService
	activeDisplayIndex int // 当前活动的显示器索引，-1表示未固定
}

// NewEnhancedVisionService 创建增强视觉分析服务
func NewEnhancedVisionService(llmService *LLMService, automationService *AutomationService) *EnhancedVisionService {
	return &EnhancedVisionService{
		llmService:         llmService,
		automationService:  automationService,
		activeDisplayIndex: -1, // 初始未固定
	}
}

// MultiDisplayAnalysis 多显示器分析结果
type MultiDisplayAnalysis struct {
	Displays              []DisplayAnalysisResult       `json:"displays"`
	RecommendedDisplay    int                           `json:"recommended_display"`
	GlobalRecommendations []domain.ActionRecommendation `json:"global_recommendations"`
}

// DisplayAnalysisResult 单个显示器分析结果
type DisplayAnalysisResult struct {
	DisplayIndex      int                    `json:"display_index"`
	Width             int                    `json:"width"`
	Height            int                    `json:"height"`
	ElementsFound     []domain.VisualElement `json:"elements_found"`
	ScreenInfo        domain.ScreenInfo      `json:"screen_info"`
	Confidence        float64                `json:"confidence"`
	HasTargetElements bool                   `json:"has_target_elements"`
}

// AnalyzeAllDisplays 分析所有显示器
func (evs *EnhancedVisionService) AnalyzeAllDisplays(context string) (*MultiDisplayAnalysis, error) {
	slog.Info("开始多显示器视觉分析", "context", context)

	// 获取所有显示器截图
	captures, err := evs.captureAllDisplays()
	if err != nil {
		return nil, fmt.Errorf("获取显示器截图失败: %v", err)
	}

	if len(captures) == 0 {
		return nil, fmt.Errorf("没有可用的显示器")
	}

	// 如果只有一个显示器，直接分析
	if len(captures) == 1 {
		result, err := evs.analyzeSingleDisplay(captures[0], context)
		if err != nil {
			return nil, err
		}

		return &MultiDisplayAnalysis{
			Displays:              []DisplayAnalysisResult{*result},
			RecommendedDisplay:    0,
			GlobalRecommendations: result.toActionRecommendations(),
		}, nil
	}

	// 多显示器分析
	return evs.analyzeMultipleDisplays(captures, context)
}

// AnalyzeActiveDisplay 分析当前活动显示器
func (evs *EnhancedVisionService) AnalyzeActiveDisplay(context string) (*domain.VisualAnalysisResponse, error) {
	if evs.activeDisplayIndex >= 0 {
		// 使用固定的显示器
		return evs.AnalyzeSpecificDisplay(evs.activeDisplayIndex, context)
	}

	// 分析所有显示器并选择最佳的
	multiResult, err := evs.AnalyzeAllDisplays(context)
	if err != nil {
		return nil, err
	}

	// 设置活动显示器
	evs.activeDisplayIndex = multiResult.RecommendedDisplay

	// 返回推荐显示器的分析结果
	if multiResult.RecommendedDisplay < len(multiResult.Displays) {
		displayResult := multiResult.Displays[multiResult.RecommendedDisplay]
		return &domain.VisualAnalysisResponse{
			ElementsFound:   displayResult.ElementsFound,
			ScreenInfo:      displayResult.ScreenInfo,
			Recommendations: multiResult.GlobalRecommendations,
		}, nil
	}

	return nil, fmt.Errorf("无法确定活动显示器")
}

// AnalyzeSpecificDisplay 分析指定显示器
func (evs *EnhancedVisionService) AnalyzeSpecificDisplay(displayIndex int, context string) (*domain.VisualAnalysisResponse, error) {
	slog.Info("分析指定显示器", "display", displayIndex, "context", context)

	capture, err := evs.captureSpecificDisplay(displayIndex)
	if err != nil {
		return nil, fmt.Errorf("获取显示器 %d 截图失败: %v", displayIndex, err)
	}

	result, err := evs.analyzeSingleDisplay(*capture, context)
	if err != nil {
		return nil, err
	}

	return &domain.VisualAnalysisResponse{
		ElementsFound:   result.ElementsFound,
		ScreenInfo:      result.ScreenInfo,
		Recommendations: result.toActionRecommendations(),
	}, nil
}

// SetActiveDisplay 设置活动显示器
func (evs *EnhancedVisionService) SetActiveDisplay(displayIndex int) {
	evs.activeDisplayIndex = displayIndex
	slog.Info("设置活动显示器", "display", displayIndex)
}

// ResetActiveDisplay 重置活动显示器（下次分析时重新选择）
func (evs *EnhancedVisionService) ResetActiveDisplay() {
	evs.activeDisplayIndex = -1
	slog.Info("重置活动显示器")
}

// GetActiveDisplayIndex 获取当前活动显示器索引
func (evs *EnhancedVisionService) GetActiveDisplayIndex() int {
	return evs.activeDisplayIndex
}

// captureAllDisplays 获取所有显示器截图
func (evs *EnhancedVisionService) captureAllDisplays() ([]core.DisplayCapture, error) {
	if evs.automationService == nil {
		return nil, fmt.Errorf("automation service not initialized")
	}

	// 通过automation引擎获取真实的截图
	engine := evs.automationService.GetEngine()
	if engine == nil {
		return nil, fmt.Errorf("automation engine not available")
	}

	// 使用标准的Screenshot接口
	result := engine.Screenshot()
	if !result.Success {
		return nil, fmt.Errorf("failed to capture screen: %s", result.Error)
	}

	// 获取屏幕尺寸（使用默认值）
	width, height := 1920, 1080 // 默认值

	// 从result中获取图像数据
	var imageData []byte
	if data, ok := result.Data.(map[string]interface{}); ok {
		if dataStr, exists := data["data"]; exists {
			if dataBytes, ok := dataStr.([]byte); ok {
				imageData = dataBytes
			}
		}
	}

	// 创建单个显示器的capture（作为多显示器的fallback）
	capture := core.DisplayCapture{
		Index:     0,
		ImageData: imageData,
		Width:     width,
		Height:    height,
		IsActive:  true,
	}

	captures := []core.DisplayCapture{capture}
	slog.Info("成功获取显示器截图", "count", len(captures), "size", fmt.Sprintf("%dx%d", width, height))
	return captures, nil
}

// captureSpecificDisplay 获取指定显示器截图
func (evs *EnhancedVisionService) captureSpecificDisplay(displayIndex int) (*core.DisplayCapture, error) {
	// 对于单显示器系统，忽略displayIndex
	if displayIndex != 0 {
		slog.Warn("请求的显示器索引超出范围，使用主显示器", "requested", displayIndex)
	}

	// 获取所有显示器（实际上是单显示器）
	captures, err := evs.captureAllDisplays()
	if err != nil {
		return nil, err
	}

	if len(captures) == 0 {
		return nil, fmt.Errorf("no displays available")
	}

	// 返回第一个（也是唯一的）显示器
	return &captures[0], nil
}

// analyzeSingleDisplay 分析单个显示器
func (evs *EnhancedVisionService) analyzeSingleDisplay(capture core.DisplayCapture, context string) (*DisplayAnalysisResult, error) {
	// 构建分析上下文
	analysisContext := fmt.Sprintf("显示器 %d (%dx%d): %s",
		capture.Index, capture.Width, capture.Height, context)

	// 调用LLM进行视觉分析（直接使用图像数据）
	response, err := evs.llmService.AnalyzeScreenshot(capture.ImageData, analysisContext)
	if err != nil {
		return nil, fmt.Errorf("LLM视觉分析失败: %v", err)
	}

	// 计算置信度（基于找到的元素数量和质量）
	confidence := evs.calculateConfidence(response)

	// 检查是否有目标元素
	hasTargetElements := len(response.ElementsFound) > 0

	result := &DisplayAnalysisResult{
		DisplayIndex:      capture.Index,
		Width:             capture.Width,
		Height:            capture.Height,
		ElementsFound:     response.ElementsFound,
		ScreenInfo:        response.ScreenInfo,
		Confidence:        confidence,
		HasTargetElements: hasTargetElements,
	}

	return result, nil
}

// analyzeMultipleDisplays 分析多个显示器
func (evs *EnhancedVisionService) analyzeMultipleDisplays(captures []core.DisplayCapture, context string) (*MultiDisplayAnalysis, error) {
	var displayResults []DisplayAnalysisResult
	var bestDisplayIndex int
	var bestConfidence float64

	// 分析每个显示器
	for _, capture := range captures {
		result, err := evs.analyzeSingleDisplay(capture, context)
		if err != nil {
			slog.Error("显示器分析失败", "display", capture.Index, "error", err)
			continue
		}

		displayResults = append(displayResults, *result)

		// 选择最佳显示器（置信度最高且有目标元素）
		if result.HasTargetElements && result.Confidence > bestConfidence {
			bestConfidence = result.Confidence
			bestDisplayIndex = len(displayResults) - 1
		}
	}

	if len(displayResults) == 0 {
		return nil, fmt.Errorf("所有显示器分析都失败")
	}

	// 如果没有找到有目标元素的显示器，选择置信度最高的
	if bestConfidence == 0 {
		for i, result := range displayResults {
			if result.Confidence > bestConfidence {
				bestConfidence = result.Confidence
				bestDisplayIndex = i
			}
		}
	}

	// 生成全局建议
	globalRecommendations := evs.generateGlobalRecommendations(displayResults, bestDisplayIndex)

	return &MultiDisplayAnalysis{
		Displays:              displayResults,
		RecommendedDisplay:    bestDisplayIndex,
		GlobalRecommendations: globalRecommendations,
	}, nil
}

// calculateConfidence 计算分析置信度
func (evs *EnhancedVisionService) calculateConfidence(response *domain.VisualAnalysisResponse) float64 {
	if len(response.ElementsFound) == 0 {
		return 0.0
	}

	totalConfidence := 0.0
	for _, element := range response.ElementsFound {
		totalConfidence += element.Confidence
	}

	// 平均置信度，并考虑元素数量的影响
	avgConfidence := totalConfidence / float64(len(response.ElementsFound))
	elementCountFactor := float64(len(response.ElementsFound)) / 10.0 // 最多10个元素为满分
	if elementCountFactor > 1.0 {
		elementCountFactor = 1.0
	}

	return avgConfidence * (0.7 + 0.3*elementCountFactor)
}

// generateGlobalRecommendations 生成全局建议
func (evs *EnhancedVisionService) generateGlobalRecommendations(results []DisplayAnalysisResult, bestIndex int) []domain.ActionRecommendation {
	var recommendations []domain.ActionRecommendation

	if bestIndex < len(results) {
		bestResult := results[bestIndex]

		// 添加显示器切换建议
		recommendations = append(recommendations, domain.ActionRecommendation{
			Action: "switch_display",
			Target: fmt.Sprintf("显示器 %d", bestResult.DisplayIndex),
			Reason: fmt.Sprintf("该显示器包含最相关的内容 (置信度: %.2f)", bestResult.Confidence),
		})

		// 添加最佳显示器的元素操作建议
		recommendations = append(recommendations, bestResult.toActionRecommendations()...)
	}

	return recommendations
}

// toActionRecommendations 将显示器分析结果转换为操作建议
func (dar *DisplayAnalysisResult) toActionRecommendations() []domain.ActionRecommendation {
	var recommendations []domain.ActionRecommendation

	for _, element := range dar.ElementsFound {
		if element.Clickable {
			recommendations = append(recommendations, domain.ActionRecommendation{
				Action: "click",
				Target: element.Description,
				Reason: fmt.Sprintf("可点击元素，置信度: %.2f", element.Confidence),
			})
		}
	}

	return recommendations
}
