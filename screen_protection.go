package main

import (
	"fmt"
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

// FindAndProtectIMEWindows 查找并保护输入法窗口 (简化实现)
func (sp *ScreenProtection) FindAndProtectIMEWindows() {
	// 清空之前的输入法窗口列表
	sp.imeWindows = make(map[win.HWND]bool)
	
	// 暂时简化实现，避免使用复杂的Windows API
	fmt.Printf("查找输入法窗口功能暂时简化\n")
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

// FindChaoxingWindow 查找超星学习通窗口 (简化实现)
func (sp *ScreenProtection) FindChaoxingWindow() win.HWND {
	// 暂时简化实现，返回0表示未找到
	// 实际实现需要枚举窗口和窗口标题匹配
	return 0
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