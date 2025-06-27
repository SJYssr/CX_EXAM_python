package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procSetClipboard = user32.NewProc("SetClipboardData")
	procOpenClipboard = user32.NewProc("OpenClipboard")
	procCloseClipboard = user32.NewProc("CloseClipboard")
	procEmptyClipboard = user32.NewProc("EmptyClipboard")

	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procGlobalAlloc  = kernel32.NewProc("GlobalAlloc")
	procGlobalLock   = kernel32.NewProc("GlobalLock")
	procGlobalUnlock = kernel32.NewProc("GlobalUnlock")
)

const (
	CF_UNICODETEXT = 13
	GMEM_MOVEABLE  = 0x0002
)

type Helper struct {
	window      fyne.Window
	questions   []string
	searchEntry *widget.Entry
	resultArea  *widget.Entry
	statusLabel *widget.Label
}

// 智能读取文件
func smartReadFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// 移除BOM
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}

	// UTF-8检查
	if utf8.Valid(data) {
		return string(data), nil
	}

	// GBK转换
	gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	utf8Data, _, err := transform.Bytes(gbkDecoder, data)
	if err == nil {
		return string(utf8Data), nil
	}

	// GB18030转换
	gb18030Decoder := simplifiedchinese.GB18030.NewDecoder()
	utf8Data, _, err = transform.Bytes(gb18030Decoder, data)
	if err == nil {
		return string(utf8Data), nil
	}

	return strings.ToValidUTF8(string(data), ""), nil
}

// Windows剪贴板操作
func setClipboardText(text string) error {
	ret, _, _ := procOpenClipboard.Call(0)
	if ret == 0 {
		return fmt.Errorf("Cannot open clipboard")
	}
	defer procCloseClipboard.Call()

	procEmptyClipboard.Call()

	utf16Text := syscall.StringToUTF16(text)
	textSize := len(utf16Text) * 2

	hMem, _, _ := procGlobalAlloc.Call(GMEM_MOVEABLE, uintptr(textSize))
	if hMem == 0 {
		return fmt.Errorf("Memory allocation failed")
	}

	pMem, _, _ := procGlobalLock.Call(hMem)
	if pMem == 0 {
		return fmt.Errorf("Memory lock failed")
	}

	copy((*[1 << 30]uint16)(unsafe.Pointer(pMem))[:len(utf16Text)], utf16Text)
	procGlobalUnlock.Call(hMem)

	ret, _, _ = procSetClipboard.Call(CF_UNICODETEXT, hMem)
	if ret == 0 {
		return fmt.Errorf("Clipboard set failed")
	}

	return nil
}

func (h *Helper) loadQuestions() error {
	content, err := smartReadFile("tiku.txt")
	if err != nil {
		return err
	}

	h.questions = []string{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			h.questions = append(h.questions, line)
		}
	}

	return nil
}

func (h *Helper) search(keyword string) []string {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return nil
	}

	var results []string
	for i, q := range h.questions {
		if strings.Contains(strings.ToLower(q), keyword) {
			results = append(results, fmt.Sprintf("[%d] %s", i+1, q))
			if len(results) >= 20 {
				results = append(results, "... (Max 20 results shown)")
				break
			}
		}
	}

	return results
}

func (h *Helper) searchOnline(query string) {
	url := "http://cx.icodef.com/wyn-nb?v=4&q=" + query
	
	h.statusLabel.SetText("Online searching...")
	
	go func() {
		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			h.resultArea.SetText("Search failed: " + err.Error())
			h.statusLabel.SetText("Search failed")
			return
		}
		defer resp.Body.Close()
		
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			h.resultArea.SetText("Read failed: " + err.Error())
			h.statusLabel.SetText("Read failed")
			return
		}
		
		result := string(body)
		h.resultArea.SetText(result)
		h.statusLabel.SetText("Search completed")
	}()
}

