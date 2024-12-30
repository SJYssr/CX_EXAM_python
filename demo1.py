import tkinter as tk
from ctypes import windll, wintypes

from pynput.keyboard import Controller
import time
# 定义SetWindowDisplayAffinity函数
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
    is_small = not is_small  # 切换窗口大小状态

def load_file_content():
    """读取 tiku.txt 文件内容到文本框中"""
    try:
        with open('tiku.txt', 'r', encoding='utf-8') as file:
            content = file.read()
            text_box.insert('1.0', content)
            text_box.config(state='disabled')  # 状态为不可编辑
    except FileNotFoundError:
        text_box.insert('1.0', "未找到tiku.txt，请确保文件夹中有此文件")

def highlight_search():
    """高亮显示文本框中所有匹配搜索框内容项"""
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
    input_entry.delete(0, 'end')  # 清空输入框

def ai_text_input():
    """AI输入功能"""
    input_text = ai_input_entry.get()
    if not input_text:  # 如果输入框为空，则跳过此函数
        return
    keyboard = Controller()
    time.sleep(5)
    keyboard.type(input_text)
    time.sleep(1)
    ai_input_entry.delete(0, 'end')  # 清空输入框
def switch_to_ai_search():
    """切换到AI搜索界面"""
    main_frame.pack_forget()  # 隐藏主界面框架
    ai_frame.pack(fill="both", expand=True)  # 显示AI搜索界面框架
    search_frame.pack_forget()  # 隐藏搜索框架
    ai_text_box.config(state='disabled')  # 禁止编辑AI回答文本框



def switch_to_main():
    """切换回主界面"""
    ai_frame.pack_forget()  # 隐藏AI搜索界面框架
    search_frame.pack(side="top", fill="x")  # 显示搜索框架
    main_frame.pack(fill="both", expand=True)  # 显示主界面框架


def AI_ask():
    ai_text_box.config(state='normal')  # 允许编辑 ai_text_box
    if ai_text_box.get('1.0', 'end-1c'):  # 检查 ai_text_box 是否有内容
        ai_text_box.delete('1.0', 'end')  # 清空AI回答文本框
    """AI搜索调用的函数"""
    # 此处填写调用ai的主函数
    ai_text_box.insert('1.0', "ai功能未完善")

    # 并且在ai回答的返回语句加入下一行代码实现显示
    # ai_text_box.insert(tk.END, 此处填写具体的回答变量)  # 将AI回答插入到AI回答文本框中
    # ai_text_box.config(state=tk.DISABLED)  # 将文本框设置为不可编辑状态

root = tk.Tk()
root.title("demo")
root.geometry("300x533+0+380")#设置窗口大小和位置
root.attributes("-alpha", 0.5)  # 设置窗口透明度为 0.5
root.configure(bg='white')  # 设置窗口背景颜色为白色
root.overrideredirect(True)  # 隐藏窗口边框

# 创建顶部框架用于放置搜索框和搜索按钮（主界面）
search_frame = tk.Frame(root)
search_frame.pack(side="top", fill="x")

search_entry = tk.Entry(search_frame)
search_entry.pack(side="left", fill="x", expand=True, padx=5, pady=5)
search_entry.configure(foreground='gray')

ai_button = tk.Button(search_frame, text="AI", command=switch_to_ai_search)
ai_button.pack(side="right", padx=5, pady=5)
ai_button.configure(foreground='gray')

search_button = tk.Button(search_frame, text="搜索", command=highlight_search)
search_button.pack(side="right", padx=5, pady=5)
search_button.configure(foreground='gray')

# 创建主界面框架和文本框、滚动条（主界面）
main_frame = tk.Frame(root)
main_frame.pack(fill="both", expand=True)

text_frame = tk.Frame(main_frame)
text_frame.pack(fill="both", expand=True)

scrollbar = tk.Scrollbar(text_frame)
scrollbar.pack(side="right", fill="y")

text_box = tk.Text(text_frame, yscrollcommand=scrollbar.set)
text_box.pack(side="left", fill="both", expand=True)
text_box.configure(foreground='gray')
scrollbar.config(command=text_box.yview)

# 创建底部框架用于放置输入框和输入按钮（主界面）
bottom_frame = tk.Frame(main_frame)
bottom_frame.pack(side="bottom", fill="x")

input_entry = tk.Entry(bottom_frame)
input_entry.pack(side="left", fill="x", expand=True, padx=5, pady=5)
input_entry.configure(foreground='gray')

submit_button = tk.Button(bottom_frame, text="输入", command=text_input)
submit_button.pack(side="right", padx=5, pady=5)
submit_button.configure(foreground='gray')

# 创建顶部框架用于放置搜索框和搜索按钮（AI搜索界面）
ai_frame = tk.Frame(root)
ai_search_frame = tk.Frame(ai_frame)
ai_search_frame.pack(side="top", fill="x")

ai_search_entry = tk.Entry(ai_search_frame)
ai_search_entry.pack(side="left", fill="x", expand=True, padx=5, pady=5)
ai_search_entry.configure(foreground='gray')

back_button = tk.Button(ai_search_frame, text="返回", command=switch_to_main)
back_button.pack(side="right", padx=5, pady=5)
back_button.configure(foreground='gray')

ai_search_button = tk.Button(ai_search_frame, text="AI搜索", command=AI_ask)
ai_search_button.pack(side="right", padx=5, pady=5)
ai_search_button.configure(foreground='gray')

# 创建底部框架用于放置输入框和输入按钮（AI搜索界面）
ai_bottom_frame = tk.Frame(ai_frame)
ai_bottom_frame.pack(side="bottom", fill="x")

ai_input_entry = tk.Entry(ai_bottom_frame)
ai_input_entry.pack(side="left", fill="x", expand=True, padx=5, pady=5)
ai_input_entry.configure(foreground='gray')

ai_submit_button = tk.Button(ai_bottom_frame, text="输入", command=ai_text_input)
ai_submit_button.pack(side="right", padx=5, pady=5)
ai_submit_button.configure(foreground='gray')

# 创建AI搜索界面的文本框用于显示AI的回答
ai_text_box = tk.Text(ai_frame, wrap="word")
ai_text_box.pack(fill="both", expand=True)
ai_text_box.configure(foreground='gray')

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
