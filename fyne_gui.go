package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/driver/desktop"
)

// FyneGUI Fyne版本的GUI组件
type FyneGUI struct {
	app         fyne.App
	window      fyne.Window
	aiService   *AIService
	protection  *ScreenProtection
	fileHandler *FileHandler
	inputHelper *InputHelper
	
	// 主界面组件
	searchEntry   *widget.Entry
	contentText   *widget.RichText
	inputEntry    *widget.Entry
	
	// AI界面组件
	aiSearchEntry *widget.Entry
	aiContentText *widget.RichText
	aiInputEntry  *widget.Entry
	
	// 容器
	mainContainer *container.AppTabs
	
	// 状态
	isSmallMode   bool
	tikuContent   string
	searchResults []string
	currentIndex  int
}

// NewFyneGUI 创建Fyne GUI
func NewFyneGUI(aiService *AIService, protection *ScreenProtection, fileHandler *FileHandler) (*FyneGUI, error) {
	fyneApp := app.NewWithID("com.sjyssr.chaoxing-helper")
	fyneApp.SetIcon(theme.ComputerIcon()) // 设置应用图标
	
	gui := &FyneGUI{
		app:         fyneApp,
		aiService:   aiService,
		protection:  protection,
		fileHandler: fileHandler,
		inputHelper: NewInputHelper(),
		searchResults: make([]string, 0),
	}
	
	return gui, nil
}

// CreateWindow 创建主窗口
func (fg *FyneGUI) CreateWindow() error {
	fg.window = fg.app.NewWindow("超星学习通辅助工具 - Fyne版")
	fg.window.Resize(fyne.NewSize(500, 700))
	fg.window.CenterOnScreen()
	
	// 设置窗口图标
	fg.window.SetIcon(theme.ComputerIcon())
	
	// 加载题库内容
	if err := fg.loadContent(); err != nil {
		log.Printf("加载题库失败: %v", err)
		fg.tikuContent = "题库加载失败，请检查tiku.txt文件"
	}
	
	fmt.Println("✅ 题库内容加载完成")
	
	// 创建界面
	fg.createMainInterface()
	fg.createAIInterface()
	
	// 设置快捷键
	fg.setupShortcuts()
	
	// 启动保护功能
	if fg.protection != nil {
		go fg.protection.StartProtection()
		fmt.Println("🛡️ 防录屏保护已启动")
	}
	
	// 延迟显示欢迎消息，确保界面完全加载
	go func() {
		time.Sleep(500 * time.Millisecond)
		fg.showWelcomeDialog()
	}()
	
	return nil
}

// createMainInterface 创建主界面
func (fg *FyneGUI) createMainInterface() {
	// 搜索区域
	fg.searchEntry = widget.NewEntry()
	fg.searchEntry.SetPlaceHolder("输入搜索关键词...")
	
	searchBtn := widget.NewButton("🔍 搜索", fg.onSearchClicked)
	searchBtn.Importance = widget.HighImportance
	
	searchContainer := container.NewBorder(nil, nil, nil, searchBtn, fg.searchEntry)
	
	// 内容显示区域
	fg.contentText = widget.NewRichText()
	fg.contentText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text:  fg.tikuContent,
			Style: widget.RichTextStyle{},
		},
	}
	fg.contentText.Wrapping = fyne.TextWrapWord
	
	contentScroll := container.NewScroll(fg.contentText)
	contentScroll.SetMinSize(fyne.NewSize(350, 300))
	
	// 输入区域
	fg.inputEntry = widget.NewEntry()
	fg.inputEntry.SetPlaceHolder("输入要自动输入的内容...")
	
	inputBtn := widget.NewButton("⌨️ 输入", fg.onInputClicked)
	inputBtn.Importance = widget.MediumImportance
	
	inputContainer := container.NewBorder(nil, nil, nil, inputBtn, fg.inputEntry)
	
	// 功能按钮区域
	toggleBtn := widget.NewButton("🪟 切换大小", fg.toggleWindowSize)
	exitBtn := widget.NewButton("❌ 退出", func() { fg.window.Close() })
	
	buttonsContainer := container.NewHBox(
		toggleBtn,
		widget.NewSeparator(),
		exitBtn,
	)
	
	// 组合主界面
	mainContent := container.NewVBox(
		widget.NewCard("", "🔍 搜索功能", searchContainer),
		widget.NewCard("", "📖 题库内容", contentScroll),
		widget.NewCard("", "⌨️ 自动输入", inputContainer),
		widget.NewSeparator(),
		buttonsContainer,
	)
	
	// 添加到标签页
	if fg.mainContainer == nil {
		fg.mainContainer = container.NewAppTabs()
	}
	fg.mainContainer.Append(container.NewTabItem("📚 主界面", mainContent))
}

