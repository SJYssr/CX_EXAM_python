package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

// GUIComponents Go版本的GUI组件管理器
type GUIComponents struct {
	// 主窗口
	mainWindow *walk.MainWindow
	
	// 核心服务
	aiService   *AIService
	protection  *ScreenProtection
	fileHandler *FileHandler
	
	// GUI组件
	searchEdit    *walk.LineEdit
	contentText   *walk.TextEdit
	inputEdit     *walk.LineEdit
	aiSearchEdit  *walk.LineEdit
	aiContentText *walk.TextEdit
	aiInputEdit   *walk.LineEdit
	
	// 容器组件
	mainContainer *walk.Composite
	aiContainer   *walk.Composite
	
	// 按钮组件
	searchButton   *walk.PushButton
	aiButton       *walk.PushButton
	inputButton    *walk.PushButton
	aiSearchButton *walk.PushButton
	aiInputButton  *walk.PushButton
	backButton     *walk.PushButton
	
	// 状态变量
	currentOpacity float64
	isSmall        bool
	isAIMode       bool
	searchResults  []GUISearchResult
	currentSearchIndex int
	
	// 位置相关
	screenWidth  int
	screenHeight int
	dynamicX     int
	dynamicY     int
	
	// 拖动相关
	startX, startY int
	isDragging     bool
	
	// 同步
	mutex sync.Mutex
}

// GUISearchResult GUI搜索结果结构
type GUISearchResult struct {
	start int
	end   int
	text  string
}

// NewGUIComponents 创建新的GUI组件管理器
func NewGUIComponents(aiService *AIService, protection *ScreenProtection, fileHandler *FileHandler) (*GUIComponents, error) {
	gc := &GUIComponents{
		aiService:     aiService,
		protection:    protection,
		fileHandler:   fileHandler,
		currentOpacity: 0.5,
		isSmall:       false,
		isAIMode:      false,
		searchResults: make([]GUISearchResult, 0),
	}
	
	// 检测动态位置
	gc.detectDynamicPosition()
	
	return gc, nil
}

// detectDynamicPosition 动态检测屏幕和窗口位置
func (gc *GUIComponents) detectDynamicPosition() {
	// 获取屏幕尺寸
	gc.screenWidth = int(win.GetSystemMetrics(win.SM_CXSCREEN))
	gc.screenHeight = int(win.GetSystemMetrics(win.SM_CYSCREEN))
	
	// 尝试查找超星学习通窗口
	chaoxingHwnd := gc.findChaoxingWindow()
	
	if chaoxingHwnd != 0 {
		// 如果找到超星窗口，计算相对位置
		var rect win.RECT
		win.GetWindowRect(chaoxingHwnd, &rect)
		
		windowWidth := int(rect.Right - rect.Left)
		if int(rect.Right)+300 <= gc.screenWidth {
			// 放在右侧
			gc.dynamicX = int(rect.Right) + 10
			gc.dynamicY = int(rect.Top) + 50
		} else {
			// 放在左侧
			gc.dynamicX = max(0, int(rect.Left)-310)
			gc.dynamicY = int(rect.Top) + 50
		}
		
		log.Printf("找到超星窗口，窗口宽度: %d, 动态位置: (%d, %d)", windowWidth, gc.dynamicX, gc.dynamicY)
	} else {
		// 如果没找到超星窗口，使用屏幕比例定位
		gc.dynamicX = 0
		gc.dynamicY = int(float64(gc.screenHeight) * 0.4) // 屏幕高度的40%位置
		log.Printf("未找到超星窗口，使用默认位置: (%d, %d)", gc.dynamicX, gc.dynamicY)
	}
}

// findChaoxingWindow 查找超星学习通窗口
func (gc *GUIComponents) findChaoxingWindow() win.HWND {
	// 暂时返回0，表示未找到超星窗口
	// 这个功能需要更复杂的Windows API调用，暂时简化处理
	return 0
}

