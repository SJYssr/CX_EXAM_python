package main

import (
	"fmt"
	"log"

	"github.com/lxn/walk"
)

type MainWindowApp struct {
	*walk.MainWindow
	config       *Config
	aiService    *AIService
	protection   *ScreenProtection
	fileHandler  *FileHandler
	guiComponents *GUIComponents
}

func main() {
	fmt.Println("超星学习通辅助工具 - Go版本")
	fmt.Println("作者: SJYssr")
	fmt.Println("Github: https://github.com/SJYssr")
	
	app := &MainWindowApp{}
	
	// 初始化各个组件
	if err := app.initialize(); err != nil {
		log.Fatal("初始化失败:", err)
	}
	
	// 创建GUI组件管理器
	if err := app.createGUIComponents(); err != nil {
		log.Fatal("创建GUI组件失败:", err)
	}
	
	// 创建主窗口
	if err := app.guiComponents.CreateMainWindow(); err != nil {
		log.Fatal("创建主窗口失败:", err)
	}
	
	// 启动防录屏保护
	app.protection.StartProtection()
	
	// 显示增强功能提示
	walk.MsgBox(app.guiComponents.mainWindow, "防录屏增强", 
		"✅ 防录屏功能已增强！\n\n"+
		"🔒 程序窗口和输入法窗口现在都受到保护\n"+
		"📹 录屏软件将无法捕获这些窗口内容\n"+
		"🛡️ 支持主流输入法：搜狗、百度、QQ、微软等\n\n"+
		"⚡ 快捷键说明：\n"+
		"🔄 F4 - 手动刷新输入法保护\n"+
		"🪟 F3 - 切换窗口大小\n"+
		"❌ ESC - 退出程序\n"+
		"🎯 Ctrl+鼠标左键 - 拖动窗口\n"+
		"🔍 Ctrl+滚轮 - 调整透明度\n"+
		"👆 右键点击 - 快速切换透明度\n"+
		"↩️ 回车键 - 跳转搜索结果", 
		walk.MsgBoxOK)
	
	// 运行GUI
	app.guiComponents.Run()
}

func (mw *MainWindowApp) initialize() error {
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

func (mw *MainWindowApp) createGUIComponents() error {
	var err error
	
	// 创建GUI组件管理器
	mw.guiComponents, err = NewGUIComponents(mw.aiService, mw.protection, mw.fileHandler)
	if err != nil {
		return fmt.Errorf("创建GUI组件管理器失败: %v", err)
	}
	
	return nil
} 