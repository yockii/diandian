package hybrid

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"diandian/background/automation/core"
)

// ExternalEngine 外部程序引擎（使用robotgo的独立程序）
type ExternalEngine struct {
	workerPath string
	available  bool
}

// AutomationRequest 自动化请求
type AutomationRequest struct {
	Action     string                 `json:"action"`
	Parameters map[string]interface{} `json:"parameters"`
}

// AutomationResponse 自动化响应
type AutomationResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// NewExternalEngine 创建外部程序引擎
func NewExternalEngine() (*ExternalEngine, error) {
	engine := &ExternalEngine{}
	
	// 查找worker程序
	workerPath, err := engine.findWorkerPath()
	if err != nil {
		return engine, err // 返回引擎但标记为不可用
	}

	engine.workerPath = workerPath
	engine.available = true
	return engine, nil
}

// findWorkerPath 查找worker程序路径
func (e *ExternalEngine) findWorkerPath() (string, error) {
	// 可能的worker程序位置
	candidates := []string{
		"./automation-worker",
		"./automation-worker.exe",
		"./bin/automation-worker",
		"./bin/automation-worker.exe",
		"../automation-worker/automation-worker",
		"../automation-worker/automation-worker.exe",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			abs, err := filepath.Abs(candidate)
			if err == nil {
				return abs, nil
			}
		}
	}

	return "", fmt.Errorf("automation worker not found")
}

// IsAvailable 检查引擎是否可用
func (e *ExternalEngine) IsAvailable() bool {
	return e.available && e.workerPath != ""
}

// executeCommand 执行外部命令
func (e *ExternalEngine) executeCommand(request AutomationRequest) *core.OperationResult {
	start := time.Now()

	if !e.IsAvailable() {
		result := core.NewErrorResult("external engine not available", fmt.Errorf("worker not found"))
		result.SetDuration(start)
		return result
	}

	// 序列化请求
	requestData, err := json.Marshal(request)
	if err != nil {
		result := core.NewErrorResult("failed to marshal request", err)
		result.SetDuration(start)
		return result
	}

	// 执行外部程序
	cmd := exec.Command(e.workerPath)
	cmd.Stdin = nil
	
	// 通过命令行参数传递请求
	cmd.Args = append(cmd.Args, string(requestData))

	output, err := cmd.Output()
	if err != nil {
		result := core.NewErrorResult(
			fmt.Sprintf("external worker failed: %v", err),
			err,
		)
		result.SetDuration(start)
		return result
	}

	// 解析响应
	var response AutomationResponse
	if err := json.Unmarshal(output, &response); err != nil {
		result := core.NewErrorResult("failed to parse worker response", err)
		result.SetDuration(start)
		return result
	}

	// 转换为OperationResult
	if response.Success {
		result := core.NewSuccessResult(response.Message, response.Data)
		result.SetDuration(start)
		return result
	} else {
		result := core.NewErrorResult(response.Message, fmt.Errorf(response.Error))
		result.SetDuration(start)
		return result
	}
}

// Click 点击操作
func (e *ExternalEngine) Click(x, y int) *core.OperationResult {
	request := AutomationRequest{
		Action: "click",
		Parameters: map[string]interface{}{
			"x": x,
			"y": y,
		},
	}
	return e.executeCommand(request)
}

// Type 输入文本
func (e *ExternalEngine) Type(text string) *core.OperationResult {
	request := AutomationRequest{
		Action: "type",
		Parameters: map[string]interface{}{
			"text": text,
		},
	}
	return e.executeCommand(request)
}

// KeyPress 按键操作
func (e *ExternalEngine) KeyPress(key string) *core.OperationResult {
	request := AutomationRequest{
		Action: "keypress",
		Parameters: map[string]interface{}{
			"key": key,
		},
	}
	return e.executeCommand(request)
}

// Screenshot 截屏
func (e *ExternalEngine) Screenshot() *core.OperationResult {
	request := AutomationRequest{
		Action:     "screenshot",
		Parameters: map[string]interface{}{},
	}
	return e.executeCommand(request)
}

// GetWorkerInfo 获取worker信息
func (e *ExternalEngine) GetWorkerInfo() map[string]interface{} {
	info := map[string]interface{}{
		"available":   e.available,
		"worker_path": e.workerPath,
	}

	if e.IsAvailable() {
		// 尝试获取worker版本信息
		request := AutomationRequest{
			Action:     "info",
			Parameters: map[string]interface{}{},
		}
		
		result := e.executeCommand(request)
		if result.Success {
			if data, ok := result.Data.(map[string]interface{}); ok {
				info["worker_info"] = data
			}
		}
	}

	return info
}
