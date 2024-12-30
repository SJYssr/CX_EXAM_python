# 该项目用于考试客户端(包括但不限于学习通考试客户端)

## 项目介绍
- 该项目用于考试客户端，主要功能包括本地题库导入，本地题库搜题，AI搜题
- 调用SetWindowDisplayAffinity函数，防止考试客户端的截屏，录屏
- 文本一键输入功能
- 界面可以隐藏(F3)
- 不同的透明度(右键0.2-0.5，CTRL+滚轮0.1-1)
- 快捷退出按键（ESC）
- 最新的[releases](https://github.com/SJYssr/CX_EXAM_python/releases/tag/V2.0.0)中已有打包好的文件可直接使用
- 如果感觉有用的话请给我一颗小星星

## 时间日历
 - 2024.12.26 项目开始，创建代码仓库
 - 2024.12.27 创建[README.md](https://github.com/SJYssr/CX_EXAM_python/blob/main/README.md)和[GPL-3.0 License](https://github.com/SJYssr/CX_EXAM_python/blob/main/LICENSE) 协议
 - 2024.12.27 创建[demo1.py](https://github.com/SJYssr/CX_EXAM_python/blob/main/demo1.py)文件，实现了界面透明度（鼠标右键透明度在0.2-0.5转换，CTRL+滚轮上/下滚动透明度增加/减少），快速退出（esc）
 - 2024.12.27 增加快捷隐藏(F3)将窗口隐藏到左侧边框
 - 2024.12.28 解决截屏/录屏问题（现在截屏/录屏无法接看到此窗口），添加自动导入题库文件（tiku.txt）和搜索功能高亮标记
 - 2024.12.28 基本功能已经完成
 - 2024.12.29 添加一键输入功能(需要为英文输入法)
 - 2024.12.30 添加AI搜题界面框架,将字体颜色改为灰色
 - 2024.12.30 完成AI功能（讯飞星火）

## 注:
- AI逻辑代码在AI_ask函数中
- 如果选择讯飞星火AI,appid、api_key、api_secret需要自行申请([申请地址](https://aiui.xfyun.cn/console))
- 预留的ai可直接使用
- 如果选择其他AI请自行修改代码
- 输入法仍然会被录屏/截屏，能力有限无法解决，如果有大佬可以解决,可以在[issues](https://github.com/SJYssr/CX_EXAM_python/issues/1)留言

# warning: 免责声明
- 本代码遵循 [GPL-3.0 License](https://github.com/SJYssr/CX_EXAM_python/blob/main/LICENSE) 协议，允许**开源/免费使用和引用/修改/衍生代码的开源/免费使用**，不允许**修改和衍生的代码作为闭源的商业软件发布和销售**，禁止**使用本代码盈利**，以此代码为基础的程序**必须**同样遵守 [GPL-3.0 License](https://github.com/SJYssr/CX_EXAM_python/blob/main/LICENSE)协议
- 本代码仅用于**学习讨论**，禁止**用于盈利**
- 他人或组织使用本代码进行的任何**违法行为**与本人无关