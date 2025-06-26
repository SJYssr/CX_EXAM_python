// +build !simple

package main

import (
	"fmt"
	"log"
	"os"
	"unicode/utf8"
)

func main() {
	// 检测字符编码支持
	testStr := "测试中文显示"
	if !utf8.ValidString(testStr) || len(testStr) != utf8.RuneCountInString(testStr)*3 {
		fmt.Println("Character encoding issue detected, switching to simple mode...")
		if gui, err := NewSimpleFyneGUI(); err == nil {
			gui.Run()
			return
		}
	}

	fmt.Println("🚀 正在启动超星学习通辅助工具 - Fyne GUI版本...")

	// 初始化配置
	config, err := NewConfig()
	if err != nil {
		log.Printf("⚠️ 配置加载失败，使用默认配置: %v", err)
		// 创建一个基本的配置结构
		config = &Config{
			AIType: 0, // 禁用AI
		}
	}

	// 初始化文件处理器
	fileHandler, err := NewFileHandler()
	if err != nil {
		fmt.Printf("❌ 文件处理器初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化AI服务（如果启用）
	var aiService *AIService
	if config.GetAIType() > 0 { // 0表示禁用AI
		aiService, err = NewAIService(config)
		if err != nil {
			log.Printf("⚠️ AI服务初始化失败: %v", err)
		} else {
			fmt.Println("✅ AI服务初始化成功")
		}
	}

	// 初始化防录屏保护
	protection := NewScreenProtection()
	fmt.Println("🛡️ 防录屏保护已启用")

	// 创建Fyne GUI
	gui, err := NewFyneGUI(aiService, protection, fileHandler)
	if err != nil {
		fmt.Printf("❌ GUI初始化失败: %v\n", err)
		fmt.Println("切换到控制台版本...")
		
		// Fallback到控制台版本
		app, err := NewConsoleApp()
		if err != nil {
			fmt.Printf("❌ 控制台版本初始化也失败: %v\n", err)
			fmt.Println("按任意键退出...")
			fmt.Scanln()
			return
		}
		app.Run()
		return
	}

	// 创建主窗口
	if err := gui.CreateWindow(); err != nil {
		fmt.Printf("❌ 窗口创建失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ GUI初始化完成，启动主界面...")
	
	// 运行GUI
	gui.Run()
} 