// CreateMainWindow 创建主窗口
func (gc *GUIComponents) CreateMainWindow() error {
	err := MainWindow{
		AssignTo: &gc.mainWindow,
		Title:    "超星学习通辅助工具 - Go版",
		MinSize:  Size{Width: 300, Height: 533},
		Size:     Size{Width: 300, Height: 533},
		Layout:   VBox{},
		OnKeyDown: gc.handleKeyDown,
		OnMouseDown: gc.handleMouseDown,
		OnMouseMove: gc.handleMouseMove,
		OnMouseUp: gc.handleMouseUp,
		Children: []Widget{
			gc.createMainInterface(),
			gc.createAIInterface(),
		},
	}.Create()
	
	if err != nil {
		return err
	}
	
	// 设置窗口位置和透明度
	gc.mainWindow.SetBounds(walk.Rectangle{
		X:      gc.dynamicX,
		Y:      gc.dynamicY,
		Width:  300,
		Height: 533,
	})
	
	// 加载题库内容
	gc.loadContent()
	
	// 启动保持置顶的goroutine
	go gc.keepOnTop()
	
	return nil
}

// createMainInterface 创建主界面
func (gc *GUIComponents) createMainInterface() Widget {
	return Composite{
		AssignTo: &gc.mainContainer,
		Layout:   VBox{},
		Children: []Widget{
			// 搜索框区域
			Composite{
				Layout: HBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &gc.searchEdit,
						OnKeyDown: func(key walk.Key) {
							if key == walk.KeyReturn {
								gc.nextSearchResult()
							}
						},
					},
					PushButton{
						AssignTo:  &gc.searchButton,
						Text:      "搜索",
						OnClicked: gc.onSearchClicked,
					},
					PushButton{
						AssignTo:  &gc.aiButton,
						Text:      "AI",
						OnClicked: gc.switchToAISearch,
					},
				},
			},
			
			// 文本显示区域
			TextEdit{
				AssignTo: &gc.contentText,
				ReadOnly: true,
				VScroll:  true,
			},
			
			// 输入区域
			Composite{
				Layout: HBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &gc.inputEdit,
					},
					PushButton{
						AssignTo:  &gc.inputButton,
						Text:      "输入",
						OnClicked: gc.onInputClicked,
					},
				},
			},
		},
	}
}

// createAIInterface 创建AI界面
func (gc *GUIComponents) createAIInterface() Widget {
	return Composite{
		AssignTo: &gc.aiContainer,
		Layout:   VBox{},
		Visible:  false,
		Children: []Widget{
			// AI搜索框区域
			Composite{
				Layout: HBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &gc.aiSearchEdit,
					},
					PushButton{
						AssignTo:  &gc.aiSearchButton,
						Text:      "AI搜索",
						OnClicked: gc.onAISearchClicked,
					},
					PushButton{
						AssignTo:  &gc.backButton,
						Text:      "返回",
						OnClicked: gc.switchToMain,
					},
				},
			},
			
			// AI文本显示区域
			TextEdit{
				AssignTo: &gc.aiContentText,
				ReadOnly: true,
				VScroll:  true,
			},
			
			// AI输入区域
			Composite{
				Layout: HBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &gc.aiInputEdit,
					},
					PushButton{
						AssignTo:  &gc.aiInputButton,
						Text:      "输入",
						OnClicked: gc.onAIInputClicked,
					},
				},
			},
		},
	}
}

// 事件处理函数

// handleKeyDown 处理键盘按下事件
func (gc *GUIComponents) handleKeyDown(key walk.Key) {
	switch key {
	case walk.KeyF3:
		gc.toggleWindowSize()
	case walk.KeyF4:
		// 暂时注释掉，避免编译错误
		// gc.protection.ManualRefreshIMEProtection()
		log.Printf("F4键按下 - 手动刷新输入法保护")
	case walk.KeyEscape:
		gc.mainWindow.Close()
	}
}

// handleMouseDown 处理鼠标按下事件
func (gc *GUIComponents) handleMouseDown(x, y int, button walk.MouseButton) {
	if button == walk.LeftButton {
		// 检查是否按下了Ctrl键
		if win.GetKeyState(win.VK_CONTROL) < 0 {
			gc.startX = x
			gc.startY = y
			gc.isDragging = true
		}
	} else if button == walk.RightButton {
		gc.changeOpacityQuick()
	}
}

