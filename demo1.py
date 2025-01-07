import tkinter as tk
from ctypes import windll, wintypes

from pynput.keyboard import Controller

import _thread as thread
import time
import base64
import datetime
import hashlib
import hmac
import json
from urllib.parse import urlparse
import ssl
from datetime import datetime
from time import mktime
from urllib.parse import urlencode
from wsgiref.handlers import format_date_time
import websocket


def keep_on_top():
    """将窗口始终保持在最上层"""
    # 定义 SetWindowDisplayAffinity 函数
    SetWindowDisplayAffinity = windll.user32.SetWindowDisplayAffinity
    SetWindowDisplayAffinity.argtypes = [wintypes.HWND, wintypes.DWORD]
    SetWindowDisplayAffinity.restype = wintypes.BOOL
    root.attributes("-topmost", True)
    hwnd = windll.user32.GetForegroundWindow()
    dwAffinity = 0x0000001  # 设置显示亲和性为 0x0000001
    SetWindowDisplayAffinity(hwnd, dwAffinity)
    root.after(1000, keep_on_top)  # 每秒钟检查一次

def change_opacity(event):
    """CTRL+滚轮改变窗口透明度"""
    global current_opacity,is_small
    if is_small != False:
        return
    else:
        if event.delta > 0:  # 滚轮向上滚动
            current_opacity += 0.1
        else:  # 滚轮向下滚动
            current_opacity -= 0.1
        current_opacity = max(0.1, min(current_opacity, 1.0))  # 限制透明度在 0.1 到 1.0 之间
        root.attributes("-alpha", current_opacity)

def change_opacity0(event):
    """右键改变窗口透明度0.2/0.5"""
    global current_opacity,is_small
    if is_small != False:
        return
    else:
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
        text_box.config(state='disabled')  # 状态为不可编辑

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
    """主界面输入功能"""
    input_text = input_entry.get()
    if not input_text:  # 如果输入框为空，则跳过此函数
        return
    keyboard = Controller()
    time.sleep(5)
    keyboard.type(input_text)
    time.sleep(1)
    input_entry.delete(0, 'end')  # 清空输入框

def ai_text_input():
    """AI界面输入功能"""
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

def change_text_size(event):
    """调整文本框内字体大小"""
    global text_size
    text_size = max(10, min(text_size, 12))  # 限制字体大小在1到12之间
    if main_frame.winfo_ismapped():
        """调整题库文本大小"""
        if event.delta > 0:  # 滚轮向上滚动
            text_size += 1
        else:  # 滚轮向下滚动
            text_size -= 1
        text_box.config(font=("Arial", text_size))
    elif ai_search_frame.winfo_ismapped():
        """调整AI回答文本大小"""
        if event.delta > 0:  # 滚轮向上滚动
            text_size += 1
        else:  # 滚轮向下滚动
            text_size -= 1
        ai_text_box.config(font=("Arial", text_size))

def start_move(event):
    global x, y
    x = event.x
    y = event.y

def stop_move(event):
    root.geometry(f"+{event.x_root - x}+{event.y_root - y}")

def AI_ask(
        # appid、api_secret、api_key三个服务认证信息请前往开放平台控制台查看（https://console.xfyun.cn/services/bm35）

        # Spark_url="wss://spark-api.xf-yun.com/v3.5/chat",   # Max环境的地址
		# Spark_url = "wss://spark-api.xf-yun.com/v4.0/chat"  # 4.0Ultra环境的地址
        # Spark_url = "wss://spark-api.xf-yun.com/v3.1/chat"  # Pro环境的地址
        # Spark_url = "wss://spark-api.xf-yun.com/v1.1/chat"  # Lite环境的地址

        # domain="generalv3.5",   # Max版本
		# domain = "4.0Ultra"     # 4.0Ultra 版本
        # domain = "generalv3"    # Pro版本
        # domain = "lite"         # Lite版本址
        ):
    global appid,api_secret,api_key,Spark_url,domain
    """AI搜索调用的函数"""
    if ai_text_box.get('1.0', 'end-1c'):  # 检查 ai_text_box 是否有内容
        ai_text_box.config(state='normal')  # 允许编辑 ai_text_box
        ai_text_box.delete('1.0', 'end')  # 清空AI回答文本框
    wsParam = Ws_Param(appid, api_key, api_secret, Spark_url)
    websocket.enableTrace(False)
    wsUrl = wsParam.create_url()
    query = ai_search_entry.get()
    ws = websocket.WebSocketApp(wsUrl, on_message=on_message, on_error=on_error, on_close=on_close, on_open=on_open)
    ws.appid = appid
    ws.query = query
    ws.domain = domain
    ws.run_forever(sslopt={"cert_reqs": ssl.CERT_NONE})


