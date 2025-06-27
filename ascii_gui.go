package main

import (
	"fmt"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ASCII GUI - Pure English interface to avoid encoding issues
type ASCIIGui struct {
	app           fyne.App
	window        fyne.Window
	fileHandler   *FileHandler
	protection    *ScreenProtection
	encodingFixer *EncodingFixer
	
	searchEntry *widget.Entry
	resultText  *widget.Entry
	inputEntry  *widget.Entry
	tikuContent string
}

// Create new ASCII GUI
func NewASCIIGui() (*ASCIIGui, error) {
	app := app.New()
	
	fileHandler, err := NewFileHandler()
	if err != nil {
		return nil, err
	}
	
	// Initialize screen protection
	protection := NewScreenProtection()
	
	// Initialize encoding fixer
	encodingFixer := NewEncodingFixer()
	
	return &ASCIIGui{
		app:           app,
		fileHandler:   fileHandler,
		protection:    protection,
		encodingFixer: encodingFixer,
	}, nil
}

// Run the ASCII GUI
func (ag *ASCIIGui) Run() {
	// Create window
	ag.window = ag.app.NewWindow("ChaoxingHelper ASCII Edition v1.0")
	ag.window.Resize(fyne.NewSize(700, 500))
	ag.window.CenterOnScreen()
	
	// Start screen protection
	if ag.protection != nil {
		go ag.protection.StartProtection()
		fmt.Println("Screen protection started")
	}
	
	// Load content
	if err := ag.loadContent(); err != nil {
		ag.tikuContent = "Error: Failed to load tiku.txt file\nPlease check if the file exists."
	}
	
	// Create interface
	ag.createInterface()
	
	// Apply screen protection to this window
	ag.protectCurrentWindow()
	
	// Show welcome message
	ag.showWelcome()
	
	// Show window
	ag.window.ShowAndRun()
}

func (ag *ASCIIGui) loadContent() error {
	// Use encoding fixer to load content with proper encoding
	content, err := ag.encodingFixer.LoadTikuWithEncoding()
	if err != nil {
		// Fallback to regular file handler
		content, err = ag.fileHandler.LoadTikuContent()
		if err != nil {
			return err
		}
		// Try to fix any encoding issues
		content = ag.encodingFixer.FixContent(content)
	}
	
	// Add a header to indicate content is loaded
	header := "=== Question Database Loaded Successfully ===\n"
	header += "Total content length: " + fmt.Sprintf("%d characters", len(content)) + "\n"
	header += "Encoding: UTF-8 (Fixed)\n"
	header += "=== Content Preview ===\n\n"
	
	ag.tikuContent = header + content
	
	fmt.Printf("Content loaded with encoding fix: %d characters\n", len(content))
	return nil
}

func (ag *ASCIIGui) createInterface() {
	// Search section
	searchLabel := widget.NewLabel("SEARCH KEYWORDS:")
	ag.searchEntry = widget.NewEntry()
	ag.searchEntry.SetPlaceHolder("Enter keywords to search...")
	
	searchBtn := widget.NewButton("SEARCH", ag.onSearch)
	clearBtn := widget.NewButton("CLEAR", ag.onClear)
	
	searchContainer := container.NewHBox(
		ag.searchEntry,
		searchBtn,
		clearBtn,
	)
	
	// Results section
	resultLabel := widget.NewLabel("SEARCH RESULTS / CONTENT:")
	ag.resultText = widget.NewMultiLineEntry()
	ag.resultText.SetText(ag.tikuContent)
	ag.resultText.Wrapping = fyne.TextWrapWord
	
	resultScroll := container.NewScroll(ag.resultText)
	resultScroll.SetMinSize(fyne.NewSize(650, 200))
	
	// Input section
	inputLabel := widget.NewLabel("AUTO INPUT:")
	ag.inputEntry = widget.NewEntry()
	ag.inputEntry.SetPlaceHolder("Content to auto-input...")
	
	inputBtn := widget.NewButton("COPY TO CLIPBOARD", ag.onAutoInput)
	
	inputContainer := container.NewHBox(
		ag.inputEntry,
		inputBtn,
	)
	
	// Control buttons
	aboutBtn := widget.NewButton("ABOUT", ag.showAbout)
	helpBtn := widget.NewButton("HELP", ag.showHelp)
	exitBtn := widget.NewButton("EXIT", func() {
		ag.window.Close()
	})
	
	controlContainer := container.NewHBox(
		aboutBtn,
		helpBtn,
		widget.NewSeparator(),
		exitBtn,
	)
	
	// Main layout
	content := container.NewVBox(
		searchLabel,
		searchContainer,
		widget.NewSeparator(),
		resultLabel,
		resultScroll,
		widget.NewSeparator(),
		inputLabel,
		inputContainer,
		widget.NewSeparator(),
		controlContainer,
	)
	
	ag.window.SetContent(content)
}

func (ag *ASCIIGui) onSearch() {
	keyword := ag.searchEntry.Text
	if keyword == "" {
		dialog.ShowInformation("Warning", "Please enter search keyword first!", ag.window)
		return
	}
	
	// Simple search
	lines := strings.Split(ag.tikuContent, "\n")
	var results []string
	resultCount := 0
	
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
			results = append(results, fmt.Sprintf("Line %d: %s", i+1, line))
			resultCount++
			if resultCount >= 50 { // Limit results
				results = append(results, "... (showing first 50 results)")
				break
			}
		}
	}
	
	if len(results) == 0 {
		ag.resultText.SetText("No results found for: " + keyword)
	} else {
		resultHeader := fmt.Sprintf("SEARCH RESULTS for '%s' (%d matches):\n\n", keyword, resultCount)
		ag.resultText.SetText(resultHeader + strings.Join(results, "\n"))
	}
}

