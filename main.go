package main

import (
	"fmt"
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MainWindow struct {
	*walk.MainWindow
	config      *Config
	aiService   *AIService
	protection  *ScreenProtection
	fileHandler *FileHandler
	gui         *GUIManager
}

func main() {
	fmt.Println("超星学习通辅助工具 - Go版本")
	fmt.Println("作者: SJYssr")
	fmt.Println("Github: https://github.com/SJYssr")
	
	app := &MainWindow{}
	
	// 初始化各个组件
	if err := app.initialize(); err != nil {
		log.Fatal("初始化失败:", err)
	}
	
	// 创建GUI
	if err := app.createMainWindow(); err != nil {
		log.Fatal("创建主窗口失败:", err)
	}
	
	// 启动防录屏保护
	app.protection.StartProtection()
	
	// 显示增强功能提示
	walk.MsgBox(app, "防录屏增强", 
		"✅ 防录屏功能已增强！\n\n"+
		"🔒 程序窗口和输入法窗口现在都受到保护\n"+
		"📹 录屏软件将无法捕获这些窗口内容\n"+
		"🛡️ 支持主流输入法：搜狗、百度、QQ、微软等\n\n"+
		"⚡ 快捷键说明：\n"+
		"🔄 F4 - 手动刷新输入法保护\n"+
		"🪟 F3 - 切换窗口大小\n"+
		"❌ ESC - 退出程序", 
		walk.MsgBoxOK)
	
	// 运行主循环
	app.Run()
}

func (mw *MainWindow) initialize() error {
	var err error
	
	// 初始化配置
	mw.config, err = NewConfig()
	if err != nil {
		return fmt.Errorf("初始化配置失败: %v", err)
	}
	
	// 初始化AI服务
	mw.aiService, err = NewAIService(mw.config)
	if err != nil {
		return fmt.Errorf("初始化AI服务失败: %v", err)
	}
	
	// 初始化屏幕保护
	mw.protection, err = NewScreenProtection()
	if err != nil {
		return fmt.Errorf("初始化屏幕保护失败: %v", err)
	}
	
	// 初始化文件处理器
	mw.fileHandler, err = NewFileHandler()
	if err != nil {
		return fmt.Errorf("初始化文件处理器失败: %v", err)
	}
	
	return nil
}

func (mw *MainWindow) createMainWindow() error {
	var searchEdit *walk.LineEdit
	var contentText *walk.TextEdit
	var inputEdit *walk.LineEdit
	
	err := MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "超星学习通辅助工具 - Go版",
		MinSize:  Size{Width: 300, Height: 533},
		Size:     Size{Width: 300, Height: 533},
		Layout:   VBox{},
		OnKeyDown: mw.handleKeyDown,
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{
					// 搜索框
					LineEdit{
						AssignTo: &searchEdit,
						Text:     "",
					},
					
					// 按钮组
					Composite{
						Layout: HBox{},
						Children: []Widget{
							PushButton{
								Text: "搜索",
								OnClicked: mw.onSearchClicked,
							},
							PushButton{
								Text: "AI搜索",
								OnClicked: mw.onAISearchClicked,
							},
						},
					},
					
					// 文本显示区域
					TextEdit{
						AssignTo: &contentText,
						ReadOnly: true,
						VScroll:  true,
					},
					
					// 输入框
					LineEdit{
						AssignTo: &inputEdit,
						Text:     "",
					},
					
					// 输入按钮
					PushButton{
						Text: "输入",
						OnClicked: mw.onInputClicked,
					},
				},
			},
		},
	}.Create()
	
	if err != nil {
		return err
	}
	
	// 初始化GUI管理器
	mw.gui, err = NewGUIManager(mw.MainWindow, mw.protection)
	if err != nil {
		return err
	}
	
	// 设置GUI组件引用
	mw.gui.SearchEdit = searchEdit
	mw.gui.ContentText = contentText
	mw.gui.InputEdit = inputEdit
	
	// 设置动态位置
	mw.gui.SetDynamicPosition()
	
	// 加载题库内容
	content, err := mw.fileHandler.LoadTikuContent()
	if err != nil {
		log.Printf("加载题库内容失败: %v", err)
	} else {
		mw.gui.ContentText.SetText(content)
	}
	
	return nil
}

func (mw *MainWindow) handleKeyDown(key walk.Key) {
	switch key {
	case walk.KeyF3:
		mw.gui.ToggleWindowSize()
	case walk.KeyF4:
		mw.protection.ManualRefreshIMEProtection()
	case walk.KeyEscape:
		mw.Close()
	}
}

func (mw *MainWindow) onSearchClicked() {
	if mw.gui.SearchEdit == nil {
		return
	}
	searchText := mw.gui.SearchEdit.Text()
	if searchText != "" {
		mw.gui.HighlightSearch(searchText)
	}
}

func (mw *MainWindow) onAISearchClicked() {
	if mw.gui.SearchEdit == nil {
		return
	}
	query := mw.gui.SearchEdit.Text()
	if query == "" {
		return
	}
	
	// 异步调用AI服务
	go func() {
		response, err := mw.aiService.Ask(query)
		if err != nil {
			walk.MsgBox(mw, "AI搜索错误", fmt.Sprintf("AI搜索失败: %v", err), walk.MsgBoxIconError)
			return
		}
		
		// 在主线程更新UI
		mw.Synchronize(func() {
			mw.gui.ShowAIResponse(response)
		})
	}()
}

func (mw *MainWindow) onInputClicked() {
	if mw.gui.InputEdit == nil {
		return
	}
	inputText := mw.gui.InputEdit.Text()
	if inputText != "" {
		go mw.gui.AutoInput(inputText)
	}
} 