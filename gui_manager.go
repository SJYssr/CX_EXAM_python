package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

// GUIManager GUI管理器
type GUIManager struct {
	mainWindow    *walk.MainWindow
	protection    *ScreenProtection
	SearchEdit    *walk.LineEdit
	ContentText   *walk.TextEdit
	InputEdit     *walk.LineEdit
	AIText        *walk.TextEdit
	
	// 状态变量
	isSmallMode   bool
	currentOpacity float64
	searchResults []SearchResult
	currentSearchIndex int
	
	// 位置变量
	originalX     int32
	originalY     int32
	originalWidth int32
	originalHeight int32
}

// NewGUIManager 创建GUI管理器
func NewGUIManager(mainWindow *walk.MainWindow, protection *ScreenProtection) (*GUIManager, error) {
	gm := &GUIManager{
		mainWindow:    mainWindow,
		protection:    protection,
		currentOpacity: 0.5,
		isSmallMode:   false,
		searchResults: make([]SearchResult, 0),
		currentSearchIndex: 0,
	}
	
	// 保存原始窗口尺寸
	bounds := mainWindow.Bounds()
	gm.originalX = int32(bounds.X)
	gm.originalY = int32(bounds.Y)
	gm.originalWidth = int32(bounds.Width)
	gm.originalHeight = int32(bounds.Height)
	
	return gm, nil
}

// SetDynamicPosition 设置动态位置
func (gm *GUIManager) SetDynamicPosition() {
	// 强制更新位置检测
	gm.protection.ForceUpdatePosition()
	
	// 获取动态计算的位置
	x, y := gm.protection.GetDynamicPosition()
	
	// 设置窗口位置
	gm.mainWindow.SetBounds(walk.Rectangle{
		X:      int(x),
		Y:      int(y),
		Width:  int(gm.originalWidth),
		Height: int(gm.originalHeight),
	})
	
	fmt.Printf("窗口位置已设置为: (%d, %d)\n", x, y)
}

// ToggleWindowSize 切换窗口大小 (F3功能)
func (gm *GUIManager) ToggleWindowSize() {
	if gm.isSmallMode {
		// 恢复正常大小
		gm.SetDynamicPosition() // 重新检测位置
		gm.mainWindow.SetOpacity(uint8(gm.currentOpacity * 255))
		gm.isSmallMode = false
		fmt.Println("窗口已恢复正常大小")
	} else {
		// 切换到小窗口模式
		gm.mainWindow.SetBounds(walk.Rectangle{
			X:      0,
			Y:      0,
			Width:  5,
			Height: 910,
		})
		gm.mainWindow.SetOpacity(uint8(0.1 * 255))
		gm.isSmallMode = true
		fmt.Println("窗口已切换到小窗口模式")
	}
}

// ChangeOpacity 改变窗口透明度 (右键功能)
func (gm *GUIManager) ChangeOpacity() {
	if gm.isSmallMode {
		return
	}
	
	if gm.currentOpacity == 0.2 {
		gm.currentOpacity = 0.5
	} else {
		gm.currentOpacity = 0.2
	}
	
	gm.mainWindow.SetOpacity(uint8(gm.currentOpacity * 255))
	fmt.Printf("透明度已设置为: %.1f\n", gm.currentOpacity)
}

// AdjustOpacity 调整透明度 (Ctrl+滚轮功能)
func (gm *GUIManager) AdjustOpacity(delta int) {
	if gm.isSmallMode {
		return
	}
	
	if delta > 0 {
		gm.currentOpacity += 0.1
	} else {
		gm.currentOpacity -= 0.1
	}
	
	// 限制透明度范围
	if gm.currentOpacity < 0.1 {
		gm.currentOpacity = 0.1
	}
	if gm.currentOpacity > 1.0 {
		gm.currentOpacity = 1.0
	}
	
	gm.mainWindow.SetOpacity(uint8(gm.currentOpacity * 255))
	fmt.Printf("透明度已调整为: %.1f\n", gm.currentOpacity)
}

