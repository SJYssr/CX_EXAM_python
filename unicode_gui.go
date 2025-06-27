// +build unicode

package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"syscall"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// UnicodeGui with enhanced Chinese text support
type UnicodeGui struct {
	app         fyne.App
	window      fyne.Window
	protection  *ScreenProtection
	
	searchEntry *widget.Entry
	resultText  *widget.Entry
	inputEntry  *widget.Entry
	tikuContent string
	rawContent  string
}

// NewUnicodeGui creates Unicode-optimized GUI
func NewUnicodeGui() (*UnicodeGui, error) {
	app := app.New()
	
	// Initialize screen protection
	protection := NewScreenProtection()
	
	return &UnicodeGui{
		app:        app,
		protection: protection,
	}, nil
}

// Run the Unicode GUI
func (ug *UnicodeGui) Run() {
	// Create window
	ug.window = ug.app.NewWindow("ChaoxingHelper Unicode Edition - Chinese Fixed")
	ug.window.Resize(fyne.NewSize(800, 600))
	ug.window.CenterOnScreen()
	
	// Start screen protection
	if ug.protection != nil {
		go ug.protection.StartProtection()
		fmt.Println("Screen protection started")
	}
	
	// Load content with proper UTF-8 handling
	if err := ug.loadUnicodeContent(); err != nil {
		ug.tikuContent = "Error loading content: " + err.Error()
		ug.rawContent = ug.tikuContent
	}
	
	// Create interface
	ug.createInterface()
	
	// Apply screen protection
	ug.protectCurrentWindow()
	
	// Show welcome
	ug.showWelcome()
	
	// Show window
	ug.window.ShowAndRun()
}

func (ug *UnicodeGui) loadUnicodeContent() error {
	// Read file with enhanced UTF-8 processing
	content, err := ioutil.ReadFile("tiku.txt")
	if err != nil {
		return err
	}
	
	// Advanced UTF-8 processing and validation
	contentStr := string(content)
	
	// Clean up any invalid UTF-8 sequences
	contentStr = strings.ToValidUTF8(contentStr, "")
	
	// Remove any BOM markers
	contentStr = strings.TrimPrefix(contentStr, "\ufeff")
	
	// Normalize line endings
	contentStr = strings.ReplaceAll(contentStr, "\r\n", "\n")
	contentStr = strings.ReplaceAll(contentStr, "\r", "\n")
	
	// Store both raw and processed content
	ug.rawContent = contentStr
	
	// Count valid questions
	lines := strings.Split(contentStr, "\n")
	validLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			validLines++
		}
	}
	
	// Create enhanced display version
	header := "超星学习通题库 - Unicode Enhanced\n"
	header += "================================\n"
	header += fmt.Sprintf("文件大小: %d 字节\n", len(content))
	header += fmt.Sprintf("有效题目: %d 条\n", validLines)
	header += fmt.Sprintf("字符总数: %d\n", utf8.RuneCountInString(contentStr))
	header += fmt.Sprintf("编码状态: %s\n", func() string {
		if utf8.ValidString(contentStr) {
			return "UTF-8 ✓"
		}
		return "UTF-8 已修复"
	}())
	header += "中文支持: 完全兼容 ✓\n"
	header += "================================\n\n"
	
	ug.tikuContent = header + contentStr
	
	fmt.Printf("题库加载完成: %d 条题目\n", validLines)
	return nil
}

func (ug *UnicodeGui) createInterface() {
	// Search section
	searchLabel := widget.NewLabel("SEARCH (supports Chinese keywords):")
	ug.searchEntry = widget.NewEntry()
	ug.searchEntry.SetPlaceHolder("Enter Chinese or English keywords...")
	
	searchBtn := widget.NewButton("SEARCH", ug.onSearch)
	clearBtn := widget.NewButton("CLEAR", ug.onClear)
	showRawBtn := widget.NewButton("SHOW RAW", ug.showRawContent)
	
	searchContainer := container.NewHBox(
		ug.searchEntry,
		searchBtn,
		clearBtn,
		showRawBtn,
	)
	
	// Results section
	resultLabel := widget.NewLabel("SEARCH RESULTS / CONTENT (Chinese text should display correctly):")
	ug.resultText = widget.NewMultiLineEntry()
	ug.resultText.SetText(ug.tikuContent)
	ug.resultText.Wrapping = fyne.TextWrapWord
	
	resultScroll := container.NewScroll(ug.resultText)
	resultScroll.SetMinSize(fyne.NewSize(750, 300))
	
	// Input section
	inputLabel := widget.NewLabel("AUTO INPUT (Chinese text will be preserved in clipboard):")
	ug.inputEntry = widget.NewEntry()
	ug.inputEntry.SetPlaceHolder("Content to copy (supports Chinese)...")
	
	inputBtn := widget.NewButton("COPY TO CLIPBOARD", ug.onAutoInput)
	testBtn := widget.NewButton("TEST CHINESE", ug.testChinese)
	
	inputContainer := container.NewHBox(
		ug.inputEntry,
		inputBtn,
		testBtn,
	)
	
	// Control buttons
	aboutBtn := widget.NewButton("ABOUT", ug.showAbout)
	helpBtn := widget.NewButton("HELP", ug.showHelp)
	exitBtn := widget.NewButton("EXIT", func() {
		ug.window.Close()
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
	
	ug.window.SetContent(content)
}

func (ug *UnicodeGui) onSearch() {
	keyword := ug.searchEntry.Text
	if keyword == "" {
		dialog.ShowInformation("Warning", "Please enter search keyword!", ug.window)
		return
	}
	
	// Search in raw content
	lines := strings.Split(ug.rawContent, "\n")
	var results []string
	resultCount := 0
	
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
			results = append(results, fmt.Sprintf("Line %d: %s", i+1, line))
			resultCount++
			if resultCount >= 50 {
				results = append(results, "... (showing first 50 results)")
				break
			}
		}
	}
	
	if len(results) == 0 {
		ug.resultText.SetText("No results found for: " + keyword)
	} else {
		resultHeader := fmt.Sprintf("SEARCH RESULTS for '%s' (%d matches):\n\n", keyword, resultCount)
		ug.resultText.SetText(resultHeader + strings.Join(results, "\n"))
	}
}

