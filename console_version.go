package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// ConsoleApp 控制台版本的应用
type ConsoleApp struct {
	fileHandler *FileHandler
	aiService   *AIService
	content     string
}

// NewConsoleApp 创建控制台应用
func NewConsoleApp() (*ConsoleApp, error) {
	app := &ConsoleApp{}
	
	// 初始化文件处理器
	var err error
	app.fileHandler, err = NewFileHandler()
	if err != nil {
		return nil, fmt.Errorf("初始化文件处理器失败: %v", err)
	}
	
	// 初始化AI服务
	config, err := NewConfig()
	if err == nil {
		app.aiService, _ = NewAIService(config)
	}
	
	// 加载题库内容
	app.content, err = app.fileHandler.LoadTikuContent()
	if err != nil {
		app.content = "题库加载失败"
	}
	
	return app, nil
}

// Run 运行控制台应用
func (app *ConsoleApp) Run() {
	fmt.Println("========================================")
	fmt.Println("  超星学习通辅助工具 - 控制台版本")
	fmt.Println("  作者: SJYssr")
	fmt.Println("  Github: https://github.com/SJYssr")
	fmt.Println("========================================")
	fmt.Println()
	
	fmt.Println("📖 题库内容预览:")
	lines := strings.Split(app.content, "\n")
	for i, line := range lines {
		if i >= 10 { // 只显示前10行
			fmt.Printf("... (还有 %d 行内容)\n", len(lines)-10)
			break
		}
		fmt.Printf("  %s\n", line)
	}
	fmt.Println()
	
	app.showMenu()
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("👉 请选择操作 (输入数字): ")
		if !scanner.Scan() {
			break
		}
		
		choice := strings.TrimSpace(scanner.Text())
		
		switch choice {
		case "1":
			app.searchFunction(scanner)
		case "2":
			app.aiFunction(scanner)
		case "3":
			app.showStatistics()
		case "4":
			app.showHelp()
		case "5":
			fmt.Println("👋 感谢使用，再见！")
			return
		default:
			fmt.Println("❌ 无效的选择，请重新输入")
		}
		
		fmt.Println()
		app.showMenu()
	}
}

// showMenu 显示菜单
func (app *ConsoleApp) showMenu() {
	fmt.Println("🔹 可用功能:")
	fmt.Println("  1️⃣  搜索题库")
	fmt.Println("  2️⃣  AI问答")
	fmt.Println("  3️⃣  查看统计")
	fmt.Println("  4️⃣  显示帮助")
	fmt.Println("  5️⃣  退出程序")
	fmt.Println()
}

// searchFunction 搜索功能
func (app *ConsoleApp) searchFunction(scanner *bufio.Scanner) {
	fmt.Print("🔍 请输入搜索关键词: ")
	if !scanner.Scan() {
		return
	}
	
	keyword := strings.TrimSpace(scanner.Text())
	if keyword == "" {
		fmt.Println("❌ 搜索关键词不能为空")
		return
	}
	
	fmt.Printf("🔎 搜索中... 关键词: %s\n", keyword)
	
	results := app.searchInContent(keyword)
	
	if len(results) == 0 {
		fmt.Println("😔 未找到匹配的内容")
		return
	}
	
	fmt.Printf("✅ 找到 %d 个匹配结果:\n\n", len(results))
	
	for i, result := range results {
		if i >= 5 { // 最多显示5个结果
			fmt.Printf("... (还有 %d 个结果)\n", len(results)-5)
			break
		}
		
		fmt.Printf("📄 结果 %d (第%d行):\n", i+1, result.LineNumber)
		fmt.Printf("   %s\n\n", result.Line)
	}
}

// aiFunction AI问答功能
func (app *ConsoleApp) aiFunction(scanner *bufio.Scanner) {
	if app.aiService == nil {
		fmt.Println("❌ AI服务未配置或初始化失败")
		return
	}
	
	fmt.Print("🤖 请输入您的问题: ")
	if !scanner.Scan() {
		return
	}
	
	question := strings.TrimSpace(scanner.Text())
	if question == "" {
		fmt.Println("❌ 问题不能为空")
		return
	}
	
	fmt.Printf("🤔 AI正在思考您的问题: %s\n", question)
	fmt.Print("⏳ 请稍候...")
	
	// 显示思考动画
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			fmt.Print(".")
		}
	}()
	
	response, err := app.aiService.Ask(question)
	fmt.Println() // 换行
	
	if err != nil {
		fmt.Printf("❌ AI回答失败: %v\n", err)
		return
	}
	
	fmt.Println("🤖 AI回答:")
	fmt.Printf("   %s\n", response)
}

// showStatistics 显示统计信息
func (app *ConsoleApp) showStatistics() {
	stats := app.fileHandler.GetFileStats()
	
	fmt.Println("📊 题库统计信息:")
	fmt.Printf("   📁 文件大小: %v 字节\n", stats["file_size"])
	fmt.Printf("   📝 总行数: %v 行\n", stats["total_lines"])
	fmt.Printf("   ❓ 题目数量: %v 个\n", stats["question_count"])
	fmt.Printf("   💾 缓存状态: %v\n", stats["cache_exists"])
	fmt.Printf("   🕒 最后修改: %v\n", stats["last_modified"])
}

// showHelp 显示帮助信息
func (app *ConsoleApp) showHelp() {
	fmt.Println("📖 使用帮助:")
	fmt.Println()
	fmt.Println("🔍 搜索功能:")
	fmt.Println("   - 输入关键词搜索题库内容")
	fmt.Println("   - 支持中英文搜索")
	fmt.Println("   - 搜索结果会显示行号和匹配内容")
	fmt.Println()
	fmt.Println("🤖 AI问答功能:")
	fmt.Println("   - 输入问题获得AI智能回答")
	fmt.Println("   - 支持学习相关的各种问题")
	fmt.Println("   - 需要正确配置AI服务")
	fmt.Println()
	fmt.Println("📊 统计功能:")
	fmt.Println("   - 查看题库文件的基本信息")
	fmt.Println("   - 包括文件大小、题目数量等")
	fmt.Println()
	fmt.Println("⚙️ 配置文件:")
	fmt.Println("   - 配置文件采用AES加密保存")
	fmt.Println("   - 支持Deepseek和星火AI配置")
	fmt.Println()
	fmt.Println("📝 题库格式:")
	fmt.Println("   - 题目格式: 问题：[题目内容]")
	fmt.Println("   - 答案格式: 答案：[答案内容]")
	fmt.Println("   - 支持Markdown格式")
}

// searchInContent 在内容中搜索
func (app *ConsoleApp) searchInContent(keyword string) []SearchResult {
	return app.fileHandler.SearchInContent(app.content, keyword)
} 