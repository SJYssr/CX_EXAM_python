package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// SimpleGUI 极简GUI结构
type SimpleGUI struct {
	mainWindow  *walk.MainWindow
	searchEdit  *walk.LineEdit
	contentText *walk.TextEdit
	inputEdit   *walk.LineEdit
}

// NewSimpleGUI 创建极简GUI
func NewSimpleGUI() (*SimpleGUI, error) {
	return &SimpleGUI{}, nil
}

// CreateWindow 创建简单窗口
func (sg *SimpleGUI) CreateWindow() error {
	log.Println("创建极简窗口...")
	
	err := MainWindow{
		AssignTo: &sg.mainWindow,
		Title:    "超星学习通辅助工具 - 简化版",
		Size:     Size{Width: 400, Height: 300},
		Layout:   VBox{},
		Children: []Widget{
			LineEdit{
				AssignTo: &sg.searchEdit,
			},
			PushButton{
				Text: "测试",
				OnClicked: func() {
					walk.MsgBox(sg.mainWindow, "测试", "GUI工作正常！", walk.MsgBoxOK)
				},
			},
			TextEdit{
				AssignTo: &sg.contentText,
				Text:     "这是一个测试文本框\n程序运行正常！\n\n快捷键：\nESC - 退出",
			},
			LineEdit{
				AssignTo: &sg.inputEdit,
			},
		},
		OnKeyDown: func(key walk.Key) {
			if key == walk.KeyEscape {
				sg.mainWindow.Close()
			}
		},
	}.Create()
	
	if err != nil {
		return err
	}
	
	log.Println("窗口创建成功")
	return nil
}

// Run 运行GUI
func (sg *SimpleGUI) Run() {
	sg.mainWindow.Run()
} 