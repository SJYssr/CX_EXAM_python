import tkinter as tk
from ctypes import windll, wintypes, byref

# 定义 SetWindowDisplayAffinity 函数
SetWindowDisplayAffinity = windll.user32.SetWindowDisplayAffinity
SetWindowDisplayAffinity.argtypes = [wintypes.HWND, wintypes.DWORD]
SetWindowDisplayAffinity.restype = wintypes.BOOL

def keep_on_top():
    """将窗口始终保持在最上层"""
    root.attributes("-topmost", True)
    hwnd = windll.user32.GetForegroundWindow()
    dwAffinity = 0x00000011  # 设置显示亲和性为 0x00000011
    SetWindowDisplayAffinity(hwnd, dwAffinity)
    root.after(1000, keep_on_top)  # 每秒钟检查一次

def change_opacity(event):
    """改变窗口透明度"""
    global current_opacity
    if event.delta > 0:  # 滚轮向上滚动
        current_opacity += 0.1
    else:  # 滚轮向下滚动
        current_opacity -= 0.1
    current_opacity = max(0.1, min(current_opacity, 1.0))  # 限制透明度在 0.1 到 1.0 之间
    root.attributes("-alpha", current_opacity)

def change_opacity0(event):
    """改变窗口透明度0.2/0.5"""
    global current_opacity
    if current_opacity == 0.2:
        current_opacity = 0.5
    else:
        current_opacity = 0.2
    root.attributes("-alpha", current_opacity)

def close_window(event):
    """关闭窗口"""
    root.destroy()

root = tk.Tk()
root.title("demo")
root.geometry("300x500")

label = tk.Label(root, text="测试")
label.pack(pady=20)

# 绑定鼠标右键点击事件
root.bind("<Button-3>", change_opacity0)  # <Button-3> 表示鼠标右键

# 绑定滚轮事件
root.bind("<Control-MouseWheel>", change_opacity)
# 绑定关闭窗口事件
root.bind("<F2>", close_window)
current_opacity = 0.5  # 初始透明度设置为 0.5
keep_on_top()  # 启动保持最上层功能

root.mainloop()
