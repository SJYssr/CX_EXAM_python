# 超星考试客户端工具  

<div align="center">

![license](https://img.shields.io/github/license/SJYssr/CX_EXAM_python?style=flat-square)
![stars](https://img.shields.io/github/stars/SJYssr/CX_EXAM_python?style=flat-square)
![release](https://img.shields.io/github/v/release/SJYssr/CX_EXAM_python?style=flat-square)
![python](https://img.shields.io/badge/python-3.7%2B-blue?style=flat-square)
![platform](https://img.shields.io/badge/platform-windows-lightgrey?style=flat-square)

</div>

---

> 🚀 **本项目旨在为考试客户端（如学习通）提供本地题库、AI搜题、界面增强等多种实用功能。**

---

<details>
<summary>📑 目录</summary>

- [项目介绍](#项目介绍)
- [主要功能](#主要功能)
- [使用说明](#使用说明)
- [常见问题](#常见问题)
- [时间日历](#时间日历)
- [赞赏](#赞赏)
- [注意事项](#注意事项)
- [配置说明](#配置说明)
- [贡献与反馈](#贡献与反馈)
- [免责声明](#免责声明)

</details>

---

## 项目介绍
本项目是一个专为考试客户端（如学习通）设计的辅助工具，集成了本地题库管理、AI智能搜题、窗口防护与界面增强等多项实用功能，提升考试答题效率与体验。

---

## 主要功能

### 题库管理与搜索
- **本地题库导入**：自动读取`tiku.txt`文件，支持大批量题目管理。
- **关键词高亮搜索**：输入关键词后，所有匹配项高亮显示，支持回车键跳转下一个结果。
- **题库内容只读保护**：防止误操作修改题库内容。

### AI 智能搜题
- **多AI平台支持**：可选[讯飞星火](https://aiui.xfyun.cn/console)或[Deepseek](https://www.deepseek.com)。
- **一键AI问答**：输入问题后，AI自动返回答案，支持多线程防止界面卡顿。
- **AI答案一键输入**：AI答案可一键自动输入到目标输入框（需英文输入法）。

### 界面与交互增强
- **窗口置顶与防录屏/截屏**：调用`SetWindowDisplayAffinity`，窗口始终置顶且无法被录屏/截屏工具捕获。
- **窗口透明度调节**：右键一键切换（0.2/0.5），Ctrl+滚轮精细调节（0.1~1.0）。
- **字体大小调节**：Alt+滚轮随时调整题库/AI答案字体大小。
- **窗口快速隐藏/显示**：F3一键隐藏到屏幕边缘，再次按下恢复。
- **窗口自由拖动**：Ctrl+鼠标左键拖动窗口到任意位置。
- **ESC/F1-F12快捷退出**：ESC或任意F1-F12键可快速关闭程序（可自定义）。

### 其他实用功能
- **在线更新检测**：自动检测新版本并提示更新。
- **多线程处理**：AI问答、输入等操作均采用多线程，保证界面流畅不卡顿。
- **详细注释与易用配置**：代码注释详细，`config.yaml`配置简单明了，便于二次开发。

---

## 使用说明

### 环境准备
1. **安装Python**：确保已安装Python 3.7+。
2. **安装依赖**：在命令行中运行以下命令安装所需依赖：
   ```bash
   pip install tkinter requests pyyaml pynput websocket-client
   ```
3. **准备文件**：在程序运行目录下创建`tiku.txt`（题库）和`config.yaml`（配置文件）。

### 配置AI平台
1. **讯飞星火**：
   - 在`config.yaml`中填写`appid`、`api_key`、`api_secret`。
   - 申请地址：[讯飞星火](https://aiui.xfyun.cn/console)。
2. **Deepseek**：
   - 在`config.yaml`中填写`api_key`和`model`。
   - 可直接使用预设，建议更换为自己的密钥。

### 启动程序
1. **运行程序**：在命令行中运行`demo1.py`。
2. **首次启动**：程序会自动检测配置和题库文件。
3. **主界面功能**：
   - 题库搜索：输入关键词，回车跳转下一个结果。
   - AI切换：点击"AI"按钮切换到AI搜题界面。
   - 快捷输入：在输入框中输入内容，点击"输入"按钮自动输入。

### 常用快捷键与操作
- **F3**：窗口隐藏/恢复
- **Ctrl+鼠标左键**：拖动窗口
- **右键**：切换透明度
- **Ctrl+滚轮**：调整透明度
- **Alt+滚轮**：调整字体大小
- **ESC/F1-F12**：快速退出
- **回车**：题库搜索下一个

---

## 常见问题

- **Q: 启动报错"缺少config.yaml文件"？**
  - A: 请确保`config.yaml`在程序同目录下，参考示例配置文件填写。

- **Q: 题库无法加载或搜索？**
  - A: 请确保`tiku.txt`存在且为UTF-8编码，每行一题。

- **Q: AI搜题无响应？**
  - A: 检查网络连接、API密钥是否正确，或更换AI平台。

- **Q: 窗口无法被录屏/截屏？**
  - A: 仅主窗口受保护，输入法弹窗等仍可能被录屏，技术有限无法完全防护。

- **Q: 如何自定义快捷键？**
  - A: 可在`demo1.py`中修改相关绑定代码，注释详细易于调整。

---

## 时间日历
| 日期         | 事件                                                         |
| ------------ | ------------------------------------------------------------ |
| 2024.12.26   | 项目开始，创建代码仓库                                       |
| 2024.12.27   | 创建README和GPL-3.0 License，demo1.py实现透明度、快捷退出等 |
| 2024.12.28   | 解决截屏/录屏，题库导入与高亮搜索                           |
| 2024.12.29   | 添加一键输入功能                                            |
| 2024.12.30   | 完成AI功能（讯飞星火），项目基本完成                        |
| 2025.1.1     | 添加Alt+滚轮调整字体大小                                    |
| 2025.1.6     | 添加窗口可移动（Ctrl+鼠标左键）                             |
| 2025.1.7     | 添加config文件，AI功能更易配置                              |
| 2025.2.13    | 添加Deepseek AI                                             |
| 2025.2.25    | 添加前置文件查找、详细注释                                   |
| 2025.2.27    | 多线程处理防止堵塞                                          |
| 2025.4.22    | 在线更新、查找下一个功能                                     |

---

## 赞赏
<p align="center">
  <img src="https://github.com/SJYssr/img/raw/main/1/zanshang.jpg" width="250" />
</p>

---

## 注意事项
> ⚠️ **请确保运行目录下有`tiku.txt`和`config.yaml`文件**
- AI逻辑代码在`AI_ask`函数中
- 讯飞星火AI需自行申请密钥
- 预留AI可直接使用，建议更换为自己的
- 其他AI请自行修改代码
- 隐藏(F3)时窗口透明度降到最低，拉成细条放左侧
- 输入法内容仍可能被录屏/截屏，能力有限无法解决，如有方案欢迎[issues](https://github.com/SJYssr/CX_EXAM_python/issues/1)留言
- 如需要破解复制粘贴功能&&篡改猴相关功能，请移步[cef_cx_copy_tool](https://github.com/SJYssr/cef_cx_copy_tool)

---

## 配置说明
- `config.yaml`文件：
  - `type=0` 未设置AI
  - `type=1` 讯飞星火AI
  - `type=2` DeepseekAI

---

## 贡献与反馈
欢迎提交 [Issues](https://github.com/SJYssr/CX_EXAM_python/issues) 反馈问题或建议，或直接 Fork/PR 参与开发。

---

## 免责声明
> **本代码仅用于学习讨论，禁止用于盈利或违法用途。**

- 遵循 [GPL-3.0 License](https://github.com/SJYssr/CX_EXAM_python/blob/main/LICENSE) 协议：
  - 允许开源/免费使用、引用、修改、衍生
  - 禁止闭源商业发布、销售及盈利
  - 基于本代码的程序**必须**同样遵守GPL-3.0协议
- 他人或组织使用本代码进行的任何违法行为与本人无关

---