import tkinter as tk
from ctypes import windll, wintypes, byref
import time
from pynput.keyboard import Controller

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

def load_file_content():
    """读取 tiku.txt 文件内容并插入到文本框中"""
    try:
        with open('tiku.txt', 'r', encoding='utf-8') as file:
            content = file.read()
            text_box.insert('1.0', content)
    except FileNotFoundError:
        text_box.insert('1.0', "未找到tiku.txt，请确保文件夹中有此文件")

def highlight_search():
    """高亮显示搜索框内容匹配的所有项"""
    search_term = search_entry.get()
    if search_term:
        text_box.tag_remove('highlight', '1.0', 'end')
        start = '1.0'
        while True:
            start = text_box.search(search_term, start, stopindex='end')
            if not start:
                break
            end = f"{start}+{len(search_term)}c"
            text_box.tag_add('highlight', start, end)
            start = end
        text_box.tag_config('highlight', background='yellow')

def text_input():
    """输入功能"""
    input_text = input_entry.get()
    if not input_text:  # 如果输入框为空，则跳过此函数
        return
    keyboard = Controller()
    time.sleep(5)
    keyboard.type(input_text)
    time.sleep(1)
    input_entry.delete(0, tk.END)  # 清空输入框

root = tk.Tk()
root.title("demo")
root.geometry("300x533+0+380")#设置窗口大小和位置
root.attributes("-alpha", 0.5)  # 设置窗口透明度为 0.5
root.overrideredirect(True)  # 隐藏窗口边框

# 创建顶部框架用于放置搜索框和搜索按钮
top_frame = tk.Frame(root)
top_frame.pack(side="top", fill="x")

search_entry = tk.Entry(top_frame)
search_entry.pack(side="left", fill="x", expand=True, padx=5, pady=5)

search_button = tk.Button(top_frame, text="搜索", command=highlight_search)
search_button.pack(side="right", padx=5, pady=5)

# 创建文本框和滚动条
text_frame = tk.Frame(root)
text_frame.pack(fill="both", expand=True)

scrollbar = tk.Scrollbar(text_frame)
scrollbar.pack(side="right", fill="y")

text_box = tk.Text(text_frame, yscrollcommand=scrollbar.set)
text_box.pack(side="left", fill="both", expand=True)
scrollbar.config(command=text_box.yview)

# 使用 place 方法固定文本框的位置和大小
text_box.place(x=0, y=0, width=285, height=455)

# 创建底部框架用于放置输入框和输入按钮
bottom_frame = tk.Frame(root)
bottom_frame.pack(side="bottom", fill="x")

input_entry = tk.Entry(bottom_frame)
input_entry.pack(side="left", fill="x", expand=True, padx=5, pady=5)

submit_button = tk.Button(bottom_frame, text="输入", command=text_input)
submit_button.pack(side="right", padx=5, pady=5)

# 绑定键盘按下事件
root.bind("<F3>", change_weight)

# 绑定鼠标右键点击事件
root.bind("<Button-3>", change_opacity0)  # <Button-3> 表示鼠标右键

# 绑定滚轮事件，仅当按下Ctrl键时生效
root.bind("<Control-MouseWheel>", change_opacity)
# 绑定关闭窗口事件
root.bind("<Escape>", close_window)
current_opacity = 0.5  # 初始透明度设置为 0.5
is_small = False  # 初始窗口大小为 300x500
keep_on_top()  # 启动保持最上层功能

# 加载文件内容到文本框中
load_file_content()

root.mainloop()