func (ag *ASCIIGui) onClear() {
	ag.searchEntry.SetText("")
	ag.resultText.SetText(ag.tikuContent)
}

func (ag *ASCIIGui) onAutoInput() {
	content := ag.inputEntry.Text
	if content == "" {
		dialog.ShowInformation("Warning", "Please enter content to copy first!", ag.window)
		return
	}
	
	// Copy to clipboard
	ag.window.Clipboard().SetContent(content)
	
	dialog.ShowInformation("Success", 
		"Content copied to clipboard successfully!\n\n"+
		"You can now paste it anywhere using Ctrl+V\n\n"+
		"Content: "+content, ag.window)
}

func (ag *ASCIIGui) showWelcome() {
	welcomeMsg := "Welcome to ChaoxingHelper ASCII Edition!\n\n" +
		"This version uses pure English interface to avoid\n" +
		"character encoding issues on different systems.\n\n" +
		"Features:\n" +
		"- Search question database\n" +
		"- View all content\n" +
		"- Copy answers to clipboard\n" +
		"- Compatible with all Windows versions\n\n" +
		"Instructions:\n" +
		"1. Enter keywords in search box\n" +
		"2. Click SEARCH to find matches\n" +
		"3. Use COPY TO CLIPBOARD for auto-input\n" +
		"4. Paste using Ctrl+V in target window\n\n" +
		"Click OK to continue..."

	dialog.ShowInformation("Welcome", welcomeMsg, ag.window)
}

func (ag *ASCIIGui) showAbout() {
	aboutMsg := "ChaoxingHelper ASCII Edition v1.0\n\n" +
		"A pure English interface version designed to solve\n" +
		"character encoding issues on different systems.\n\n" +
		"This program helps with:\n" +
		"- Question database searching\n" +
		"- Content viewing and browsing\n" +
		"- Clipboard operations\n\n" +
		"Compatible with:\n" +
		"- All Windows versions\n" +
		"- Different language settings\n" +
		"- Various encoding configurations\n\n" +
		"Build: 2025-06-26"

	dialog.ShowInformation("About", aboutMsg, ag.window)
}

func (ag *ASCIIGui) showHelp() {
	helpMsg := "HELP - How to use this program:\n\n" +
		"SEARCH:\n" +
		"1. Type keywords in the search box\n" +
		"2. Click SEARCH button\n" +
		"3. Results will show in the text area\n" +
		"4. Click CLEAR to reset\n\n" +
		"AUTO INPUT:\n" +
		"1. Copy desired text to input box\n" +
		"2. Click COPY TO CLIPBOARD\n" +
		"3. Go to target window/application\n" +
		"4. Press Ctrl+V to paste\n\n" +
		"TIPS:\n" +
		"- Search is case-insensitive\n" +
		"- Use simple keywords for better results\n" +
		"- Results are limited to 50 matches\n" +
		"- All text operations use clipboard\n\n" +
		"If you encounter any issues, try:\n" +
		"- Restart the program\n" +
		"- Check if tiku.txt file exists\n" +
		"- Run as administrator"

	dialog.ShowInformation("Help", helpMsg, ag.window)
}

// protectCurrentWindow applies screen recording protection to current window
func (ag *ASCIIGui) protectCurrentWindow() {
	// Get current process window handle
	user32 := syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow := user32.NewProc("GetForegroundWindow")
	procSetWindowDisplayAffinity := user32.NewProc("SetWindowDisplayAffinity")
	
	// Get current window handle
	hwnd, _, _ := procGetForegroundWindow.Call()
	
	if hwnd != 0 {
		// Apply WDA_EXCLUDEFROMCAPTURE to prevent screen recording
		dwAffinity := uintptr(0x00000011)
		ret, _, _ := procSetWindowDisplayAffinity.Call(hwnd, dwAffinity)
		
		if ret != 0 {
			fmt.Println("Screen protection applied to main window")
		} else {
			fmt.Println("Failed to apply screen protection")
		}
	}
} 