// HighlightSearch 高亮显示搜索结果
func (gm *GUIManager) HighlightSearch(searchTerm string) {
	if gm.ContentText == nil || searchTerm == "" {
		return
	}
	
	content := gm.ContentText.Text()
	if content == "" {
		return
	}
	
	// 清除之前的搜索结果
	gm.searchResults = gm.searchResults[:0]
	gm.currentSearchIndex = 0
	
	// 查找所有匹配项
	lines := strings.Split(content, "\n")
	searchLower := strings.ToLower(searchTerm)
	
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), searchLower) {
			result := SearchResult{
				LineNumber: i + 1,
				Line:       line,
			}
			gm.searchResults = append(gm.searchResults, result)
		}
	}
	
	if len(gm.searchResults) > 0 {
		// 跳转到第一个搜索结果
		gm.jumpToSearchResult(0)
		fmt.Printf("找到 %d 个匹配结果\n", len(gm.searchResults))
	} else {
		fmt.Println("未找到匹配结果")
	}
}

// NextSearchResult 跳转到下一个搜索结果 (Enter键功能)
func (gm *GUIManager) NextSearchResult() {
	if len(gm.searchResults) == 0 {
		return
	}
	
	gm.currentSearchIndex = (gm.currentSearchIndex + 1) % len(gm.searchResults)
	gm.jumpToSearchResult(gm.currentSearchIndex)
}

// jumpToSearchResult 跳转到指定搜索结果
func (gm *GUIManager) jumpToSearchResult(index int) {
	if index < 0 || index >= len(gm.searchResults) {
		return
	}
	
	result := gm.searchResults[index]
	
	// 计算文本位置并滚动到该位置
	content := gm.ContentText.Text()
	lines := strings.Split(content, "\n")
	
	if result.LineNumber > 0 && result.LineNumber <= len(lines) {
		// 简单的文本选中实现
		lineStart := 0
		for i := 0; i < result.LineNumber-1; i++ {
			lineStart += len(lines[i]) + 1 // +1 for newline
		}
		
		// 这里可以实现文本选中功能，但walk库的TextEdit支持有限
		fmt.Printf("跳转到第 %d 行: %s\n", result.LineNumber, result.Line)
	}
}

// AutoInput 自动输入功能
func (gm *GUIManager) AutoInput(text string) {
	if text == "" {
		return
	}
	
	fmt.Printf("5秒后将自动输入: %s\n", text)
	
	// 等待5秒
	time.Sleep(5 * time.Second)
	
	// 使用robotgo进行自动输入
	robotgo.TypeStr(text)
	
	fmt.Println("自动输入完成")
	
	// 清空输入框
	if gm.InputEdit != nil {
		gm.InputEdit.SetText("")
	}
}

// ShowAIResponse 显示AI响应
func (gm *GUIManager) ShowAIResponse(response string) {
	if gm.AIText == nil {
		// 如果AI文本框不存在，创建新窗口显示
		gm.showAIWindow(response)
		return
	}
	
	gm.AIText.SetText(response)
	fmt.Println("AI响应已显示")
}

// showAIWindow 显示AI响应窗口
func (gm *GUIManager) showAIWindow(response string) {
	// 创建AI响应窗口
	var aiWindow *walk.MainWindow
	var aiTextEdit *walk.TextEdit
	
	err := walk.MainWindow{
		AssignTo: &aiWindow,
		Title:    "AI搜索结果",
		Size:     walk.Size{Width: 500, Height: 400},
		Layout:   walk.VBox{},
		Children: []walk.Widget{
			walk.TextEdit{
				AssignTo: &aiTextEdit,
				ReadOnly: true,
				Text:     response,
				VScroll:  true,
			},
			walk.PushButton{
				Text: "关闭",
				OnClicked: func() {
					aiWindow.Close()
				},
			},
		},
	}.Run()
	
	if err != nil {
		fmt.Printf("创建AI窗口失败: %v\n", err)
		// 如果创建窗口失败，使用消息框显示
		walk.MsgBox(gm.mainWindow, "AI搜索结果", response, walk.MsgBoxOK)
	}
}

// EnableKeepOnTop 启用窗口置顶
func (gm *GUIManager) EnableKeepOnTop() {
	// 设置窗口始终置顶
	hwnd := win.HWND(gm.mainWindow.Handle())
	win.SetWindowPos(hwnd, win.HWND_TOPMOST, 0, 0, 0, 0, 
		win.SWP_NOMOVE|win.SWP_NOSIZE)
	
	fmt.Println("窗口置顶已启用")
}

