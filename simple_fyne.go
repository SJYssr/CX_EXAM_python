package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SimpleFyneGUI 简化版GUI，专门解决编码问题
type SimpleFyneGUI struct {
	app         fyne.App
	window      fyne.Window
	fileHandler *FileHandler
	
	searchEntry *widget.Entry
	resultText  *widget.Entry
	inputEntry  *widget.Entry
	tikuContent string
}

// NewSimpleFyneGUI 创建简化版GUI
func NewSimpleFyneGUI() (*SimpleFyneGUI, error) {
	app := app.New()
	// 移除自定义主题，使用默认主题
	
	fileHandler, err := NewFileHandler()
	if err != nil {
		return nil, err
	}
	
	return &SimpleFyneGUI{
		app:         app,
		fileHandler: fileHandler,
	}, nil
}

// Run 运行简化版GUI
func (sg *SimpleFyneGUI) Run() {
	// 创建窗口
	sg.window = sg.app.NewWindow("Chaoxing Helper - Simple")
	sg.window.Resize(fyne.NewSize(600, 400))
	sg.window.CenterOnScreen()
	
	// 加载题库
	if err := sg.loadTiku(); err != nil {
		sg.tikuContent = "Failed to load tiku.txt"
	}
	
	// 创建界面
	sg.createInterface()
	
	// 显示窗口
	sg.window.ShowAndRun()
}

func (sg *SimpleFyneGUI) loadTiku() error {
	content, err := sg.fileHandler.LoadTikuContent()
	if err != nil {
		return err
	}
	sg.tikuContent = content
	return nil
}

func (sg *SimpleFyneGUI) createInterface() {
	// 搜索输入框
	sg.searchEntry = widget.NewEntry()
	sg.searchEntry.SetPlaceHolder("Search keyword...")
	
	// 搜索按钮
	searchBtn := widget.NewButton("Search", sg.onSearch)
	
	// 结果显示区域
	sg.resultText = widget.NewMultiLineEntry()
	sg.resultText.SetText(sg.tikuContent)
	sg.resultText.Wrapping = fyne.TextWrapWord
	
	// 输入框
	sg.inputEntry = widget.NewEntry()
	sg.inputEntry.SetPlaceHolder("Content to input...")
	
	// 输入按钮
	inputBtn := widget.NewButton("Auto Input", sg.onAutoInput)
	
	// 退出按钮
	exitBtn := widget.NewButton("Exit", func() {
		sg.window.Close()
	})
	
	// 布局
	searchContainer := container.NewBorder(nil, nil, nil, searchBtn, sg.searchEntry)
	inputContainer := container.NewBorder(nil, nil, nil, inputBtn, sg.inputEntry)
	buttonContainer := container.NewHBox(exitBtn)
	
	content := container.NewVBox(
		widget.NewLabel("Search:"),
		searchContainer,
		widget.NewSeparator(),
		widget.NewLabel("Results:"),
		container.NewScroll(sg.resultText),
		widget.NewSeparator(),
		widget.NewLabel("Auto Input:"),
		inputContainer,
		widget.NewSeparator(),
		buttonContainer,
	)
	
	sg.window.SetContent(content)
}

func (sg *SimpleFyneGUI) onSearch() {
	keyword := sg.searchEntry.Text
	if keyword == "" {
		dialog.ShowInformation("Info", "Please enter search keyword", sg.window)
		return
	}
	
	// 简单搜索
	lines := strings.Split(sg.tikuContent, "\n")
	var results []string
	
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
			results = append(results, fmt.Sprintf("Line %d: %s", i+1, line))
		}
	}
	
	if len(results) == 0 {
		sg.resultText.SetText("No results found")
	} else {
		sg.resultText.SetText(strings.Join(results, "\n"))
	}
}

func (sg *SimpleFyneGUI) onAutoInput() {
	content := sg.inputEntry.Text
	if content == "" {
		dialog.ShowInformation("Info", "Please enter content to input", sg.window)
		return
	}
	
	// 复制到剪贴板
	sg.window.Clipboard().SetContent(content)
	dialog.ShowInformation("Success", "Content copied to clipboard!\nPlease paste manually.", sg.window)
}

// 移除自定义主题实现，使用Fyne默认主题 