func (ug *UnicodeGui) onClear() {
	ug.searchEntry.SetText("")
	ug.resultText.SetText(ug.tikuContent)
}

func (ug *UnicodeGui) showRawContent() {
	// Show the raw content without headers
	ug.resultText.SetText(ug.rawContent)
}

func (ug *UnicodeGui) testChinese() {
	// Enhanced Chinese test with common exam keywords
	testText := "测试完成 ✓ 中文字符正常显示：数据库、算法、程序设计、软件工程、计算机网络、操作系统"
	ug.inputEntry.SetText(testText)
	
	// Also copy to clipboard immediately for testing
	ug.window.Clipboard().SetContent(testText)
	
	dialog.ShowInformation("中文测试", 
		"测试文本已加载到输入框并复制到剪贴板！\n\n"+
		"测试内容: "+testText+"\n\n"+
		"请按 Ctrl+V 粘贴到考试窗口验证中文显示效果。\n"+
		"即使界面显示乱码，剪贴板中的中文依然正确！", 
		ug.window)
}

func (ug *UnicodeGui) onAutoInput() {
	content := ug.inputEntry.Text
	if content == "" {
		dialog.ShowInformation("提示", "请在输入框中输入要复制的内容！", ug.window)
		return
	}
	
	// Enhanced UTF-8 clipboard processing
	cleanContent := strings.ToValidUTF8(content, "")
	
	// Copy to clipboard with UTF-8 preservation
	ug.window.Clipboard().SetContent(cleanContent)
	
	// Count Chinese characters for verification
	chineseCount := 0
	for _, r := range cleanContent {
		if r >= 0x4e00 && r <= 0x9fff {
			chineseCount++
		}
	}
	
	statusMsg := fmt.Sprintf("复制成功！\n\n内容: %s\n\n统计信息:\n- 总字符数: %d\n- 中文字符: %d\n- UTF-8编码: ✓\n\n剪贴板已更新，请按 Ctrl+V 粘贴到考试界面。\n即使此处显示乱码，剪贴板内容100%%正确！", 
		cleanContent, utf8.RuneCountInString(cleanContent), chineseCount)
	
	dialog.ShowInformation("复制完成", statusMsg, ug.window)
}

func (ug *UnicodeGui) showWelcome() {
	welcomeMsg := "Welcome to ChaoxingHelper Unicode Edition!\n\n" +
		"This version has enhanced Chinese text support:\n\n" +
		"Features:\n" +
		"- Advanced UTF-8 text handling\n" +
		"- Chinese keyword searching\n" +
		"- Screen recording protection\n" +
		"- Unicode-safe clipboard operations\n\n" +
		"If Chinese text still appears garbled in the interface,\n" +
		"don't worry - it will copy and paste correctly!\n\n" +
		"Click TEST CHINESE to verify functionality."

	dialog.ShowInformation("Welcome", welcomeMsg, ug.window)
}

func (ug *UnicodeGui) showAbout() {
	aboutMsg := "ChaoxingHelper Unicode Edition v2.0\n\n" +
		"Special features for Chinese text handling:\n\n" +
		"- UTF-8 validation and cleanup\n" +
		"- Unicode-safe text operations\n" +
		"- Screen recording protection\n" +
		"- Multi-encoding support\n\n" +
		"Even if Chinese characters appear as symbols\n" +
		"in this interface, they are preserved correctly\n" +
		"in clipboard operations and will paste properly\n" +
		"in the exam window.\n\n" +
		"Build: 2025-06-26 Unicode"

	dialog.ShowInformation("About", aboutMsg, ug.window)
}

func (ug *UnicodeGui) showHelp() {
	helpMsg := "HELP - Unicode Edition:\n\n" +
		"SEARCH:\n" +
		"- Enter Chinese or English keywords\n" +
		"- Search works on original UTF-8 content\n" +
		"- Click SHOW RAW to see unprocessed text\n\n" +
		"AUTO INPUT:\n" +
		"- Chinese text is preserved in clipboard\n" +
		"- Click TEST CHINESE to verify functionality\n" +
		"- Paste using Ctrl+V in exam window\n\n" +
		"TROUBLESHOOTING:\n" +
		"- If Chinese shows as symbols: this is display issue only\n" +
		"- Content will paste correctly regardless of display\n" +
		"- For better display: enable UTF-8 in Windows settings\n\n" +
		"IMPORTANT: Even garbled text copies correctly!"

	dialog.ShowInformation("Help", helpMsg, ug.window)
}

func (ug *UnicodeGui) protectCurrentWindow() {
	user32 := syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow := user32.NewProc("GetForegroundWindow")
	procSetWindowDisplayAffinity := user32.NewProc("SetWindowDisplayAffinity")
	
	hwnd, _, _ := procGetForegroundWindow.Call()
	
	if hwnd != 0 {
		dwAffinity := uintptr(0x00000011)
		ret, _, _ := procSetWindowDisplayAffinity.Call(hwnd, dwAffinity)
		
		if ret != 0 {
			fmt.Println("Screen protection applied successfully")
		}
	}
} 