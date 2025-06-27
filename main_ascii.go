// +build ascii

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting ChaoxingHelper ASCII Edition...")
	fmt.Println("Pure English interface to avoid encoding issues")
	
	// Create ASCII GUI
	gui, err := NewASCIIGui()
	if err != nil {
		fmt.Printf("ERROR: Failed to create ASCII GUI: %v\n", err)
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		os.Exit(1)
	}
	
	fmt.Println("ASCII GUI created successfully!")
	fmt.Println("Launching interface...")
	
	// Run GUI
	gui.Run()
} 