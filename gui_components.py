# _*_coding : UTF_8 _*_
# GUI组件模块
import tkinter as tk
import time
import json
from pynput.keyboard import Controller
import _thread as thread
from tkinter import messagebox
from ctypes import windll, wintypes, WINFUNCTYPE


class GUIComponents:
    def __init__(self, root, ai_service_manager, screen_protection, file_handler):
        self.root = root
        self.ai_service_manager = ai_service_manager
        self.screen_protection = screen_protection
        self.file_handler = file_handler
        
        # 全局变量
        self.current_opacity = 0.5
        self.text_size = 10
        self.is_small = False
        self.current_search_index = 0
        self.search_results = []
        self.x = 0
        self.y = 0
        
        # 动态位置变量
        self.screen_width = 0
        self.screen_height = 0
        self.dynamic_x = 0
        self.dynamic_y = 0
        self.position_check_counter = 0  # 位置检测计数器
        
        # 初始化Windows API
        self._init_windows_api()
        
        # GUI组件
        self.text_box = None
        self.ai_text_box = None
        self.search_entry = None
        self.input_entry = None
        self.ai_search_entry = None
        self.ai_input_entry = None
        self.ai_search_button = None
        self.main_frame = None
        self.ai_frame = None
        self.search_frame = None
        self.ai_search_frame = None
        
        self._detect_dynamic_position()
        self._setup_window()
        self._create_widgets()
        self._bind_events()
        self._load_content()
    
    def _init_windows_api(self):
        """初始化Windows API函数"""
        # 获取屏幕尺寸的API
        self.GetSystemMetrics = windll.user32.GetSystemMetrics
        # 枚举窗口的API
        self.EnumWindows = windll.user32.EnumWindows
        self.EnumWindowsProc = WINFUNCTYPE(wintypes.BOOL, wintypes.HWND, wintypes.LPARAM)
        self.GetWindowTextW = windll.user32.GetWindowTextW
        self.GetWindowTextLengthW = windll.user32.GetWindowTextLengthW
        self.IsWindowVisible = windll.user32.IsWindowVisible
        self.GetWindowRect = windll.user32.GetWindowRect
    
    def _detect_dynamic_position(self):
        """动态检测屏幕和窗口位置"""
        # 获取屏幕尺寸
        self.screen_width = self.GetSystemMetrics(0)  # SM_CXSCREEN
        self.screen_height = self.GetSystemMetrics(1)  # SM_CYSCREEN
        
        # 尝试检测超星学习通窗口
        chaoxing_window = self._find_chaoxing_window()
        
        if chaoxing_window:
            # 如果找到超星窗口，计算相对位置
            rect = wintypes.RECT()
            self.GetWindowRect(chaoxing_window, wintypes.byref(rect))
            
            # 将窗口放置在超星窗口右侧或适当位置
            window_width = rect.right - rect.left
            if rect.right + 300 <= self.screen_width:
                # 放在右侧
                self.dynamic_x = rect.right + 10
                self.dynamic_y = rect.top + 50
            else:
                # 放在左侧
                self.dynamic_x = max(0, rect.left - 310)
                self.dynamic_y = rect.top + 50
        else:
            # 如果没找到超星窗口，使用屏幕比例定位
            # 左边缘，垂直居中偏下
            self.dynamic_x = 0
            self.dynamic_y = int(self.screen_height * 0.4)  # 屏幕高度的40%位置
    
    def _find_chaoxing_window(self):
        """查找超星学习通窗口"""
        chaoxing_hwnd = None
        
        # 超星相关的窗口标题关键词
        chaoxing_keywords = [
            '超星学习通',
            '学习通',
            'Chaoxing',
            '超星',
            'StudyTong',
            '超星客户端',
            '学习通客户端'
        ]
        
        def enum_proc(hwnd, lParam):
            nonlocal chaoxing_hwnd
            try:
                # 检查窗口是否可见
                if not self.IsWindowVisible(hwnd):
                    return True
                
                # 获取窗口标题
                title_length = self.GetWindowTextLengthW(hwnd)
                if title_length > 0:
                    title_buffer = wintypes.create_unicode_buffer(title_length + 1)
                    self.GetWindowTextW(hwnd, title_buffer, title_length + 1)
                    title_str = title_buffer.value
                    
                    # 检查是否包含超星相关关键词
                    for keyword in chaoxing_keywords:
                        if keyword in title_str:
                            chaoxing_hwnd = hwnd
                            return False  # 停止枚举
                            
            except Exception:
                pass
            return True
        
        # 枚举所有顶层窗口
        enum_proc_func = self.EnumWindowsProc(enum_proc)
        self.EnumWindows(enum_proc_func, 0)
        
        return chaoxing_hwnd

    def _setup_window(self):
        """设置窗口属性"""
        self.root.geometry(f"300x533+{self.dynamic_x}+{self.dynamic_y}")
        self.root.attributes("-alpha", 0.5)
        self.root.configure(bg='white')
        self.root.overrideredirect(True)
    
    def _create_widgets(self):
        """创建GUI组件"""
        self._create_main_interface()
        self._create_ai_interface()
    
    def _create_main_interface(self):
        """创建主界面"""
        # 创建顶部框架用于放置搜索框和搜索按钮（主界面）
        self.search_frame = tk.Frame(self.root)
        self.search_frame.pack(side="top", fill="x")

        self.search_entry = tk.Entry(self.search_frame)
        self.search_entry.pack(side="left", fill="x", expand=True, padx=5, pady=5)
        self.search_entry.configure(foreground='gray')

        ai_button = tk.Button(self.search_frame, text="AI", command=self.switch_to_ai_search)
        ai_button.pack(side="right", padx=5, pady=5)
        ai_button.configure(foreground='gray')

        search_button = tk.Button(self.search_frame, text="搜索", command=self.highlight_search)
        search_button.pack(side="right", padx=5, pady=5)
        search_button.configure(foreground='gray')

        # 创建主界面框架和文本框（主界面）
        self.main_frame = tk.Frame(self.root)
        self.main_frame.pack(fill="both", expand=True)

        text_frame = tk.Frame(self.main_frame)
        text_frame.pack(fill="both", expand=True)

        self.text_box = tk.Text(text_frame, wrap='word')
        self.text_box.pack(side="left", fill="both", expand=True)
        self.text_box.configure(foreground='gray')

        # 创建底部框架用于放置输入框和输入按钮（主界面）
        bottom_frame = tk.Frame(self.main_frame)
        bottom_frame.pack(side="bottom", fill="x")

        self.input_entry = tk.Entry(bottom_frame)
        self.input_entry.pack(side="left", fill="x", expand=True, padx=5, pady=5)
        self.input_entry.configure(foreground='gray')

        submit_button = tk.Button(bottom_frame, text="输入", command=self.text_input)
        submit_button.pack(side="right", padx=5, pady=5)
        submit_button.configure(foreground='gray')
    
    def _create_ai_interface(self):
        """创建AI界面"""
        # 创建顶部框架用于放置搜索框和搜索按钮（AI搜索界面）
        self.ai_frame = tk.Frame(self.root)
        self.ai_search_frame = tk.Frame(self.ai_frame)
        self.ai_search_frame.pack(side="top", fill="x")

        self.ai_search_entry = tk.Entry(self.ai_search_frame)
        self.ai_search_entry.pack(side="left", fill="x", expand=True, padx=5, pady=5)
        self.ai_search_entry.configure(foreground='gray')

        back_button = tk.Button(self.ai_search_frame, text="返回", command=self.switch_to_main)
        back_button.pack(side="right", padx=5, pady=5)
        back_button.configure(foreground='gray')

        self.ai_search_button = tk.Button(self.ai_search_frame, text="AI搜索", command=self.ai_ask)
        self.ai_search_button.pack(side="right", padx=5, pady=5)
        self.ai_search_button.configure(foreground='gray')

        # 创建AI搜索界面的文本框用于显示AI的回答
        self.ai_text_box = tk.Text(self.ai_frame, wrap="word")
        self.ai_text_box.pack(fill="both", expand=True)
        self.ai_text_box.configure(foreground='gray')

        # 创建底部框架用于放置输入框和输入按钮（AI搜索界面）
        ai_bottom_frame = tk.Frame(self.ai_frame)
        ai_bottom_frame.pack(side="bottom", fill="x")

        self.ai_input_entry = tk.Entry(ai_bottom_frame)
        self.ai_input_entry.pack(side="left", fill="x", expand=True, padx=5, pady=5)
        self.ai_input_entry.configure(foreground='gray')

        ai_submit_button = tk.Button(ai_bottom_frame, text="输入", command=self.ai_text_input)
        ai_submit_button.pack(side="right", padx=5, pady=5)
        ai_submit_button.configure(foreground='gray')
    
    def _bind_events(self):
        """绑定事件"""
        # 绑定鼠标事件
        self.root.bind("<Control-Button-1>", self.start_move)
        self.root.bind("<Control-B1-Motion>", self.stop_move)
        # 绑定键盘按下事件
        self.root.bind("<F3>", self.change_weight)
        # 绑定手动刷新输入法保护事件
        self.root.bind("<F4>", self.manual_refresh_ime_protection)
        # 绑定鼠标右键点击事件
        self.root.bind("<Button-3>", self.change_opacity0)
        # 绑定滚轮事件，仅当按下Ctrl键时生效
        self.root.bind("<Control-MouseWheel>", self.change_opacity)
        # 绑定滚轮事件，仅当按下ALT键时生效
        self.root.bind("<Alt-MouseWheel>", self.change_text_size)
        # 绑定关闭窗口事件
        self.root.bind("<Escape>", self.close_window)
        # 绑定搜索导航事件
        self.root.bind("<Return>", lambda e: self.next_search_result())
    
    def _load_content(self):
        """加载文件内容"""
        content, success = self.file_handler.load_tiku_content()
        self.text_box.insert('1.0', content)
        self.text_box.config(state='disabled')
    
    def change_opacity(self, event):
        """CTRL+滚轮改变窗口透明度"""
        if self.is_small:
            return
        
        if event.delta > 0:
            self.current_opacity += 0.1
        else:
            self.current_opacity -= 0.1
        self.current_opacity = max(0.1, min(self.current_opacity, 1.0))
        self.root.attributes("-alpha", self.current_opacity)
    
    def change_opacity0(self, event):
        """右键改变窗口透明度0.2/0.5"""
        if self.is_small:
            return
        
        if self.current_opacity == 0.2:
            self.current_opacity = 0.5
        else:
            self.current_opacity = 0.2
        self.root.attributes("-alpha", self.current_opacity)
    
    def close_window(self, event):
        """关闭窗口"""
        self.root.destroy()
    
    def manual_refresh_ime_protection(self, event):
        """手动刷新输入法保护 (快捷键: F4)"""
        self.screen_protection.manual_refresh_ime_protection()
    
    def change_weight(self, event):
        """键盘按下事件处理函数"""
        if self.is_small:
            # 恢复正常大小时，重新检测动态位置
            self._detect_dynamic_position()
            self.root.geometry(f"300x533+{self.dynamic_x}+{self.dynamic_y}")
            self.root.attributes("-alpha", self.current_opacity)
        else:
            # 缩小时放在屏幕左上角
            self.root.geometry("5x910+0+0")
            self.current_opacity = self.root.attributes("-alpha")
            self.root.attributes("-alpha", 0.1)
        self.is_small = not self.is_small
    
    def highlight_search(self):
        """高亮显示文本框中匹配搜索框内容项，并支持单个跳转"""
        search_term = self.search_entry.get()
        if not search_term:
            return

        # 移除所有现有高亮标记
        self.text_box.tag_remove('highlight', '1.0', 'end')
        self.text_box.tag_remove('current_highlight', '1.0', 'end')

        # 收集所有匹配项的位置
        self.search_results = []
        start = '1.0'
        while True:
            start = self.text_box.search(search_term, start, stopindex='end')
            if not start:
                break
            end = f"{start}+{len(search_term)}c"
            self.search_results.append((start, end))
            start = end

        if not self.search_results:
            return

        # 重置当前索引
        self.current_search_index = 0

        # 高亮显示第一个匹配项
        start, end = self.search_results[self.current_search_index]
        self.text_box.tag_add('current_highlight', start, end)
        self.text_box.tag_config('current_highlight', background='yellow')

        # 确保当前高亮项可见
        self.text_box.see(start)
    
    def next_search_result(self):
        """跳转到下一个搜索结果"""
        if not self.search_results:
            return

        # 移除当前高亮
        self.text_box.tag_remove('current_highlight', '1.0', 'end')

        # 移动到下一个结果
        self.current_search_index = (self.current_search_index + 1) % len(self.search_results)
        start, end = self.search_results[self.current_search_index]

        # 高亮显示新的当前项
        self.text_box.tag_add('current_highlight', start, end)
        self.text_box.see(start)
    
    def text_input(self):
        """主界面输入功能"""
        input_text = self.input_entry.get()
        if not input_text:
            return

        def input_thread():
            keyboard = Controller()
            time.sleep(5)
            keyboard.type(input_text)
            time.sleep(1)
            self.input_entry.delete(0, 'end')

        thread.start_new_thread(input_thread, ())
    
    def ai_text_input(self):
        """AI界面输入功能"""
        input_text = self.ai_input_entry.get()
        if not input_text:
            return

        def input_thread():
            keyboard = Controller()
            time.sleep(5)
            keyboard.type(input_text)
            time.sleep(1)
            self.ai_input_entry.delete(0, 'end')

        thread.start_new_thread(input_thread, ())
    
    def switch_to_ai_search(self):
        """切换到AI搜索界面"""
        self.main_frame.pack_forget()
        self.ai_frame.pack(fill="both", expand=True)
        self.search_frame.pack_forget()
        self.ai_text_box.config(state='disabled')
    
    def switch_to_main(self):
        """切换回主界面"""
        self.ai_frame.pack_forget()
        self.search_frame.pack(side="top", fill="x")
        self.main_frame.pack(fill="both", expand=True)
    
    def change_text_size(self, event):
        """调整文本框内字体大小"""
        self.text_size = max(10, min(self.text_size, 12))
        if self.main_frame.winfo_ismapped():
            if event.delta > 0:
                self.text_size += 1
            else:
                self.text_size -= 1
            self.text_box.config(font=("Arial", self.text_size))
        elif self.ai_search_frame.winfo_ismapped():
            if event.delta > 0:
                self.text_size += 1
            else:
                self.text_size -= 1
            self.ai_text_box.config(font=("Arial", self.text_size))
    
    def start_move(self, event):
        """开始移动窗口"""
        self.x = event.x
        self.y = event.y
    
    def stop_move(self, event):
        """停止移动窗口"""
        self.root.geometry(f"+{event.x_root - self.x}+{event.y_root - self.y}")
    
    def ai_ask(self):
        """AI搜索调用的函数"""
        self.ai_search_button.config(state='disabled')
        self.ai_text_box.config(state='normal')
        self.ai_text_box.delete('1.0', tk.END)
        
        ai_type = self.ai_service_manager.get_ai_type()
        
        if ai_type == 1:
            self.ai_text_box.config(state='disabled')
        else:
            self.ai_text_box.insert(tk.END, "正在思考中，请稍候...")
            self.ai_text_box.config(state='disabled')

        def run_ai():
            try:
                query = self.ai_search_entry.get()
                if ai_type == 1:
                    # SparkAI需要回调函数
                    callbacks = {
                        'on_message': self._on_spark_message,
                        'on_error': self._on_spark_error,
                        'on_close': self._on_spark_close
                    }
                    self.ai_service_manager.ask(query, callbacks)
                elif ai_type == 2:
                    # DeepseekAI直接返回结果
                    response = self.ai_service_manager.ask(query)
                    self.root.after(0, self._update_ai_text, response)
            except Exception as e:
                self.root.after(0, self._update_ai_text, f"发生错误: {str(e)}")
            finally:
                self.root.after(0, lambda: self.ai_search_button.config(state='normal'))

        thread.start_new_thread(run_ai, ())
    
    def _update_ai_text(self, content):
        """更新AI文本框内容"""
        self.ai_text_box.config(state='normal')
        self.ai_text_box.delete('1.0', tk.END)
        self.ai_text_box.insert(tk.END, content)
        self.ai_text_box.config(state='disabled')
    
    def _on_spark_message(self, ws, message):
        """处理Spark AI消息"""
        data = json.loads(message)
        self.ai_text_box.config(state='normal')
        if data['header']['code'] != 0:
            self.root.after(0, self._update_ai_text, f"请求错误: {data['header']['code']}")
            ws.close()
        else:
            content = data["payload"]["choices"]["text"][0]["content"]
            self.root.after(0, lambda: self.ai_text_box.insert(tk.END, content))
            if data["payload"]["choices"]["status"] == 2:
                ws.close()
        self.ai_text_box.config(state='disabled')
    
    def _on_spark_error(self, ws, error):
        """处理Spark AI错误"""
        self.root.after(0, self._update_ai_text, f"连接错误: {str(error)}")
    
    def _on_spark_close(self, ws):
        """处理Spark AI连接关闭"""
        self.root.after(0, lambda: self.ai_text_box.insert(tk.END, "\n##连接已关闭##"))
    
    def keep_on_top(self):
        """将窗口始终保持在最上层，并保护输入法窗口不被录屏"""
        # 设置窗口属性，使其始终保持在最上层
        self.root.attributes("-topmost", True)
        
        # 保护当前程序窗口
        self.screen_protection.protect_current_window()
        
        # 查找并保护输入法窗口
        self.screen_protection.find_ime_windows()
        self.screen_protection.protect_all_ime_windows()
        
        # 每10秒钟重新检测一次超星窗口位置（只在正常大小时）
        self.position_check_counter += 1
        if not self.is_small and self.position_check_counter >= 5:  # 5 * 2秒 = 10秒
            self.position_check_counter = 0
            old_x, old_y = self.dynamic_x, self.dynamic_y
            self._detect_dynamic_position()
            # 如果位置发生变化，更新窗口位置
            if old_x != self.dynamic_x or old_y != self.dynamic_y:
                self.root.geometry(f"300x533+{self.dynamic_x}+{self.dynamic_y}")
        
        # 每2秒钟检查一次
        self.root.after(2000, self.keep_on_top) 