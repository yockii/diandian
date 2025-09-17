package hybrid

import (
	"fmt"
	"runtime"

	"diandian/background/automation/core"
)

// HybridEngine 混合自动化引擎
// 优先使用纯Go实现，必要时回退到外部程序
type HybridEngine struct {
	pureGoEngine   AutomationEngine
	externalEngine AutomationEngine
	preferPureGo   bool
}

// AutomationEngine 自动化引擎接口
type AutomationEngine interface {
	Click(x, y int) *core.OperationResult
	Type(text string) *core.OperationResult
	KeyPress(key string) *core.OperationResult
	Screenshot() *core.OperationResult
	IsAvailable() bool
}

// NewHybridEngine 创建混合引擎
func NewHybridEngine() (*HybridEngine, error) {
	engine := &HybridEngine{
		preferPureGo: true,
	}

	// 初始化纯Go引擎
	pureGo, err := NewPureGoEngine()
	if err == nil && pureGo.IsAvailable() {
		engine.pureGoEngine = pureGo
	}

	// 初始化外部程序引擎作为回退
	external, err := NewExternalEngine()
	if err == nil && external.IsAvailable() {
		engine.externalEngine = external
	}

	// 至少需要一个可用的引擎
	if engine.pureGoEngine == nil && engine.externalEngine == nil {
		return nil, fmt.Errorf("no available automation engine")
	}

	return engine, nil
}

// getEngine 获取当前应该使用的引擎
func (h *HybridEngine) getEngine() AutomationEngine {
	if h.preferPureGo && h.pureGoEngine != nil && h.pureGoEngine.IsAvailable() {
		return h.pureGoEngine
	}
	if h.externalEngine != nil && h.externalEngine.IsAvailable() {
		return h.externalEngine
	}
	// 如果首选不可用，尝试另一个
	if h.pureGoEngine != nil && h.pureGoEngine.IsAvailable() {
		return h.pureGoEngine
	}
	return nil
}

// Click 点击操作
func (h *HybridEngine) Click(x, y int) *core.OperationResult {
	engine := h.getEngine()
	if engine == nil {
		return core.NewErrorResult("no available engine for click operation", fmt.Errorf("no engine"))
	}
	return engine.Click(x, y)
}

// Type 输入文本
func (h *HybridEngine) Type(text string) *core.OperationResult {
	engine := h.getEngine()
	if engine == nil {
		return core.NewErrorResult("no available engine for type operation", fmt.Errorf("no engine"))
	}
	return engine.Type(text)
}

// KeyPress 按键操作
func (h *HybridEngine) KeyPress(key string) *core.OperationResult {
	engine := h.getEngine()
	if engine == nil {
		return core.NewErrorResult("no available engine for keypress operation", fmt.Errorf("no engine"))
	}
	return engine.KeyPress(key)
}

// Screenshot 截屏
func (h *HybridEngine) Screenshot() *core.OperationResult {
	engine := h.getEngine()
	if engine == nil {
		return core.NewErrorResult("no available engine for screenshot operation", fmt.Errorf("no engine"))
	}
	return engine.Screenshot()
}

// SetPreferPureGo 设置是否优先使用纯Go引擎
func (h *HybridEngine) SetPreferPureGo(prefer bool) {
	h.preferPureGo = prefer
}

// GetEngineInfo 获取当前引擎信息
func (h *HybridEngine) GetEngineInfo() map[string]interface{} {
	info := map[string]interface{}{
		"platform":       runtime.GOOS,
		"prefer_pure_go": h.preferPureGo,
	}

	if h.pureGoEngine != nil {
		info["pure_go_available"] = h.pureGoEngine.IsAvailable()
	} else {
		info["pure_go_available"] = false
	}

	if h.externalEngine != nil {
		info["external_available"] = h.externalEngine.IsAvailable()
	} else {
		info["external_available"] = false
	}

	currentEngine := h.getEngine()
	if currentEngine == h.pureGoEngine {
		info["current_engine"] = "pure_go"
	} else if currentEngine == h.externalEngine {
		info["current_engine"] = "external"
	} else {
		info["current_engine"] = "none"
	}

	return info
}
