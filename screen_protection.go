package main

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/lxn/walk"
	"github.com/lxn/win"
)

// ScreenProtection 屏幕保护结构体
type ScreenProtection struct {
	user32           *syscall.LazyDLL
	kernel32         *syscall.LazyDLL
	imeWindows       map[win.HWND]bool
	screenWidth      int32
	screenHeight     int32
	dynamicX         int32
	dynamicY         int32
	chaoxingWindow   win.HWND
	protectionActive bool
}

// WindowRect 窗口矩形结构
type WindowRect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

// NewScreenProtection 创建屏幕保护实例
func NewScreenProtection() (*ScreenProtection, error) {
	sp := &ScreenProtection{
		user32:     syscall.NewLazyDLL("user32.dll"),
		kernel32:   syscall.NewLazyDLL("kernel32.dll"),
		imeWindows: make(map[win.HWND]bool),
	}
	
	// 获取屏幕尺寸
	sp.screenWidth = win.GetSystemMetrics(win.SM_CXSCREEN)
	sp.screenHeight = win.GetSystemMetrics(win.SM_CYSCREEN)
	
	return sp, nil
}

// StartProtection 启动保护功能
func (sp *ScreenProtection) StartProtection() {
	sp.protectionActive = true
	go sp.protectionLoop()
}

// StopProtection 停止保护功能
func (sp *ScreenProtection) StopProtection() {
	sp.protectionActive = false
}

// protectionLoop 保护循环
func (sp *ScreenProtection) protectionLoop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	for sp.protectionActive {
		select {
		case <-ticker.C:
			// 检测并保护窗口
			sp.DetectAndProtectWindows()
			
			// 每10秒重新检测超星窗口位置
			sp.DetectDynamicPosition()
		}
	}
}

// DetectAndProtectWindows 检测并保护窗口
func (sp *ScreenProtection) DetectAndProtectWindows() {
	// 保护当前程序窗口
	foregroundWindow := win.GetForegroundWindow()
	sp.ProtectWindow(foregroundWindow)
	
	// 查找并保护输入法窗口
	sp.FindAndProtectIMEWindows()
}

// FindAndProtectIMEWindows 查找并保护输入法窗口
func (sp *ScreenProtection) FindAndProtectIMEWindows() {
	// 清空之前的输入法窗口列表
	sp.imeWindows = make(map[win.HWND]bool)
	
	// 输入法相关的窗口类名和标题关键词
	imeKeywords := []string{
		"搜狗输入法", "百度输入法", "QQ输入法", "谷歌拼音输入法",
		"微软拼音", "紫光拼音", "讯飞输入法", "Sogou", "Baidu",
		"InputMethod", "IME", "SOGOU_IME_WND", "SogouPY",
		"BaiduIME", "QQPinyin", "CiceroUIWndFrame",
		"MSCTFIME UI", "IMEUI", "SunPinyin", "SCIM",
		"GooglePinyinIME", "BaiduPY", "iFlyInput",
	}
	
	// 枚举所有窗口
	win.EnumWindows(func(hwnd win.HWND) bool {
		if !win.IsWindowVisible(hwnd) {
			return true
		}
		
		// 获取窗口标题
		titleLength := win.GetWindowTextLength(hwnd)
		if titleLength > 0 {
			title := make([]uint16, titleLength+1)
			win.GetWindowText(hwnd, &title[0], int32(len(title)))
			titleStr := syscall.UTF16ToString(title)
			
			// 获取窗口类名
			className := make([]uint16, 256)
			win.GetClassName(hwnd, &className[0], int32(len(className)))
			classNameStr := syscall.UTF16ToString(className)
			
			// 检查是否为输入法相关窗口
			for _, keyword := range imeKeywords {
				if strings.Contains(strings.ToLower(titleStr), strings.ToLower(keyword)) ||
					strings.Contains(strings.ToLower(classNameStr), strings.ToLower(keyword)) {
					sp.imeWindows[hwnd] = true
					sp.ProtectWindow(hwnd)
					fmt.Printf("找到并保护输入法窗口: %s (类名: %s)\n", titleStr, classNameStr)
					break
				}
			}
		}
		
		return true
	}, 0)
	
	fmt.Printf("总共找到并保护了 %d 个输入法窗口\n", len(sp.imeWindows))
}

// ProtectWindow 保护指定窗口免受录屏
func (sp *ScreenProtection) ProtectWindow(hwnd win.HWND) bool {
	if hwnd == 0 {
		return false
	}
	
	// 使用SetWindowDisplayAffinity API
	procSetWindowDisplayAffinity := sp.user32.NewProc("SetWindowDisplayAffinity")
	dwAffinity := uintptr(0x00000011) // WDA_EXCLUDEFROMCAPTURE
	
	ret, _, _ := procSetWindowDisplayAffinity.Call(
		uintptr(hwnd),
		dwAffinity,
	)
	
	return ret != 0
}