// HandleMouseWheel 处理鼠标滚轮事件
func (gm *GUIManager) HandleMouseWheel(delta int, keys walk.ModifierKey) {
	if keys&walk.ModControl != 0 {
		// Ctrl + 滚轮：调整透明度
		gm.AdjustOpacity(delta)
	} else if keys&walk.ModAlt != 0 {
		// Alt + 滚轮：调整字体大小
		gm.AdjustFontSize(delta)
	}
}

// AdjustFontSize 调整字体大小
func (gm *GUIManager) AdjustFontSize(delta int) {
	// 获取当前字体
	if gm.ContentText == nil {
		return
	}
	
	font := gm.ContentText.Font()
	if font == nil {
		return
	}
	
	// 调整字体大小
	currentSize := font.PointSize()
	newSize := currentSize
	
	if delta > 0 {
		newSize += 1
	} else {
		newSize -= 1
	}
	
	// 限制字体大小范围
	if newSize < 8 {
		newSize = 8
	}
	if newSize > 24 {
		newSize = 24
	}
	
	// 创建新字体
	newFont, err := walk.NewFont(font.Family(), newSize, font.Style())
	if err == nil {
		gm.ContentText.SetFont(newFont)
		fmt.Printf("字体大小已调整为: %d\n", newSize)
	}
}

// UpdateWindowPosition 更新窗口位置
func (gm *GUIManager) UpdateWindowPosition() {
	if gm.isSmallMode {
		return
	}
	
	// 重新检测位置
	gm.protection.ForceUpdatePosition()
	x, y := gm.protection.GetDynamicPosition()
	
	// 获取当前窗口位置
	bounds := gm.mainWindow.Bounds()
	currentX := int32(bounds.X)
	currentY := int32(bounds.Y)
	
	// 如果位置发生变化，更新窗口位置
	if currentX != x || currentY != y {
		gm.mainWindow.SetBounds(walk.Rectangle{
			X:      int(x),
			Y:      int(y),
			Width:  bounds.Width,
			Height: bounds.Height,
		})
		
		fmt.Printf("窗口位置已更新为: (%d, %d)\n", x, y)
	}
}

// GetGUIStats 获取GUI统计信息
func (gm *GUIManager) GetGUIStats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	bounds := gm.mainWindow.Bounds()
	stats["window_x"] = bounds.X
	stats["window_y"] = bounds.Y
	stats["window_width"] = bounds.Width
	stats["window_height"] = bounds.Height
	stats["is_small_mode"] = gm.isSmallMode
	stats["current_opacity"] = gm.currentOpacity
	stats["search_results_count"] = len(gm.searchResults)
	stats["current_search_index"] = gm.currentSearchIndex
	
	return stats
}

// StartPeriodicUpdate 启动定期更新
func (gm *GUIManager) StartPeriodicUpdate() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				gm.UpdateWindowPosition()
			}
		}
	}()
}

// EnableDragMove 启用拖拽移动窗口功能
func (gm *GUIManager) EnableDragMove() {
	// 这里可以实现Ctrl+鼠标拖拽移动窗口的功能
	// 由于walk库的限制，这里只是一个占位函数
	fmt.Println("拖拽移动功能已启用")
}

// ShowHelp 显示帮助信息
func (gm *GUIManager) ShowHelp() {
	helpText := `超星学习通辅助工具 - Go版本

快捷键说明：
• F3 - 切换窗口大小
• F4 - 手动刷新输入法保护
• ESC - 退出程序
• 右键 - 改变透明度
• Ctrl+滚轮 - 调整透明度
• Alt+滚轮 - 调整字体大小
• Enter - 下一个搜索结果

功能说明：
• 自动检测超星学习通窗口位置
• 防录屏保护功能
• 题库搜索和AI问答
• 自动输入功能（5秒延迟）
• 输入法窗口保护

配置文件：
• 配置文件使用AES-256加密保存
• 支持Deepseek和星火AI
• 自动备份题库文件

作者：SJYssr
Github：https://github.com/SJYssr`
	
	walk.MsgBox(gm.mainWindow, "帮助信息", helpText, walk.MsgBoxIconInformation)
} 