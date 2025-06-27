package main

import (
	"bufio"
	"fmt"
	"image/color"
	"io/ioutil"
	"strings"
	"syscall"
	"unicode/utf8"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Windows剪贴板相关
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

// 自定义主题，强制使用系统字体
type chineseTheme struct {
	fyne.Theme
}

func (t *chineseTheme) Font(style fyne.TextStyle) fyne.Resource {
	// 返回nil让系统自动选择最佳中文字体
	return nil
}

func (t *chineseTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14 // 增大字体大小
	default:
		return theme.DefaultTheme().Size(name)
	}
}

func (t *chineseTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) (c color.Color) {
	return theme.DefaultTheme().Color(name, variant)
}

func (t *chineseTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
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

func main() {
	fmt.Println("超星学习通辅助工具 - GUI中文字体修复版")
	fmt.Println("正在启动GUI界面...")

	// 加载题库
	content, err := smartReadFile("tiku.txt")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		fmt.Println("请确保 tiku.txt 文件存在于程序同一目录")
		fmt.Println("按任意键退出...")
		fmt.Scanln()
		return
	}

	var questions []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			questions = append(questions, line)
		}
	}

	fmt.Printf("成功加载 %d 条题目\n", len(questions))

	// 创建Fyne应用
	myApp := app.New()
	
	// 应用自定义主题
	myApp.Settings().SetTheme(&chineseTheme{})
	
	myWindow := myApp.NewWindow("超星学习通辅助工具 - GUI修复版")
	myWindow.Resize(fyne.NewSize(800, 650))

	// 搜索功能
	searchFunc := func(keyword string) []string {
		keyword = strings.ToLower(strings.TrimSpace(keyword))
		if keyword == "" {
			return nil
		}

		var results []string
		for i, q := range questions {
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

	// 创建UI组件 - 使用简化的ASCII文本
	statusLabel := widget.NewLabel(fmt.Sprintf("题库已加载: %d 条 | UTF-8编码支持 | 剪贴板兼容", len(questions)))
	statusLabel.Importance = widget.MediumImportance
	
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("请输入中文或英文关键词搜索...")
	
	resultsText := widget.NewMultiLineEntry()
	resultsText.SetText("搜索结果将在这里显示...\n\n提示：\n1. 输入关键词进行搜索\n2. 支持中文和英文搜索\n3. 点击复制按钮将结果复制到剪贴板\n4. 在考试窗口按 Ctrl+V 粘贴")
	resultsText.Resize(fyne.NewSize(750, 350))

	// 按钮功能
	searchBtn := widget.NewButton("搜索", func() {
		keyword := searchEntry.Text
		if keyword == "" {
			resultsText.SetText("请输入搜索关键词")
			return
		}

		results := searchFunc(keyword)
		if len(results) > 0 {
			resultsText.SetText(strings.Join(results, "\n\n"))
		} else {
			resultsText.SetText("未找到相关题目，请尝试其他关键词")
		}
	})

	copyBtn := widget.NewButton("复制到剪贴板", func() {
		content := resultsText.Text
		if content == "" {
			resultsText.SetText("没有内容可复制")
			return
		}

		if err := setClipboardText(content); err != nil {
			resultsText.SetText("复制失败: " + err.Error())
		} else {
			// 在现有内容前添加成功提示
			resultsText.SetText("成功！内容已复制到剪贴板！请在考试窗口按 Ctrl+V 粘贴\n\n" + content)
		}
	})

	testBtn := widget.NewButton("测试中文", func() {
		testText := "测试成功！中文显示正常：\n数据库管理、算法设计、程序开发、软件工程\n计算机网络、操作系统、数据结构、人工智能"
		resultsText.SetText(testText)
		
		if err := setClipboardText(testText); err != nil {
			resultsText.SetText("复制失败: " + err.Error())
		} else {
			resultsText.SetText("成功！中文测试文本已复制到剪贴板！请粘贴验证功能\n\n" + testText)
		}
	})

	showBtn := widget.NewButton("显示题库", func() {
		count := 10
		if len(questions) < 10 {
			count = len(questions)
		}
		
		var preview []string
		for i := 0; i < count; i++ {
			preview = append(preview, fmt.Sprintf("[%d] %s", i+1, questions[i]))
		}
		
		resultsText.SetText(fmt.Sprintf("题库预览（前%d题）：\n\n%s", count, strings.Join(preview, "\n\n")))
	})

	clearBtn := widget.NewButton("清空", func() {
		searchEntry.SetText("")
		resultsText.SetText("搜索结果将在这里显示...\n\n提示：\n1. 输入关键词进行搜索\n2. 支持中文和英文搜索\n3. 点击复制按钮将结果复制到剪贴板\n4. 在考试窗口按 Ctrl+V 粘贴")
	})

	// 回车键搜索
	searchEntry.OnSubmitted = func(text string) {
		searchBtn.OnTapped()
	}

	// 布局
	searchBox := container.NewBorder(nil, nil, nil, searchBtn, searchEntry)
	buttonBox := container.NewHBox(copyBtn, testBtn, showBtn, clearBtn)
	
	mainContent := container.NewVBox(
		statusLabel,
		widget.NewSeparator(),
		widget.NewLabel("搜索功能:"),
		searchBox,
		widget.NewSeparator(),
		buttonBox,
		widget.NewSeparator(),
		widget.NewLabel("搜索结果:"),
		resultsText,
	)

	myWindow.SetContent(mainContent)
	
	fmt.Println("GUI界面已启动，应用了中文字体修复！")
	fmt.Println("如果仍有乱码，这是字体渲染问题，但剪贴板功能是正常的！")

	// 运行应用
	myWindow.ShowAndRun()
} 