// DetectDynamicPosition 动态检测屏幕和超星窗口位置
func (sp *ScreenProtection) DetectDynamicPosition() {
	// 查找超星学习通窗口
	chaoxingWindow := sp.FindChaoxingWindow()
	
	if chaoxingWindow != 0 {
		sp.chaoxingWindow = chaoxingWindow
		
		// 获取超星窗口位置
		var rect WindowRect
		if sp.getWindowRect(chaoxingWindow, &rect) {
			// 计算程序窗口的最佳位置
			windowWidth := int32(300)
			
			if rect.Right+windowWidth <= sp.screenWidth {
				// 放在超星窗口右侧
				sp.dynamicX = rect.Right + 10
				sp.dynamicY = rect.Top + 50
			} else {
				// 放在超星窗口左侧
				sp.dynamicX = rect.Left - windowWidth - 10
				if sp.dynamicX < 0 {
					sp.dynamicX = 0
				}
				sp.dynamicY = rect.Top + 50
			}
		}
	} else {
		// 如果没找到超星窗口，使用默认位置
		sp.dynamicX = 0
		sp.dynamicY = sp.screenHeight * 40 / 100 // 屏幕高度的40%位置
	}
}

// FindChaoxingWindow 查找超星学习通窗口
func (sp *ScreenProtection) FindChaoxingWindow() win.HWND {
	var chaoxingHwnd win.HWND
	
	// 超星相关的窗口标题关键词
	chaoxingKeywords := []string{
		"超星学习通", "学习通", "Chaoxing", "超星",
		"StudyTong", "超星客户端", "学习通客户端",
		"超星阅读器", "chaoxing", "xuexitong",
	}
	
	// 枚举所有顶层窗口
	win.EnumWindows(func(hwnd win.HWND) bool {
		if !win.IsWindowVisible(hwnd) {
			return true
		}
		
		// 获取窗口标题
		titleLength := win.GetWindowTextLength(hwnd)
		if titleLength > 0 {
			title := make([]uint16, titleLength+1)
			win.GetWindowText(hwnd, &title[0], int32(len(title)))
			titleStr := syscall.UTF16ToString(title)
			
			// 检查是否包含超星相关关键词
			for _, keyword := range chaoxingKeywords {
				if strings.Contains(strings.ToLower(titleStr), strings.ToLower(keyword)) {
					chaoxingHwnd = hwnd
					fmt.Printf("找到超星窗口: %s\n", titleStr)
					return false // 停止枚举
				}
			}
		}
		
		return true
	}, 0)
	
	return chaoxingHwnd
}

// getWindowRect 获取窗口矩形
func (sp *ScreenProtection) getWindowRect(hwnd win.HWND, rect *WindowRect) bool {
	procGetWindowRect := sp.user32.NewProc("GetWindowRect")
	
	ret, _, _ := procGetWindowRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(rect)),
	)
	
	return ret != 0
}

// GetDynamicPosition 获取动态计算的位置
func (sp *ScreenProtection) GetDynamicPosition() (int32, int32) {
	return sp.dynamicX, sp.dynamicY
}

// ManualRefreshIMEProtection 手动刷新输入法保护
func (sp *ScreenProtection) ManualRefreshIMEProtection() {
	sp.FindAndProtectIMEWindows()
	
	message := fmt.Sprintf("✅ 检测到 %d 个输入法窗口\n🔒 已对所有输入法窗口应用防录屏保护", len(sp.imeWindows))
	if len(sp.imeWindows) == 0 {
		message = "⚠️ 未检测到活动的输入法窗口\n💡 请确保输入法正在运行并处于可见状态"
	}
	
	walk.MsgBox(nil, "输入法保护", message, walk.MsgBoxIconInformation)
}

// GetIMEWindowsCount 获取当前检测到的输入法窗口数量
func (sp *ScreenProtection) GetIMEWindowsCount() int {
	return len(sp.imeWindows)
}

// GetChaoxingWindow 获取超星窗口句柄
func (sp *ScreenProtection) GetChaoxingWindow() win.HWND {
	return sp.chaoxingWindow
}

// IsProtectionActive 检查保护是否激活
func (sp *ScreenProtection) IsProtectionActive() bool {
	return sp.protectionActive
}

// GetScreenSize 获取屏幕尺寸
func (sp *ScreenProtection) GetScreenSize() (int32, int32) {
	return sp.screenWidth, sp.screenHeight
}

// ForceUpdatePosition 强制更新位置
func (sp *ScreenProtection) ForceUpdatePosition() {
	sp.DetectDynamicPosition()
}

// GetProtectionStats 获取保护统计信息
func (sp *ScreenProtection) GetProtectionStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["protection_active"] = sp.protectionActive
	stats["ime_windows_count"] = len(sp.imeWindows)
	stats["chaoxing_window_found"] = sp.chaoxingWindow != 0
	stats["screen_width"] = sp.screenWidth
	stats["screen_height"] = sp.screenHeight
	stats["dynamic_x"] = sp.dynamicX
	stats["dynamic_y"] = sp.dynamicY
	
	return stats
} 