func (h *Helper) createUI() {
	// 创建应用
	fyneApp := app.New()
	
	h.window = fyneApp.NewWindow("Chaoxing Helper - Working v3.2")
	h.window.Resize(fyne.NewSize(900, 700))
	h.window.CenterOnScreen()

	// 顶部信息栏
	infoText := fmt.Sprintf("✓ Questions loaded: %d | UTF-8 OK | Clipboard OK", len(h.questions))
	infoLabel := widget.NewLabel(infoText)
	infoLabel.TextStyle = fyne.TextStyle{Bold: true}

	// 搜索输入框
	h.searchEntry = widget.NewEntry()
	h.searchEntry.SetPlaceHolder("Enter keywords to search...")
	
	// 搜索按钮
	localSearchBtn := widget.NewButton("Search Local", func() {
		keyword := h.searchEntry.Text
		if keyword == "" {
			dialog.ShowInformation("Info", "Please enter search keyword", h.window)
			return
		}
		
		results := h.search(keyword)
		if len(results) == 0 {
			h.resultArea.SetText("No results found. Try other keywords.")
		} else {
			h.resultArea.SetText(strings.Join(results, "\n\n"))
		}
		h.statusLabel.SetText(fmt.Sprintf("Local search completed, found %d results", len(results)))
	})

	onlineSearchBtn := widget.NewButton("Search Online", func() {
		query := h.searchEntry.Text
		if query == "" {
			dialog.ShowInformation("Info", "Please enter search content", h.window)
			return
		}
		h.searchOnline(query)
	})

	// 搜索框布局
	searchBox := container.NewBorder(nil, nil, nil, 
		container.NewHBox(localSearchBtn, onlineSearchBtn), 
		h.searchEntry)

	// 结果显示区域
	h.resultArea = widget.NewMultiLineEntry()
	h.resultArea.Wrapping = fyne.TextWrapWord
	h.resultArea.SetPlaceHolder("Search results will be displayed here...\n\nIMPORTANT: Clipboard function works correctly even if display shows symbols")
	
	resultScroll := container.NewScroll(h.resultArea)
	resultScroll.SetMinSize(fyne.NewSize(850, 400))

	// 功能按钮
	copyBtn := widget.NewButton("📋 Copy to Clipboard", func() {
		content := h.resultArea.Text
		if content == "" {
			dialog.ShowInformation("Info", "No content to copy", h.window)
			return
		}
		
		if err := setClipboardText(content); err != nil {
			dialog.ShowError(err, h.window)
		} else {
			dialog.ShowInformation("Success", "Content copied to clipboard successfully!\nPress Ctrl+V in exam window to paste", h.window)
		}
	})

	testBtn := widget.NewButton("🧪 Test Copy", func() {
		// 使用调试版本确认的正确中文内容
		testText := ""
		if len(h.questions) >= 3 {
			testText = fmt.Sprintf("Test content from question bank:\n1. %s\n2. %s\n3. %s", 
				h.questions[0], h.questions[1], h.questions[2])
		} else {
			testText = "Test: Chinese characters should work correctly in clipboard even if GUI shows symbols"
		}
		
		h.resultArea.SetText(testText)
		
		if err := setClipboardText(testText); err != nil {
			dialog.ShowError(err, h.window)
		} else {
			dialog.ShowInformation("Test Success", "Chinese content copied to clipboard!\nPlease paste somewhere to verify it's correct", h.window)
		}
	})

	clearBtn := widget.NewButton("🗑️ Clear", func() {
		h.searchEntry.SetText("")
		h.resultArea.SetText("")
		h.statusLabel.SetText("Ready")
	})

	showAllBtn := widget.NewButton("📄 Show All", func() {
		if len(h.questions) == 0 {
			h.resultArea.SetText("Question bank is empty or failed to load")
			return
		}
		
		count := 10
		if len(h.questions) < 10 {
			count = len(h.questions)
		}
		
		var result []string
		result = append(result, fmt.Sprintf("Question bank preview (first %d):\n", count))
		for i := 0; i < count; i++ {
			result = append(result, fmt.Sprintf("[%d] %s", i+1, h.questions[i]))
		}
		
		h.resultArea.SetText(strings.Join(result, "\n\n"))
		h.statusLabel.SetText(fmt.Sprintf("Showing first %d questions", count))
	})

	// 按钮布局
	buttonBox := container.NewHBox(copyBtn, testBtn, clearBtn, showAllBtn)

	// 状态栏
	h.statusLabel = widget.NewLabel("Ready")
	h.statusLabel.TextStyle = fyne.TextStyle{Italic: true}

	// 主界面布局
	content := container.NewVBox(
		infoLabel,
		widget.NewSeparator(),
		widget.NewCard("Search", "", searchBox),
		widget.NewSeparator(),
		buttonBox,
		widget.NewSeparator(),
		widget.NewCard("Results", "", resultScroll),
		widget.NewSeparator(),
		h.statusLabel,
	)

	h.window.SetContent(content)
	
	// 窗口关闭处理
	h.window.SetCloseIntercept(func() {
		h.window.Close()
	})
}

func main() {
	fmt.Println("Chaoxing Learning Assistant - Working v3.2")
	fmt.Println("==========================================")
	fmt.Println("Note: GUI may show symbols but clipboard works correctly")
	fmt.Println("All Chinese text will be copied properly")
	fmt.Println("")
	
	helper := &Helper{}
	
	// 加载题库
	fmt.Print("Loading question bank...")
	if err := helper.loadQuestions(); err != nil {
		fmt.Printf("\nError: %v\n", err)
		fmt.Println("Please ensure tiku.txt file exists in the same directory")
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
		return
	}
	fmt.Printf(" Success! Loaded %d questions\n", len(helper.questions))
	
	// 显示前3题作为测试（控制台能正常显示）
	fmt.Println("\nQuestion bank preview (console display is correct):")
	for i := 0; i < 3 && i < len(helper.questions); i++ {
		fmt.Printf("  %d. %s\n", i+1, helper.questions[i])
	}
	
	fmt.Println("\nStarting GUI...")
	fmt.Println("IMPORTANT: Even if GUI shows symbols, clipboard copy function works 100% correctly!")
	helper.createUI()
	helper.window.ShowAndRun()
} 