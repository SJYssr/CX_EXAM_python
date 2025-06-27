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
	user32                       = syscall.NewLazyDLL("user32.dll")
	procSetClipboard             = user32.NewProc("SetClipboardData")
	procOpenClipboard            = user32.NewProc("OpenClipboard")
	procCloseClipboard           = user32.NewProc("CloseClipboard")
	procEmptyClipboard           = user32.NewProc("EmptyClipboard")
	procSetWindowDisplayAffinity = user32.NewProc("SetWindowDisplayAffinity")

	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procGlobalAlloc  = kernel32.NewProc("GlobalAlloc")
	procGlobalLock   = kernel32.NewProc("GlobalLock")
	procGlobalUnlock = kernel32.NewProc("GlobalUnlock")
)

const (
	CF_UNICODETEXT         = 13
	GMEM_MOVEABLE          = 0x0002
	WDA_EXCLUDEFROMCAPTURE = 0x00000011
)

type ChaoxingApp struct {
	window      fyne.Window
	questions   []string
	searchEntry *widget.Entry
	resultArea  *widget.Entry
	statusLabel *widget.Label
}

// 智能文件读取 - 支持多种编码
func smartReadFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// 移除UTF-8 BOM
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}

	// 首先尝试UTF-8
	if utf8.Valid(data) {
		return string(data), nil
	}

	// 尝试GBK转UTF-8
	gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	utf8Data, _, err := transform.Bytes(gbkDecoder, data)
	if err == nil {
		return string(utf8Data), nil
	}

	// 尝试GB18030转UTF-8
	gb18030Decoder := simplifiedchinese.GB18030.NewDecoder()
	utf8Data, _, err = transform.Bytes(gb18030Decoder, data)
	if err == nil {
		return string(utf8Data), nil
	}

	// 最后清理无效字符
	return strings.ToValidUTF8(string(data), ""), nil
}

// Windows剪贴板操作 - 确保中文正确复制
func setClipboardText(text string) error {
	ret, _, _ := procOpenClipboard.Call(0)
	if ret == 0 {
		return fmt.Errorf("无法打开剪贴板")
	}
	defer procCloseClipboard.Call()

	procEmptyClipboard.Call()

	// 转换为UTF-16（Windows标准）
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

	// 安全复制数据
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
	
	app.statusLabel.SetText("在线搜索中...")
	
	go func() {
		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			app.resultArea.SetText("在线搜索失败: " + err.Error())
			app.statusLabel.SetText("搜索失败")
			return
		}
		defer resp.Body.Close()
		
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			app.resultArea.SetText("读取搜索结果失败: " + err.Error())
			app.statusLabel.SetText("读取失败")
			return
		}
		
		result := string(body)
		app.resultArea.SetText(result)
		app.statusLabel.SetText("在线搜索完成")
	}()
}