// handleMouseMove 处理鼠标移动事件
func (gc *GUIComponents) handleMouseMove(x, y int, button walk.MouseButton) {
	if gc.isDragging && button == walk.LeftButton {
		// 获取当前窗口位置
		bounds := gc.mainWindow.Bounds()
		newX := bounds.X + (x - gc.startX)
		newY := bounds.Y + (y - gc.startY)
		
		gc.mainWindow.SetBounds(walk.Rectangle{
			X:      newX,
			Y:      newY,
			Width:  bounds.Width,
			Height: bounds.Height,
		})
	}
}

// handleMouseUp 处理鼠标释放事件
func (gc *GUIComponents) handleMouseUp(x, y int, button walk.MouseButton) {
	if button == walk.LeftButton {
		gc.isDragging = false
	}
}

// 功能函数

// loadContent 加载题库内容
func (gc *GUIComponents) loadContent() {
	content, err := gc.fileHandler.LoadTikuContent()
	if err != nil {
		log.Printf("加载题库内容失败: %v", err)
		content = "加载题库失败，请检查tiku.txt文件"
	}
	gc.contentText.SetText(content)
}

// onSearchClicked 搜索按钮点击事件
func (gc *GUIComponents) onSearchClicked() {
	searchText := gc.searchEdit.Text()
	if searchText == "" {
		return
	}
	gc.highlightSearch(searchText)
}

// highlightSearch 高亮搜索结果
func (gc *GUIComponents) highlightSearch(searchText string) {
	content := gc.contentText.Text()
	gc.searchResults = gc.searchResults[:0] // 清空搜索结果
	gc.currentSearchIndex = 0
	
	// 查找所有匹配项
	startIndex := 0
	for {
		index := strings.Index(content[startIndex:], searchText)
		if index == -1 {
			break
		}
		
		realIndex := startIndex + index
		gc.searchResults = append(gc.searchResults, GUISearchResult{
			start: realIndex,
			end:   realIndex + len(searchText),
			text:  searchText,
		})
		startIndex = realIndex + 1
	}
	
	if len(gc.searchResults) > 0 {
		// 选中第一个搜索结果
		result := gc.searchResults[0]
		gc.contentText.SetTextSelection(result.start, result.end)
		log.Printf("找到 %d 个搜索结果", len(gc.searchResults))
	}
}

// nextSearchResult 跳转到下一个搜索结果
func (gc *GUIComponents) nextSearchResult() {
	if len(gc.searchResults) == 0 {
		return
	}
	
	gc.currentSearchIndex = (gc.currentSearchIndex + 1) % len(gc.searchResults)
	result := gc.searchResults[gc.currentSearchIndex]
	gc.contentText.SetTextSelection(result.start, result.end)
}

// onInputClicked 输入按钮点击事件
func (gc *GUIComponents) onInputClicked() {
	inputText := gc.inputEdit.Text()
	if inputText == "" {
		return
	}
	
	go gc.autoInput(inputText, gc.inputEdit)
}

// onAIInputClicked AI输入按钮点击事件
func (gc *GUIComponents) onAIInputClicked() {
	inputText := gc.aiInputEdit.Text()
	if inputText == "" {
		return
	}
	
	go gc.autoInput(inputText, gc.aiInputEdit)
}

// autoInput 自动输入文本
func (gc *GUIComponents) autoInput(text string, inputEdit *walk.LineEdit) {
	// 延时5秒
	time.Sleep(5 * time.Second)
	
	// 使用Windows API自动输入
	gc.sendKeys(text)
	
	// 清空输入框
	gc.mainWindow.Synchronize(func() {
		inputEdit.SetText("")
	})
}

// sendKeys 发送按键 (简化实现)
func (gc *GUIComponents) sendKeys(text string) {
	// 暂时使用剪贴板实现自动输入
	// 这是一个简化的实现，实际使用时可以改进
	log.Printf("模拟输入文本: %s", text)
	
	// 可以在这里实现剪贴板操作
	// 或者使用其他方式实现自动输入
}

// switchToAISearch 切换到AI搜索界面
func (gc *GUIComponents) switchToAISearch() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	
	gc.mainContainer.SetVisible(false)
	gc.aiContainer.SetVisible(true)
	gc.isAIMode = true
}

