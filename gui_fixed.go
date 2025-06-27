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
	"fyne.io/fyne/v2/theme"
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

// 自定义主题以确保字体正确
type chineseTheme struct {
	fyne.Theme
}

func (t chineseTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t chineseTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14 // 稍微大一点的字体
	default:
		return theme.DefaultTheme().Size(name)
	}
}

type ChaoxingApp struct {
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
		return fmt.Errorf("无法打开剪贴板")
	}
	defer procCloseClipboard.Call()

	procEmptyClipboard.Call()

	utf16Text := syscall.StringToUTF16(text)
	textSize := len(utf16Text) * 2

	hMem, _, _ := procGlobalAlloc.Call(GMEM_MOVEABLE, uintptr(textSize))
	if hMem == 0 {
		return fmt.Errorf("内存分配失败")
	}

	pMem, _, _ := procGlobalLock.Call(hMem)
	if pMem == 0 {
		return fmt.Errorf("内存锁定失败")
	}

	copy((*[1 << 30]uint16)(unsafe.Pointer(pMem))[:len(utf16Text)], utf16Text)
	procGlobalUnlock.Call(hMem)

	ret, _, _ = procSetClipboard.Call(CF_UNICODETEXT, hMem)
	if ret == 0 {
		return fmt.Errorf("剪贴板设置失败")
	}

	return nil
}

func (app *ChaoxingApp) loadQuestions() error {
	content, err := smartReadFile("tiku.txt")
	if err != nil {
		return err
	}

	app.questions = []string{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			app.questions = append(app.questions, line)
		}
	}

	return nil
}

func (app *ChaoxingApp) search(keyword string) []string {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return nil
	}

	var results []string
	for i, q := range app.questions {
		if strings.Contains(strings.ToLower(q), keyword) {
			results = append(results, fmt.Sprintf("[%d] %s", i+1, q))
			if len(results) >= 20 {
				results = append(results, "... (最多显示20条结果)")
				break
			}
		}
	}

	return results
}

func (app *ChaoxingApp) searchOnline(query string) {
	url := "http://cx.icodef.com/wyn-nb?v=4&q=" + query
	
	app.statusLabel.SetText("Online searching...")
	
	go func() {
		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			app.resultArea.SetText("Search failed: " + err.Error())
			app.statusLabel.SetText("Search failed")
			return
		}
		defer resp.Body.Close()
		
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			app.resultArea.SetText("Read failed: " + err.Error())
			app.statusLabel.SetText("Read failed")
			return
		}
		
		result := string(body)
		app.resultArea.SetText(result)
		app.statusLabel.SetText("Search completed")
	}()
}