func (chaoxingApp *ChaoxingApp) createUI() {
	fyneApp := app.New()
	chaoxingApp.window = fyneApp.NewWindow("超星学习通辅助工具 - 中文修复版 v3.0")
	app.window.Resize(fyne.NewSize(900, 700))
	app.window.CenterOnScreen()

	// 顶部信息栏
	infoText := fmt.Sprintf("✅ 题库已加载: %d 条 | UTF-8编码支持 | 剪贴板100%%兼容", len(app.questions))
	infoLabel := widget.NewLabel(infoText)
	infoLabel.TextStyle = fyne.TextStyle{Bold: true}

	// 搜索输入框
	app.searchEntry = widget.NewEntry()
	app.searchEntry.SetPlaceHolder("输入中文或英文关键词搜索...")
	
	// 搜索按钮
	localSearchBtn := widget.NewButton("搜索题库", func() {
		keyword := app.searchEntry.Text
		if keyword == "" {
			dialog.ShowInformation("提示", "请输入搜索关键词", app.window)
			return
		}
		
		results := app.search(keyword)
		if len(results) == 0 {
			app.resultArea.SetText("未找到相关题目，请尝试其他关键词")
		} else {
			app.resultArea.SetText(strings.Join(results, "\n\n"))
		}
		app.statusLabel.SetText(fmt.Sprintf("本地搜索完成，找到 %d 条结果", len(results)))
	})

	onlineSearchBtn := widget.NewButton("在线搜索", func() {
		query := app.searchEntry.Text
		if query == "" {
			dialog.ShowInformation("提示", "请输入搜索内容", app.window)
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
	app.resultArea.SetPlaceHolder("搜索结果将在这里显示...\n\n提示：即使界面显示异常，复制功能100%正常")
	
	resultScroll := container.NewScroll(app.resultArea)
	resultScroll.SetMinSize(fyne.NewSize(850, 400))

	// 功能按钮
	copyBtn := widget.NewButton("📋 复制到剪贴板", func() {
		content := app.resultArea.Text
		if content == "" {
			dialog.ShowInformation("提示", "没有内容可复制", app.window)
			return
		}
		
		if err := setClipboardText(content); err != nil {
			dialog.ShowError(err, app.window)
		} else {
			dialog.ShowInformation("复制成功", "内容已复制到剪贴板！\n请在考试窗口按 Ctrl+V 粘贴", app.window)
		}
	})

	testBtn := widget.NewButton("🧪 测试中文", func() {
		testText := "测试成功！中文显示正常：\n" +
			"数据库管理、算法设计、程序开发、软件工程、\n" +
			"计算机网络、操作系统、数据结构、人工智能"
		app.resultArea.SetText(testText)
		
		if err := setClipboardText(testText); err != nil {
			dialog.ShowError(err, app.window)
		} else {
			dialog.ShowInformation("测试成功", "中文测试文本已复制到剪贴板！\n请粘贴验证功能是否正常", app.window)
		}
	})

	clearBtn := widget.NewButton("🗑️ 清空", func() {
		app.searchEntry.SetText("")
		app.resultArea.SetText("")
		app.statusLabel.SetText("就绪")
	})

	showAllBtn := widget.NewButton("📄 显示题库", func() {
		if len(app.questions) == 0 {
			app.resultArea.SetText("题库为空或加载失败")
			return
		}
		
		count := 10
		if len(app.questions) < 10 {
			count = len(app.questions)
		}
		
		var result []string
		result = append(result, fmt.Sprintf("题库预览（前%d题）：\n", count))
		for i := 0; i < count; i++ {
			result = append(result, fmt.Sprintf("[%d] %s", i+1, app.questions[i]))
		}
		
		app.resultArea.SetText(strings.Join(result, "\n\n"))
		app.statusLabel.SetText(fmt.Sprintf("显示前 %d 题", count))
	})

	// 按钮布局
	buttonBox1 := container.NewHBox(copyBtn, testBtn, clearBtn)
	buttonBox2 := container.NewHBox(showAllBtn)

	// 状态栏
	app.statusLabel = widget.NewLabel("就绪")
	app.statusLabel.TextStyle = fyne.TextStyle{Italic: true}

	// 主界面布局
	content := container.NewVBox(
		infoLabel,
		widget.NewSeparator(),
		widget.NewCard("搜索", "", searchBox),
		widget.NewSeparator(),
		buttonBox1,
		buttonBox2,
		widget.NewSeparator(),
		widget.NewCard("结果", "", resultScroll),
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
	// 强制设置UTF-8控制台（Windows专用）
	
	fmt.Println("超星学习通辅助工具 - 中文修复版 v3.0")
	fmt.Println("==========================================")
	fmt.Println("🎯 完全解决中文乱码问题")
	fmt.Println("📚 智能编码检测和转换")
	fmt.Println("📋 Windows剪贴板100%支持")
	fmt.Println("")
	
	app := &ChaoxingApp{}
	
	// 加载题库
	fmt.Print("正在加载题库文件...")
	if err := app.loadQuestions(); err != nil {
		fmt.Printf("\n❌ 错误: %v\n", err)
		fmt.Println("请确保 tiku.txt 文件存在于程序同一目录")
		fmt.Println("\n按回车键退出...")
		fmt.Scanln()
		return
	}
	fmt.Printf(" ✅ 成功！共加载 %d 条题目\n", len(app.questions))
	
	// 显示前3题作为测试
	fmt.Println("\n题库预览（前3题）：")
	for i := 0; i < 3 && i < len(app.questions); i++ {
		fmt.Printf("  %d. %s\n", i+1, app.questions[i])
	}
	
	fmt.Println("\n正在启动GUI界面...")
	app.createUI()
	app.window.ShowAndRun()
} 