// createAIInterface 创建AI界面
func (fg *FyneGUI) createAIInterface() {
	// AI搜索区域
	fg.aiSearchEntry = widget.NewEntry()
	fg.aiSearchEntry.SetPlaceHolder("输入您的问题...")
	
	aiSearchBtn := widget.NewButton("🤖 AI搜索", fg.onAISearchClicked)
	aiSearchBtn.Importance = widget.HighImportance
	
	aiSearchContainer := container.NewBorder(nil, nil, nil, aiSearchBtn, fg.aiSearchEntry)
	
	// AI回复显示区域
	fg.aiContentText = widget.NewRichText()
	fg.aiContentText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text:  "AI功能已就绪，请输入问题获取智能回答。\n\n支持学习相关的各种问题...",
			Style: widget.RichTextStyle{},
		},
	}
	fg.aiContentText.Wrapping = fyne.TextWrapWord
	
	aiContentScroll := container.NewScroll(fg.aiContentText)
	aiContentScroll.SetMinSize(fyne.NewSize(350, 300))
	
	// AI输入区域  
	fg.aiInputEntry = widget.NewEntry()
	fg.aiInputEntry.SetPlaceHolder("AI回答的内容将自动输入...")
	
	aiInputBtn := widget.NewButton("⌨️ 输入AI回答", fg.onAIInputClicked)
	aiInputBtn.Importance = widget.MediumImportance
	
	aiInputContainer := container.NewBorder(nil, nil, nil, aiInputBtn, fg.aiInputEntry)
	
	// 组合AI界面
	aiContent := container.NewVBox(
		widget.NewCard("", "🤖 AI问答", aiSearchContainer),
		widget.NewCard("", "💬 AI回答", aiContentScroll),
		widget.NewCard("", "⌨️ 自动输入", aiInputContainer),
	)
	
	// 添加到标签页
	fg.mainContainer.Append(container.NewTabItem("🤖 AI助手", aiContent))
	
	// 设置窗口内容
	fg.window.SetContent(fg.mainContainer)
}

// setupShortcuts 设置快捷键
func (fg *FyneGUI) setupShortcuts() {
	// F3 - 切换窗口大小
	fg.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyF3,
		Modifier: 0,
	}, func(shortcut fyne.Shortcut) {
		fg.toggleWindowSize()
	})
	
	// ESC - 退出
	fg.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyEscape,
		Modifier: 0,
	}, func(shortcut fyne.Shortcut) {
		fg.window.Close()
	})
	
	// Enter - 搜索
	fg.searchEntry.OnSubmitted = func(text string) {
		fg.onSearchClicked()
	}
	
	// AI搜索Enter
	fg.aiSearchEntry.OnSubmitted = func(text string) {
		fg.onAISearchClicked()
	}
}

// loadContent 加载题库内容
func (fg *FyneGUI) loadContent() error {
	content, err := fg.fileHandler.LoadTikuContent()
	if err != nil {
		return err
	}
	fg.tikuContent = content
	return nil
}

// onSearchClicked 搜索按钮点击
func (fg *FyneGUI) onSearchClicked() {
	searchText := fg.searchEntry.Text
	if searchText == "" {
		dialog.ShowInformation("提示", "请输入搜索关键词", fg.window)
		return
	}
	
	// 执行搜索
	results := fg.searchInContent(searchText)
	
	if len(results) == 0 {
		dialog.ShowInformation("搜索结果", "未找到匹配的内容", fg.window)
		return
	}
	
	// 显示搜索结果
	fg.showSearchResults(results, searchText)
}

// searchInContent 搜索内容
func (fg *FyneGUI) searchInContent(keyword string) []string {
	var results []string
	lines := strings.Split(fg.tikuContent, "\n")
	keyword = strings.ToLower(keyword)
	
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), keyword) {
			results = append(results, fmt.Sprintf("第%d行: %s", i+1, line))
		}
	}
	
	return results
}

// showSearchResults 显示搜索结果
func (fg *FyneGUI) showSearchResults(results []string, keyword string) {
	fg.searchResults = results
	fg.currentIndex = 0
	
	// 更新内容显示
	highlightedContent := fg.highlightKeyword(fg.tikuContent, keyword)
	fg.contentText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text:  highlightedContent,
			Style: widget.RichTextStyle{},
		},
	}
	fg.contentText.Refresh()
	
	// 显示结果统计
	dialog.ShowInformation("搜索结果", 
		fmt.Sprintf("找到 %d 个匹配结果\n关键词: %s", len(results), keyword), 
		fg.window)
}

// highlightKeyword 高亮关键词
func (fg *FyneGUI) highlightKeyword(content, keyword string) string {
	// 简单的高亮实现，用【】标记关键词
	return strings.ReplaceAll(content, keyword, fmt.Sprintf("【%s】", keyword))
}