class Ws_Param(object):
    # 初始化
    def __init__(self, APPID, APIKey, APISecret, gpt_url):
        self.APPID = APPID
        self.APIKey = APIKey
        self.APISecret = APISecret
        self.host = urlparse(gpt_url).netloc
        self.path = urlparse(gpt_url).path
        self.gpt_url = gpt_url

    # 生成url
    def create_url(self):
        # 生成RFC1123格式的时间戳
        now = datetime.now()
        date = format_date_time(mktime(now.timetuple()))

        # 拼接字符串
        signature_origin = "host: " + self.host + "\n"
        signature_origin += "date: " + date + "\n"
        signature_origin += "GET " + self.path + " HTTP/1.1"

        # 进行hmac-sha256进行加密
        signature_sha = hmac.new(self.APISecret.encode('utf-8'), signature_origin.encode('utf-8'),
                                 digestmod=hashlib.sha256).digest()

        signature_sha_base64 = base64.b64encode(signature_sha).decode(encoding='utf-8')

        authorization_origin = f'api_key="{self.APIKey}", algorithm="hmac-sha256", headers="host date request-line", signature="{signature_sha_base64}"'

        authorization = base64.b64encode(authorization_origin.encode('utf-8')).decode(encoding='utf-8')

        # 将请求的鉴权参数组合为字典
        v = {
            "authorization": authorization,
            "date": date,
            "host": self.host
        }
        # 拼接鉴权参数，生成url
        url = self.gpt_url + '?' + urlencode(v)
        # 此处打印出建立连接时候的url,参考本demo的时候可取消上方打印的注释，比对相同参数时生成的url与自己代码生成的url是否一致
        return url


# 收到websocket错误的处理
def on_error(ws, error):
    ai_text_box.insert('1.0', error)


# 收到websocket关闭的处理
def on_close(ws):
    ai_text_box.insert('1.0', '##连接已关闭##')


# 收到websocket连接建立的处理
def on_open(ws):
    thread.start_new_thread(run, (ws,))


def run(ws, *args):
    data = json.dumps(gen_params(appid=ws.appid, query=ws.query, domain=ws.domain))
    ws.send(data)


# 收到websocket消息的处理
def on_message(ws, message):
    data = json.loads(message)
    code = data['header']['code']
    if code != 0:
        ai_text_box.insert('1.0', f'请求错误: {code}, {data}')
        ws.close()
    else:
        ai_text_box.config(state='normal')  # 允许编辑 ai_text_box
        choices = data["payload"]["choices"]
        status = choices["status"]
        content = choices["text"][0]["content"]
        ai_text_box.insert(tk.END, content)
        if status == 2:
            ws.close()
    ai_text_box.config(state='disabled')  # 将文本框设置为不可编辑状态



def gen_params(appid, query, domain):
    """
    通过appid和用户的提问来生成请参数
    """

    data = {
        "header": {
            "app_id": appid,
            "uid": "1234",
            # "patch_id": []    #接入微调模型，对应服务发布后的resourceid
        },
        "parameter": {
            "chat": {
                "domain": domain,
                "temperature": 0.5,
                "max_tokens": 4096,
                "auditing": "default",
            }
        },
        "payload": {
            "message": {
                "text": [{"role": "user", "content": query}]
            }
        }
    }
    return data

config = configparser.ConfigParser()
config.read('config.ini', encoding='utf-8')
appid = config.get('SPARK', 'appid')
api_secret = config.get('SPARK', 'api_secret')
api_key = config.get('SPARK', 'api_key')
Spark_url = config.get('SPARK', 'Spark_url')
domain = config.get('SPARK', 'domain')

# 界面部件
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

# 创建主界面框架和文本框（主界面）
main_frame = tk.Frame(root)
main_frame.pack(fill="both", expand=True)

text_frame = tk.Frame(main_frame)
text_frame.pack(fill="both", expand=True)

text_box = tk.Text(text_frame,wrap='word')
text_box.pack(side="left", fill="both", expand=True)
text_box.configure(foreground='gray')

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

# 创建AI搜索界面的文本框用于显示AI的回答
ai_text_box = tk.Text(ai_frame, wrap="word")
ai_text_box.pack(fill="both", expand=True)
ai_text_box.configure(foreground='gray')

# 创建底部框架用于放置输入框和输入按钮（AI搜索界面）
ai_bottom_frame = tk.Frame(ai_frame)
ai_bottom_frame.pack(side="bottom", fill="x")

ai_input_entry = tk.Entry(ai_bottom_frame)
ai_input_entry.pack(side="left", fill="x", expand=True, padx=5, pady=5)
ai_input_entry.configure(foreground='gray')

ai_submit_button = tk.Button(ai_bottom_frame, text="输入", command=ai_text_input)
ai_submit_button.pack(side="right", padx=5, pady=5)
ai_submit_button.configure(foreground='gray')


# 绑定鼠标事件
root.bind("<Button-1>", start_move)
root.bind("<B1-Motion>", stop_move)
# 绑定键盘按下事件
root.bind("<F3>", change_weight)
# 绑定鼠标右键点击事件
root.bind("<Button-3>", change_opacity0)  # <Button-3> 表示鼠标右键
# 绑定滚轮事件，仅当按下Ctrl键时生效
root.bind("<Control-MouseWheel>", change_opacity)
# 绑定滚轮事件，仅当按下ALT键时生效
root.bind("<Alt-MouseWheel>", change_text_size)
# 绑定关闭窗口事件
root.bind("<Escape>", close_window)
current_opacity = 0.5  # 初始透明度设置为 0.5
text_size = 10  # 初始文字大小为 10
is_small = False  # 初始窗口大小为 300x500
keep_on_top()  # 启动保持最上层功能
# 加载文件内容到文本框中
load_file_content()
root.mainloop()
