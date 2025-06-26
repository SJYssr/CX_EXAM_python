// +build simple

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting Chaoxing Helper - Simple Edition...")
	
	// 创建简化版GUI
	gui, err := NewSimpleFyneGUI()
	if err != nil {
		fmt.Printf("Failed to create simple GUI: %v\n", err)
		fmt.Println("Press any key to exit...")
		fmt.Scanln()
		os.Exit(1)
	}
	
	fmt.Println("Simple GUI created successfully!")
	
	// 运行GUI
	gui.Run()
} 