// onInputClicked 输入按钮点击
func (fg *FyneGUI) onInputClicked() {
	inputText := fg.inputEntry.Text
	if inputText == "" {
		dialog.ShowInformation("提示", "请输入要自动输入的内容", fg.window)
		return
	}
	
	// 显示倒计时
	fg.showCountdown(func() {
		fg.autoInput(inputText)
		fg.inputEntry.SetText("") // 清空输入框
	})
}

// onAISearchClicked AI搜索按钮点击
func (fg *FyneGUI) onAISearchClicked() {
	query := fg.aiSearchEntry.Text
	if query == "" {
		dialog.ShowInformation("提示", "请输入您的问题", fg.window)
		return
	}
	
	if fg.aiService == nil {
		dialog.ShowError(fmt.Errorf("AI服务未配置"), fg.window)
		return
	}
	
	// 显示加载状态
	fg.aiContentText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text:  "🤔 AI正在思考中，请稍候...",
			Style: widget.RichTextStyle{},
		},
	}
	fg.aiContentText.Refresh()
	
	// 异步调用AI
	go func() {
		response, err := fg.aiService.Ask(query)
		
		// 更新UI
		if err != nil {
			fg.aiContentText.Segments = []widget.RichTextSegment{
				&widget.TextSegment{
					Text:  fmt.Sprintf("❌ AI搜索失败: %v", err),
					Style: widget.RichTextStyle{},
				},
			}
		} else {
			fg.aiContentText.Segments = []widget.RichTextSegment{
				&widget.TextSegment{
					Text:  fmt.Sprintf("🤖 AI回答:\n\n%s", response),
					Style: widget.RichTextStyle{},
				},
			}
			// 自动填充到输入框
			fg.aiInputEntry.SetText(response)
		}
		fg.aiContentText.Refresh()
	}()
}

// onAIInputClicked AI输入按钮点击
func (fg *FyneGUI) onAIInputClicked() {
	inputText := fg.aiInputEntry.Text
	if inputText == "" {
		dialog.ShowInformation("提示", "没有可输入的AI回答内容", fg.window)
		return
	}
	
	// 显示倒计时
	fg.showCountdown(func() {
		fg.autoInput(inputText)
	})
}

// showCountdown 显示倒计时
func (fg *FyneGUI) showCountdown(callback func()) {
	// 先显示一个通知
	dialog.ShowInformation("准备中", "5秒后开始自动输入，请将鼠标移到目标位置...", fg.window)
	
	go func() {
		// 简单的延迟，不使用动态更新的对话框
		time.Sleep(5 * time.Second)
		callback()
	}()
}

// autoInput 自动输入
func (fg *FyneGUI) autoInput(text string) {
	// 最小化窗口，让用户点击需要输入的地方
	fg.window.Hide()
	
	// 延迟一秒让用户切换到目标窗口
	time.Sleep(1 * time.Second)
	
	// 执行自动输入
	if err := fg.inputHelper.TypeText(text); err != nil {
		// 如果输入失败，显示错误并恢复窗口
		fg.window.Show()
		dialog.ShowError(fmt.Errorf("自动输入失败: %v\n已复制到剪贴板作为备选", err), fg.window)
		fg.window.Clipboard().SetContent(text)
	} else {
		// 输入成功，短暂延迟后恢复窗口
		time.Sleep(500 * time.Millisecond)
		fg.window.Show()
		dialog.ShowInformation("完成", 
			fmt.Sprintf("✅ 自动输入完成\n输入内容: %s", text), 
			fg.window)
	}
}

// toggleWindowSize 切换窗口大小
func (fg *FyneGUI) toggleWindowSize() {
	if fg.isSmallMode {
		fg.window.Resize(fyne.NewSize(400, 600))
		fg.isSmallMode = false
	} else {
		fg.window.Resize(fyne.NewSize(200, 100))
		fg.isSmallMode = true
	}
}

// showWelcomeDialog 显示欢迎对话框
func (fg *FyneGUI) showWelcomeDialog() {
	welcomeText := "✅ 超星学习通辅助工具已启动！\n\n" +
		"🔍 主要功能:\n" +
		"• 题库搜索和查看\n" +
		"• AI智能问答\n" +
		"• 自动输入功能\n\n" +
		"⌨️ 快捷键:\n" +
		"• F3 - 切换窗口大小\n" +
		"• ESC - 退出程序\n" +
		"• Enter - 执行搜索\n\n" +
		"🛡️ 防录屏保护已启动\n" +
		"📚 题库已加载完成\n\n" +
		"祝您学习愉快！"

	dialog.ShowInformation("欢迎使用 - 超星学习通辅助工具", welcomeText, fg.window)
}

// Run 运行GUI
func (fg *FyneGUI) Run() {
	fg.window.ShowAndRun()
} 