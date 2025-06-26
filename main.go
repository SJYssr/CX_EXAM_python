package main

import (
	"fmt"
)

func main() {
	// 使用控制台版本作为fallback
	fmt.Println("🔍 检测到GUI组件问题，启动控制台版本...")
	
	app, err := NewConsoleApp()
	if err != nil {
		fmt.Printf("❌ 初始化失败: %v\n", err)
		fmt.Println("按任意键退出...")
		fmt.Scanln()
		return
	}
	
	app.Run()
} 