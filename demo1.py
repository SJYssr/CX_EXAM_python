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

def change_weight(event):
    """键盘按下事件处理函数"""
    global is_small, current_opacity
    if is_small:
        root.geometry("300x533+0+380")  # 将窗口大小改为 300x533
        root.attributes("-alpha", current_opacity)  # 设置窗口透明度为之前的透明度
    else:
        root.geometry("5x910+0+0")  # 将窗口大小改为 5x910
        current_opacity = root.attributes("-alpha")  # 保存当前的透明度
        root.attributes("-alpha", 0.1)  # 设置窗口透明度为 0.1
    is_small = not is_small  # 切换窗口大小

root = tk.Tk()
root.title("demo")
root.geometry("300x533+0+380")#设置窗口大小和位置
root.attributes("-alpha", 0.5)  # 设置窗口透明度为 0.5
root.overrideredirect(True)  # 隐藏窗口边框
label = tk.Label(root, text="测试")
label.pack(pady=20)

# 绑定键盘按下事件
root.bind("<F3>", change_weight)

# 绑定鼠标右键点击事件
root.bind("<Button-3>", change_opacity0)  # <Button-3> 表示鼠标右键

# 绑定滚轮事件
root.bind("<Control-MouseWheel>", change_opacity)
# 绑定关闭窗口事件
root.bind("<Escape>", close_window)
current_opacity = 0.5  # 初始透明度设置为 0.5
is_small = False  # 初始窗口大小为 300x500
keep_on_top()  # 启动保持最上层功能

root.mainloop()