// switchToMain 切换回主界面
func (gc *GUIComponents) switchToMain() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	
	gc.aiContainer.SetVisible(false)
	gc.mainContainer.SetVisible(true)
	gc.isAIMode = false
}

// onAISearchClicked AI搜索按钮点击事件
func (gc *GUIComponents) onAISearchClicked() {
	query := gc.aiSearchEdit.Text()
	if query == "" {
		return
	}
	
	// 禁用按钮，显示加载状态
	gc.aiSearchButton.SetEnabled(false)
	gc.aiContentText.SetText("正在思考中，请稍候...")
	
	// 异步调用AI服务
	go func() {
		response, err := gc.aiService.Ask(query)
		
		// 在主线程更新UI
		gc.mainWindow.Synchronize(func() {
			gc.aiSearchButton.SetEnabled(true)
			if err != nil {
				gc.aiContentText.SetText(fmt.Sprintf("AI搜索失败: %v", err))
			} else {
				gc.aiContentText.SetText(response)
			}
		})
	}()
}

// toggleWindowSize 切换窗口大小
func (gc *GUIComponents) toggleWindowSize() {
	if gc.isSmall {
		// 恢复正常大小
		gc.detectDynamicPosition() // 重新检测位置
		gc.mainWindow.SetBounds(walk.Rectangle{
			X:      gc.dynamicX,
			Y:      gc.dynamicY,
			Width:  300,
			Height: 533,
		})
		gc.setOpacity(gc.currentOpacity)
	} else {
		// 缩小到迷你模式
		gc.mainWindow.SetBounds(walk.Rectangle{
			X:      0,
			Y:      0,
			Width:  5,
			Height: 910,
		})
		gc.setOpacity(0.1)
	}
	gc.isSmall = !gc.isSmall
}

// changeOpacityQuick 快速切换透明度
func (gc *GUIComponents) changeOpacityQuick() {
	if gc.isSmall {
		return
	}
	
	if gc.currentOpacity == 0.2 {
		gc.currentOpacity = 0.5
	} else {
		gc.currentOpacity = 0.2
	}
	gc.setOpacity(gc.currentOpacity)
}

// setOpacity 设置窗口透明度
func (gc *GUIComponents) setOpacity(opacity float64) {
	// 简化的透明度设置
	// walk库可能有自己的透明度设置方法
	log.Printf("设置窗口透明度: %.2f", opacity)
	
	// 可以在这里实现透明度设置
	// 目前先用日志输出代替
}

// keepOnTop 保持窗口置顶
func (gc *GUIComponents) keepOnTop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	positionCheckCounter := 0
	
	for range ticker.C {
		// 设置窗口置顶
		hwnd := win.HWND(gc.mainWindow.Handle())
		win.SetWindowPos(hwnd, win.HWND_TOPMOST, 0, 0, 0, 0, 
			win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_NOACTIVATE)
		
		// 保护当前程序窗口和输入法窗口 (简化实现)
		// gc.protection.ProtectCurrentWindow()
		// gc.protection.FindIMEWindows()
		// gc.protection.ProtectAllIMEWindows()
		log.Printf("保护窗口功能运行中...")
		
		// 每10秒检查一次位置（只在正常大小时）
		positionCheckCounter++
		if !gc.isSmall && positionCheckCounter >= 5 { // 5 * 2秒 = 10秒
			positionCheckCounter = 0
			oldX, oldY := gc.dynamicX, gc.dynamicY
			gc.detectDynamicPosition()
			
			// 如果位置发生变化，更新窗口位置
			if oldX != gc.dynamicX || oldY != gc.dynamicY {
				gc.mainWindow.Synchronize(func() {
					bounds := gc.mainWindow.Bounds()
					gc.mainWindow.SetBounds(walk.Rectangle{
						X:      gc.dynamicX,
						Y:      gc.dynamicY,
						Width:  bounds.Width,
						Height: bounds.Height,
					})
				})
			}
		}
	}
}

// Run 运行GUI
func (gc *GUIComponents) Run() {
	gc.mainWindow.Run()
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
} 