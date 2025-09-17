package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-vgo/robotgo"
)

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

func main() {
	if len(os.Args) < 2 {
		response := AutomationResponse{
			Success: false,
			Error:   "no request provided",
		}
		output, _ := json.Marshal(response)
		fmt.Print(string(output))
		return
	}

	// 解析请求
	var request AutomationRequest
	err := json.Unmarshal([]byte(os.Args[1]), &request)
	if err != nil {
		response := AutomationResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to parse request: %v", err),
		}
		output, _ := json.Marshal(response)
		fmt.Print(string(output))
		return
	}

	// 执行操作
	response := executeAction(request)
	output, _ := json.Marshal(response)
	fmt.Print(string(output))
}

func executeAction(request AutomationRequest) AutomationResponse {
	switch request.Action {
	case "click":
		return handleClick(request.Parameters)
	case "type":
		return handleType(request.Parameters)
	case "keypress":
		return handleKeyPress(request.Parameters)
	case "screenshot":
		return handleScreenshot(request.Parameters)
	case "info":
		return handleInfo(request.Parameters)
	default:
		return AutomationResponse{
			Success: false,
			Error:   fmt.Sprintf("unknown action: %s", request.Action),
		}
	}
}

func handleClick(params map[string]interface{}) AutomationResponse {
	x, ok1 := params["x"].(float64)
	y, ok2 := params["y"].(float64)
	
	if !ok1 || !ok2 {
		return AutomationResponse{
			Success: false,
			Error:   "invalid x or y coordinates",
		}
	}

	robotgo.Click(int(x), int(y))
	
	return AutomationResponse{
		Success: true,
		Message: fmt.Sprintf("clicked at (%d, %d)", int(x), int(y)),
		Data: map[string]interface{}{
			"x": int(x),
			"y": int(y),
		},
	}
}

func handleType(params map[string]interface{}) AutomationResponse {
	text, ok := params["text"].(string)
	if !ok {
		return AutomationResponse{
			Success: false,
			Error:   "invalid text parameter",
		}
	}

	robotgo.TypeStr(text)
	
	return AutomationResponse{
		Success: true,
		Message: fmt.Sprintf("typed text: %s", text),
		Data: map[string]interface{}{
			"text": text,
		},
	}
}

func handleKeyPress(params map[string]interface{}) AutomationResponse {
	key, ok := params["key"].(string)
	if !ok {
		return AutomationResponse{
			Success: false,
			Error:   "invalid key parameter",
		}
	}

	robotgo.KeyTap(key)
	
	return AutomationResponse{
		Success: true,
		Message: fmt.Sprintf("pressed key: %s", key),
		Data: map[string]interface{}{
			"key": key,
		},
	}
}

func handleScreenshot(params map[string]interface{}) AutomationResponse {
	bitmap := robotgo.CaptureScreen()
	if bitmap == nil {
		return AutomationResponse{
			Success: false,
			Error:   "failed to capture screen",
		}
	}

	// 将bitmap转换为base64
	// 这里简化处理，实际应用中可能需要更复杂的图像处理
	return AutomationResponse{
		Success: true,
		Message: "screenshot taken",
		Data: map[string]interface{}{
			"format": "bitmap",
			"width":  bitmap.Width,
			"height": bitmap.Height,
		},
	}
}

func handleInfo(params map[string]interface{}) AutomationResponse {
	return AutomationResponse{
		Success: true,
		Message: "automation worker info",
		Data: map[string]interface{}{
			"version": "1.0.0",
			"engine":  "robotgo",
			"capabilities": []string{
				"click",
				"type", 
				"keypress",
				"screenshot",
			},
		},
	}
}
