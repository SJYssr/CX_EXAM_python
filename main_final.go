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

type ChaoxingHelper struct {
	window      fyne.Window
	questions   []string
	searchEntry *widget.Entry
	resultArea  *widget.Entry
	statusLabel *widget.Label
}

// 智能读取文件 - 自动检测编码
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

	// 最后的字符清理
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

	// 转换为UTF-16
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

func (helper *ChaoxingHelper) loadQuestions() error {
	content, err := smartReadFile("tiku.txt")
	if err != nil {
		return err
	}

	helper.questions = []string{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			helper.questions = append(helper.questions, line)
		}
	}

	return nil
}

func (helper *ChaoxingHelper) search(keyword string) []string {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return nil
	}

	var results []string
	for i, q := range helper.questions {
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

func (helper *ChaoxingHelper) searchOnline(query string) {
	url := "http://cx.icodef.com/wyn-nb?v=4&q=" + query
	
	helper.statusLabel.SetText("在线搜索中...")
	
	go func() {
		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			helper.resultArea.SetText("在线搜索失败: " + err.Error())
			helper.statusLabel.SetText("搜索失败")
			return
		}
		defer resp.Body.Close()
		
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			helper.resultArea.SetText("读取搜索结果失败: " + err.Error())
			helper.statusLabel.SetText("读取失败")
			return
		}
		
		result := string(body)
		helper.resultArea.SetText(result)
		helper.statusLabel.SetText("在线搜索完成")
	}()
}

func (helper *ChaoxingHelper) createUI() {
	fyneApp := app.New()
	helper.window = fyneApp.NewWindow("超星学习通辅助工具 - 中文修复版 v3.0")
	helper.window.Resize(fyne.NewSize(900, 700))
	helper.window.CenterOnScreen()

	// 顶部信息栏
	infoText := fmt.Sprintf("✅ 题库已加载: %d 条 | UTF-8编码支持 | 剪贴板100%%兼容", len(helper.questions))
	infoLabel := widget.NewLabel(infoText)
	infoLabel.TextStyle = fyne.TextStyle{Bold: true}

	// 搜索输入框
	helper.searchEntry = widget.NewEntry()
	helper.searchEntry.SetPlaceHolder("输入中文或英文关键词搜索...")
	
	// 搜索按钮
	localSearchBtn := widget.NewButton("搜索题库", func() {
		keyword := helper.searchEntry.Text
		if keyword == "" {
			dialog.ShowInformation("提示", "请输入搜索关键词", helper.window)
			return
		}
		
		results := helper.search(keyword)
		if len(results) == 0 {
			helper.resultArea.SetText("未找到相关题目，请尝试其他关键词")
		} else {
			helper.resultArea.SetText(strings.Join(results, "\n\n"))
		}
		helper.statusLabel.SetText(fmt.Sprintf("本地搜索完成，找到 %d 条结果", len(results)))
	})

	onlineSearchBtn := widget.NewButton("在线搜索", func() {
		query := helper.searchEntry.Text
		if query == "" {
			dialog.ShowInformation("提示", "请输入搜索内容", helper.window)
			return
		}
		helper.searchOnline(query)
	})

	// 搜索框布局
	searchBox := container.NewBorder(nil, nil, nil, 
		container.NewHBox(localSearchBtn, onlineSearchBtn), 
		helper.searchEntry)

	// 结果显示区域
	helper.resultArea = widget.NewMultiLineEntry()
	helper.resultArea.Wrapping = fyne.TextWrapWord
	helper.resultArea.SetPlaceHolder("搜索结果将在这里显示...\n\n提示：剪贴板功能100%正常")
	
	resultScroll := container.NewScroll(helper.resultArea)
	resultScroll.SetMinSize(fyne.NewSize(850, 400))

	// 功能按钮
	copyBtn := widget.NewButton("📋 复制到剪贴板", func() {
		content := helper.resultArea.Text
		if content == "" {
			dialog.ShowInformation("提示", "没有内容可复制", helper.window)
			return
		}
		
		if err := setClipboardText(content); err != nil {
			dialog.ShowError(err, helper.window)
		} else {
			dialog.ShowInformation("复制成功", "内容已复制到剪贴板！\n请在考试窗口按 Ctrl+V 粘贴", helper.window)
		}
	})

	testBtn := widget.NewButton("🧪 测试中文", func() {
		testText := "测试成功！中文显示正常：\n" +
			"数据库管理、算法设计、程序开发、软件工程、\n" +
			"计算机网络、操作系统、数据结构、人工智能"
		helper.resultArea.SetText(testText)
		
		if err := setClipboardText(testText); err != nil {
			dialog.ShowError(err, helper.window)
		} else {
			dialog.ShowInformation("测试成功", "中文测试文本已复制到剪贴板！\n请粘贴验证功能是否正常", helper.window)
		}
	})

	clearBtn := widget.NewButton("🗑️ 清空", func() {
		helper.searchEntry.SetText("")
		helper.resultArea.SetText("")
		helper.statusLabel.SetText("就绪")
	})

	showAllBtn := widget.NewButton("📄 显示题库", func() {
		if len(helper.questions) == 0 {
			helper.resultArea.SetText("题库为空或加载失败")
			return
		}
		
		count := 10
		if len(helper.questions) < 10 {
			count = len(helper.questions)
		}
		
		var result []string
		result = append(result, fmt.Sprintf("题库预览（前%d题）：\n", count))
		for i := 0; i < count; i++ {
			result = append(result, fmt.Sprintf("[%d] %s", i+1, helper.questions[i]))
		}
		
		helper.resultArea.SetText(strings.Join(result, "\n\n"))
		helper.statusLabel.SetText(fmt.Sprintf("显示前 %d 题", count))
	})

	// 按钮布局
	buttonBox1 := container.NewHBox(copyBtn, testBtn, clearBtn)
	buttonBox2 := container.NewHBox(showAllBtn)

	// 状态栏
	helper.statusLabel = widget.NewLabel("就绪")
	helper.statusLabel.TextStyle = fyne.TextStyle{Italic: true}

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
		helper.statusLabel,
	)

	helper.window.SetContent(content)
	
	// 窗口关闭处理
	helper.window.SetCloseIntercept(func() {
		helper.window.Close()
	})
}

func main() {
	fmt.Println("超星学习通辅助工具 - 中文修复版 v3.0")
	fmt.Println("==========================================")
	fmt.Println("🎯 完全解决中文乱码问题")
	fmt.Println("📚 智能编码检测和转换")
	fmt.Println("📋 Windows剪贴板100%支持")
	fmt.Println("")
	
	helper := &ChaoxingHelper{}
	
	// 加载题库
	fmt.Print("正在加载题库文件...")
	if err := helper.loadQuestions(); err != nil {
		fmt.Printf("\n❌ 错误: %v\n", err)
		fmt.Println("请确保 tiku.txt 文件存在于程序同一目录")
		fmt.Println("\n按回车键退出...")
		fmt.Scanln()
		return
	}
	fmt.Printf(" ✅ 成功！共加载 %d 条题目\n", len(helper.questions))
	
	// 显示前3题作为测试
	fmt.Println("\n题库预览（前3题）：")
	for i := 0; i < 3 && i < len(helper.questions); i++ {
		fmt.Printf("  %d. %s\n", i+1, helper.questions[i])
	}
	
	fmt.Println("\n正在启动GUI界面...")
	helper.createUI()
	helper.window.ShowAndRun()
} 