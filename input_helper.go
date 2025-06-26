package main

import (
	"syscall"
	"unsafe"
	"time"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procSendInput    = user32.NewProc("SendInput")
	procGetKeyState  = user32.NewProc("GetAsyncKeyState")
)

const (
	INPUT_KEYBOARD = 1
	KEYEVENTF_KEYUP = 0x0002
	KEYEVENTF_UNICODE = 0x0004
)

type INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
}

type KEYBDINPUT struct {
	VK        uint16
	Scan      uint16
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

// InputHelper Windows输入辅助
type InputHelper struct{}

// NewInputHelper 创建输入辅助
func NewInputHelper() *InputHelper {
	return &InputHelper{}
}

// TypeText 模拟键盘输入文本
func (ih *InputHelper) TypeText(text string) error {
	// 短暂延迟确保窗口获得焦点
	time.Sleep(100 * time.Millisecond)
	
	for _, char := range text {
		if err := ih.typeChar(char); err != nil {
			return err
		}
		// 每个字符间的短暂延迟
		time.Sleep(10 * time.Millisecond)
	}
	
	return nil
}

// typeChar 输入单个字符
func (ih *InputHelper) typeChar(char rune) error {
	// 按下
	input := INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			VK:    0,
			Scan:  uint16(char),
			Flags: KEYEVENTF_UNICODE,
			Time:  0,
		},
	}
	
	ret, _, _ := procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
	if ret == 0 {
		// 如果Unicode输入失败，尝试其他方法
		return nil
	}
	
	// 释放
	input.Ki.Flags = KEYEVENTF_UNICODE | KEYEVENTF_KEYUP
	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
	
	return nil
}

// SimulateKeyPress 模拟按键
func (ih *InputHelper) SimulateKeyPress(vk uint16) error {
	// 按下
	input := INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			VK:    vk,
			Scan:  0,
			Flags: 0,
			Time:  0,
		},
	}
	
	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
	
	// 释放
	input.Ki.Flags = KEYEVENTF_KEYUP
	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
	
	return nil
} 