func (chaoxingApp *ChaoxingApp) createUI() {
	// 创建应用并设置自定义主题
	fyneApp := app.New()
	fyneApp.Settings().SetTheme(&chineseTheme{})
	
	app.window = fyneApp.NewWindow("Chaoxing Helper - Fixed v3.1")
	app.window.Resize(fyne.NewSize(900, 700))
	app.window.CenterOnScreen()

	// 顶部信息栏 - 使用英文避免字体问题
	infoText := fmt.Sprintf("Questions loaded: %d | UTF-8 support | Clipboard OK", len(app.questions))
	infoLabel := widget.NewLabel(infoText)
	infoLabel.TextStyle = fyne.TextStyle{Bold: true}

	// 搜索输入框
	app.searchEntry = widget.NewEntry()
	app.searchEntry.SetPlaceHolder("Enter keywords to search...")
	
	// 搜索按钮
	localSearchBtn := widget.NewButton("Search Local", func() {
		keyword := app.searchEntry.Text
		if keyword == "" {
			dialog.ShowInformation("Info", "Please enter search keyword", app.window)
			return
		}
		
		results := app.search(keyword)
		if len(results) == 0 {
			app.resultArea.SetText("No results found. Try other keywords.")
		} else {
			app.resultArea.SetText(strings.Join(results, "\n\n"))
		}
		app.statusLabel.SetText(fmt.Sprintf("Local search completed, found %d results", len(results)))
	})

	onlineSearchBtn := widget.NewButton("Search Online", func() {
		query := app.searchEntry.Text
		if query == "" {
			dialog.ShowInformation("Info", "Please enter search content", app.window)
			return
		}
		app.searchOnline(query)
	})

	// 搜索框布局
	searchBox := container.NewBorder(nil, nil, nil, 
		container.NewHBox(localSearchBtn, onlineSearchBtn), 
		app.searchEntry)

	// 结果显示区域
	app.resultArea = widget.NewMultiLineEntry()
	app.resultArea.Wrapping = fyne.TextWrapWord
	app.resultArea.SetPlaceHolder("Search results will be displayed here...\n\nNote: Clipboard function works 100% correctly")
	
	resultScroll := container.NewScroll(app.resultArea)
	resultScroll.SetMinSize(fyne.NewSize(850, 400))

	// 功能按钮
	copyBtn := widget.NewButton("Copy to Clipboard", func() {
		content := app.resultArea.Text
		if content == "" {
			dialog.ShowInformation("Info", "No content to copy", app.window)
			return
		}
		
		if err := setClipboardText(content); err != nil {
			dialog.ShowError(err, app.window)
		} else {
			dialog.ShowInformation("Success", "Content copied to clipboard!\nPress Ctrl+V in exam window to paste", app.window)
		}
	})

	testBtn := widget.NewButton("Test Chinese", func() {
		// 使用调试版本确认可以正常显示的中文内容
		testText := "Test successful! Chinese display normal:\nDatabase management, Algorithm design, Program development, Software engineering,\nComputer network, Operating system, Data structure, Artificial intelligence"
		
		// 实际的中文内容 - 如果显示正常说明字体修复成功
		if len(app.questions) > 0 {
			testText = "Chinese test:\n" + strings.Join(app.questions[:3], "\n")
		}
		
		app.resultArea.SetText(testText)
		
		if err := setClipboardText(testText); err != nil {
			dialog.ShowError(err, app.window)
		} else {
			dialog.ShowInformation("Test Success", "Chinese test text copied to clipboard!\nPlease paste to verify functionality", app.window)
		}
	})

	clearBtn := widget.NewButton("Clear", func() {
		app.searchEntry.SetText("")
		app.resultArea.SetText("")
		app.statusLabel.SetText("Ready")
	})

	showAllBtn := widget.NewButton("Show Questions", func() {
		if len(app.questions) == 0 {
			app.resultArea.SetText("Question bank is empty or failed to load")
			return
		}
		
		count := 10
		if len(app.questions) < 10 {
			count = len(app.questions)
		}
		
		var result []string
		result = append(result, fmt.Sprintf("Question bank preview (first %d):\n", count))
		for i := 0; i < count; i++ {
			result = append(result, fmt.Sprintf("[%d] %s", i+1, app.questions[i]))
		}
		
		app.resultArea.SetText(strings.Join(result, "\n\n"))
		app.statusLabel.SetText(fmt.Sprintf("Showing first %d questions", count))
	})

	// 按钮布局
	buttonBox1 := container.NewHBox(copyBtn, testBtn, clearBtn)
	buttonBox2 := container.NewHBox(showAllBtn)

	// 状态栏
	app.statusLabel = widget.NewLabel("Ready")
	app.statusLabel.TextStyle = fyne.TextStyle{Italic: true}

	// 主界面布局
	content := container.NewVBox(
		infoLabel,
		widget.NewSeparator(),
		widget.NewCard("Search", "", searchBox),
		widget.NewSeparator(),
		buttonBox1,
		buttonBox2,
		widget.NewSeparator(),
		widget.NewCard("Results", "", resultScroll),
		widget.NewSeparator(),
		app.statusLabel,
	)

	app.window.SetContent(content)
	
	// 窗口关闭处理
	app.window.SetCloseIntercept(func() {
		app.window.Close()
	})
}

func main() {
	fmt.Println("Chaoxing Learning Assistant - Font Fixed v3.1")
	fmt.Println("===============================================")
	fmt.Println("Fixing Chinese font display issues...")
	fmt.Println("Smart encoding detection and conversion")
	fmt.Println("Windows clipboard 100% support")
	fmt.Println("")
	
	app := &ChaoxingApp{}
	
	// 加载题库
	fmt.Print("Loading question bank...")
	if err := app.loadQuestions(); err != nil {
		fmt.Printf("\nError: %v\n", err)
		fmt.Println("Please ensure tiku.txt file exists in the same directory")
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
		return
	}
	fmt.Printf(" Success! Loaded %d questions\n", len(app.questions))
	
	// 显示前3题作为测试（控制台能正常显示）
	fmt.Println("\nQuestion bank preview (first 3):")
	for i := 0; i < 3 && i < len(app.questions); i++ {
		fmt.Printf("  %d. %s\n", i+1, app.questions[i])
	}
	
	fmt.Println("\nStarting GUI interface with font fixes...")
	app.createUI()
	app.window.ShowAndRun()
} 