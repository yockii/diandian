package operation

// 点击操作结构
type ClickOperation struct {
	X      int    `json:"x"`      // X坐标
	Y      int    `json:"y"`      // Y坐标
	Button string `json:"button"` // 按钮类型: left, right, middle
}

// 输入操作结构
type TypeOperation struct {
	Text string `json:"text"` // 要输入的文本
}

// 文件操作结构
type FileOperation struct {
	Operation  string `json:"operation"`   // 操作类型: create, delete, move, copy
	SourcePath string `json:"source_path"` // 源文件路径
	TargetPath string `json:"target_path"` // 目标路径（移动/复制时需要）
	Content    string `json:"content"`     // 文件内容（创建时需要）
}

// 视觉分析响应结构
type VisualAnalysisResponse struct {
	ElementsFound   []VisualElement        `json:"elements_found"`  // 找到的元素
	ScreenInfo      ScreenInfo             `json:"screen_info"`     // 屏幕信息
	Recommendations []ActionRecommendation `json:"recommendations"` // 操作建议
}

// 视觉元素
type VisualElement struct {
	Type        string      `json:"type"`         // 元素类型
	Description string      `json:"description"`  // 元素描述
	Coordinates Coordinates `json:"coordinates"`  // 坐标信息
	Confidence  float64     `json:"confidence"`   // 置信度
	TextContent string      `json:"text_content"` // 文本内容
	Clickable   bool        `json:"clickable"`    // 是否可点击
}

// 屏幕信息
type ScreenInfo struct {
	Width  int `json:"width"`  // 屏幕宽度
	Height int `json:"height"` // 屏幕高度
}

// 操作建议
type ActionRecommendation struct {
	Type        string `json:"type"`        // 建议类型
	Description string `json:"description"` // 建议描述
	Priority    int    `json:"priority"`    // 优先级
}

// 坐标结构
type Coordinates struct {
	X      int `json:"x"`      // X坐标
	Y      int `json:"y"`      // Y坐标
	Width  int `json:"width"`  // 宽度
	Height int `json:"height"` // 高度
}
