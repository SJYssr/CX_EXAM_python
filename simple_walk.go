package main

import (
	"github.com/lxn/walk"
)

func simpleWalkTest() error {
	// 直接使用walk的API而不是声明式语法
	mainWindow, err := walk.NewMainWindow()
	if err != nil {
		return err
	}
	
	mainWindow.SetTitle("超星学习通辅助工具 - 简化版")
	mainWindow.SetSize(walk.Size{Width: 400, Height: 300})
	
	// 创建垂直布局
	layout := walk.NewVBoxLayout()
	layout.SetMargins(walk.Margins{10, 10, 10, 10})
	layout.SetSpacing(10)
	mainWindow.SetLayout(layout)
	
	// 添加标签
	label, err := walk.NewLabel(mainWindow)
	if err != nil {
		return err
	}
	label.SetText("程序运行正常！\n\n快捷键：ESC - 退出程序")
	
	// 添加按钮
	button, err := walk.NewPushButton(mainWindow)
	if err != nil {
		return err
	}
	button.SetText("测试按钮")
	button.Clicked().Attach(func() {
		walk.MsgBox(mainWindow, "测试", "按钮点击成功！", walk.MsgBoxOK)
	})
	
	// 绑定键盘事件
	mainWindow.KeyDown().Attach(func(key walk.Key) {
		if key == walk.KeyEscape {
			mainWindow.Close()
		}
	})
	
	// 运行窗口
	mainWindow.Run()
	return nil
} 