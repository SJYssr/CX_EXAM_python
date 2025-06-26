package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func minimalTest() {
	var mw *walk.MainWindow
	
	err := MainWindow{
		AssignTo: &mw,
		Title:    "最简测试",
		Size:     Size{Width: 300, Height: 200},
		Children: []Widget{
			Label{
				Text: "如果你看到这个窗口，说明GUI库工作正常！",
			},
		},
	}.Create()
	
	if err != nil {
		panic(err)
	}
	
	mw.Run()
} 