package operation

import (
	"fmt"
	"log/slog"
	"strings"

	"diandian/background/database"
	"diandian/background/model"

	"github.com/sashabaranov/go-openai"
)

// BaseGenerator 操作生成器基础结构
type BaseGenerator struct {
	textClient   *openai.Client
	textModel    string
	visionClient *openai.Client
	visionModel  string
}

// NewBaseGenerator 创建基础生成器
func NewBaseGenerator() *BaseGenerator {
	return &BaseGenerator{}
}

// createTextClient 创建文本模型客户端
func (g *BaseGenerator) createTextClient() (*openai.Client, string, error) {
	if g.textClient != nil {
		return g.textClient, g.textModel, nil
	}

	config, err := g.getTextModelConfig()
	if err != nil {
		return nil, "", err
	}

	clientConfig := openai.DefaultConfig(config.Token)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	g.textClient = openai.NewClientWithConfig(clientConfig)
	g.textModel = config.Model
	return g.textClient, g.textModel, nil
}

// createVisionClient 创建视觉模型客户端
func (g *BaseGenerator) createVisionClient() (*openai.Client, string, error) {
	if g.visionClient != nil {
		return g.visionClient, g.visionModel, nil
	}

	config, err := g.getVisionModelConfig()
	if err != nil {
		return nil, "", err
	}

	clientConfig := openai.DefaultConfig(config.Token)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	g.visionClient = openai.NewClientWithConfig(clientConfig)
	g.visionModel = config.Model
	return g.visionClient, g.visionModel, nil
}

// TextModelConfig 文本模型配置
type TextModelConfig struct {
	Model   string
	Token   string
	BaseURL string
}

// VisionModelConfig 视觉模型配置
type VisionModelConfig struct {
	Model   string
	Token   string
	BaseURL string
}

// getTextModelConfig 获取文本模型配置
func (g *BaseGenerator) getTextModelConfig() (*TextModelConfig, error) {
	var settings []*model.Setting
	err := database.DB.Where("key IN ?", []string{
		model.SettingKeyLlmTextModel,
		model.SettingKeyLlmTextToken,
		model.SettingKeyLlmTextBaseUrl,
	}).Find(&settings).Error
	if err != nil {
		slog.Error("获取文本模型配置失败", "error", err)
		return nil, fmt.Errorf("获取文本模型配置失败: %v", err)
	}

	config := &TextModelConfig{}
	for _, setting := range settings {
		if setting.Value == nil {
			continue
		}
		switch setting.Key {
		case model.SettingKeyLlmTextModel:
			config.Model = *setting.Value
		case model.SettingKeyLlmTextToken:
			config.Token = *setting.Value
		case model.SettingKeyLlmTextBaseUrl:
			config.BaseURL = *setting.Value
		}
	}

	if config.Model == "" || config.Token == "" {
		return nil, fmt.Errorf("文本模型配置不完整")
	}

	return config, nil
}

// getVisionModelConfig 获取视觉模型配置
func (g *BaseGenerator) getVisionModelConfig() (*VisionModelConfig, error) {
	var settings []*model.Setting
	err := database.DB.Where("key IN ?", []string{
		model.SettingKeyLlmVlModel,
		model.SettingKeyLlmVlToken,
		model.SettingKeyLlmVlBaseUrl,
	}).Find(&settings).Error
	if err != nil {
		slog.Error("获取视觉模型配置失败", "error", err)
		return nil, fmt.Errorf("获取视觉模型配置失败: %v", err)
	}

	config := &VisionModelConfig{}
	for _, setting := range settings {
		if setting.Value == nil {
			continue
		}
		switch setting.Key {
		case model.SettingKeyLlmVlModel:
			config.Model = *setting.Value
		case model.SettingKeyLlmVlToken:
			config.Token = *setting.Value
		case model.SettingKeyLlmVlBaseUrl:
			config.BaseURL = *setting.Value
		}
	}

	if config.Model == "" || config.Token == "" {
		return nil, fmt.Errorf("视觉模型配置不完整")
	}

	return config, nil
}

// retryLLMCall 重试LLM调用的通用方法
func (g *BaseGenerator) retryLLMCall(
	callFunc func() (string, error),
	validateFunc func(content string) error,
	maxRetries int,
	operation string,
) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// 调用LLM
		content, err := callFunc()
		if err != nil {
			lastErr = fmt.Errorf("LLM调用失败: %v", err)
			if attempt < maxRetries {
				fmt.Printf("🔄 %s第%d次尝试失败，重试中... 错误: %v\n", operation, attempt, err)
				continue
			}
			break
		}

		// 清理内容
		cleanedContent := cleanMarkdownCodeBlock(content)

		// 验证内容
		if validateFunc != nil {
			if err := validateFunc(cleanedContent); err != nil {
				lastErr = fmt.Errorf("内容验证失败: %v", err)
				if attempt < maxRetries {
					fmt.Printf("🔄 %s第%d次尝试验证失败，重试中... 错误: %v\n", operation, attempt, err)
					continue
				}
				break
			}
		}

		// 成功
		if attempt > 1 {
			fmt.Printf("✅ %s在第%d次尝试后成功\n", operation, attempt)
		}
		return cleanedContent, nil
	}

	return "", fmt.Errorf("%s在%d次尝试后仍然失败，最后错误: %v", operation, maxRetries, lastErr)
}

// cleanMarkdownCodeBlock 清理markdown代码块标记
func cleanMarkdownCodeBlock(content string) string {
	content = strings.TrimSpace(content)

	// 移除开头的```json或```
	if strings.HasPrefix(content, "```json") {
		content = content[7:]
	} else if strings.HasPrefix(content, "```") {
		content = content[3:]
	}

	// 移除结尾的```
	if strings.HasSuffix(content, "```") {
		content = content[:len(content)-3]
	}

	return strings.TrimSpace(content)
}
