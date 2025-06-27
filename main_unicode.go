// +build unicode

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting ChaoxingHelper Unicode Edition...")
	fmt.Println("Enhanced Chinese text support + Screen protection")
	
	// Create Unicode GUI
	gui, err := NewUnicodeGui()
	if err != nil {
		fmt.Printf("ERROR: Failed to create Unicode GUI: %v\n", err)
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		os.Exit(1)
	}
	
	fmt.Println("Unicode GUI created successfully!")
	fmt.Println("Loading with advanced UTF-8 support...")
	
	// Run GUI
	gui.Run()
} 