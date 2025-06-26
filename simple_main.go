package main

import (
	"fmt"
	"log"
	"os"
)

func simpleMain() {
	fmt.Println("超星学习通辅助工具 - 极简测试版")
	
	// 创建简单GUI
	gui, err := NewSimpleGUI()
	if err != nil {
		log.Printf("创建GUI失败: %v", err)
		os.Exit(1)
	}
	
	// 创建窗口
	if err := gui.CreateWindow(); err != nil {
		log.Printf("创建窗口失败: %v", err)
		fmt.Printf("错误: %v\n", err)
		fmt.Println("按任意键退出...")
		fmt.Scanln()
		os.Exit(1)
	}
	
	fmt.Println("GUI创建成功，启动窗口...")
	
	// 运行GUI
	gui.Run()
} 