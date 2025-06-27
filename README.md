# fuck超星 🎯

一个功能丰富的WPF桌面应用程序，专为学习和考试辅助设计，集成了AI问答、题库搜索、防录屏、文本输入等多种实用功能。

## ✨ 核心功能

### 🛡️ 防录屏保护
- **主窗口防录屏**：防止考试界面被录屏软件捕获
- **多系统兼容**：支持Windows 7及以上版本
- **智能检测**：自动检测系统版本并使用最佳防录屏方案
- **一键切换**：F4快速开启/关闭防录屏功能

### 🤖 AI智能问答
- **双AI引擎**：支持讯飞星火AI和腾讯云Deepseek AI
- **实时问答**：输入问题即可获得AI回答
- **流式输出**：支持实时显示AI回答过程
- **配置灵活**：可在配置文件中切换不同AI服务

### 📚 题库搜索系统
- **关键词搜索**：快速在题库中查找相关内容
- **高亮显示**：搜索结果自动高亮标记
- **下一个结果**：Enter键快速跳转到下一个搜索结果
- **自定义题库**：支持自定义题库内容

### ⌨️ 文本输入辅助
- **延迟输入**：模拟人工输入，避免检测
- **AI答案输入**：将AI回答直接输入到目标窗口
- **智能识别**：自动识别当前焦点窗口

### 🖥️ 窗口管理
- **置顶显示**：窗口始终保持在最前端
- **透明度调节**：Ctrl+鼠标滚轮调节透明度
- **尺寸切换**：F3在正常/迷你模式间切换
- **拖拽移动**：Ctrl+鼠标左键拖拽移动窗口

### 🔧 调试工具
- **窗口选择器**：查看系统中所有窗口信息
- **防录屏测试**：测试防录屏API是否正常工作
- **实时检测**：实时监控窗口状态

## 📋 系统要求

- **操作系统**：Windows 7 及以上版本
- **框架要求**：.NET 6.0 Windows Desktop Runtime
- **推荐配置**：Windows 10 Build 19041+ （获得最佳防录屏效果）

## 🛠️ 技术栈

- **框架**：.NET 6.0 + WPF
- **配置管理**：YamlDotNet
- **JSON处理**：Newtonsoft.Json
- **WebSocket**：WebSocketSharp-netstandard
- **AI服务**：讯飞星火 WebSocket API、Deepseek HTTP API

## 🚀 安装与使用

### 1. 下载程序
```bash
git clone https://github.com/maoqijie/CX_EXAM_python.git
cd fuck超星
```

### 2. 配置AI服务
编辑 `config.yaml` 文件：

```yaml
# AI 服务类型设置
AI_set:
  # AI服务类型: 0=未设置, 1=Spark, 2=Deepseek
  type: 2

# Deepseek AI配置 - 腾讯云免费版本
deepseek:
  # 去腾讯云申请API密钥：https://console.cloud.tencent.com/lke
  api_key: "sk-你的API密钥"
  model: "deepseek-chat"
  api_url: "https://api.lkeap.cloud.tencent.com/v1"

# 讯飞星火AI配置
SPARK:
  Spark_url: "wss://spark-api.xf-yun.com/v3.5/chat"
  api_key: "你的API密钥"
  api_secret: "你的API密钥"
  appid: "你的应用ID"
  domain: "generalv3.5"
```

### 3. 自定义题库
编辑 `tiku.txt` 文件，添加你的题库内容：
```text
你的题库内容
可以是任何文本格式
支持中文和英文
每行一个知识点或问答
```

### 4. 编译运行
```bash
# 编译项目
dotnet build

# 运行程序
dotnet run
```

或者直接运行编译后的可执行文件：
```bash
./bin/Debug/net6.0-windows/CxExamClient.exe
```

## 🎯 使用说明

### 基本操作

1. **题库搜索**
   - 在搜索框输入关键词
   - 按回车或点击搜索按钮
   - 使用回车键跳转到下一个搜索结果

2. **AI问答**
   - 点击"AI"按钮切换到AI界面
   - 输入问题并提交
   - 等待AI回答并查看结果

3. **自动输入**
   - 在输入框中输入要发送的文本
   - 点击输入按钮或按回车
   - 系统将自动输入到目标位置

## 🎮 快捷键说明

| 快捷键 | 功能 |
|--------|------|
| **F3** | 切换窗口大小（正常/迷你模式） |
| **F4** | 切换防录屏状态 |
| **F5** | 退出程序 |
| **Ctrl + 鼠标左键** | 拖拽移动窗口 |
| **Ctrl + 鼠标滚轮** | 调节窗口透明度 |
| **Alt + 鼠标滚轮** | 调节文字大小 |
| **鼠标右键** | 快速切换透明度（0.2/0.5） |
| **Enter** | 搜索模式下跳转下一个结果 |
| **Escape** | 退出程序 |

### 窗口模式

- **正常模式**：300x533像素，半透明显示
- **最小化模式**：5x910像素，贴边显示，几乎隐形

## 📁 项目结构

```
CxExamClient/
├── App.xaml              # 应用程序入口
├── MainWindow.xaml       # 主窗口界面
├── MainWindow.xaml.cs    # 主窗口逻辑
├── config.yaml           # 配置文件
├── tiku.txt             # 题库文件
├── Models/
│   └── ConfigModel.cs    # 配置模型
├── Services/
│   ├── AiService.cs      # AI服务
│   ├── ConfigManager.cs  # 配置管理
│   ├── InputService.cs   # 输入服务
│   ├── SearchHelper.cs   # 搜索助手
│   ├── TikuManager.cs    # 题库管理
│   └── WindowManager.cs  # 窗口管理
└── build.bat            # 构建脚本
```

## 🔒 防录屏技术

### 技术原理
- **WDA_EXCLUDEFROMCAPTURE**：Windows 10+ 最强防录屏API
- **WDA_MONITOR**：Windows 7+ 兼容防录屏方案
- **DWM属性**：桌面窗口管理器预览排除
- **分层窗口**：降级兼容方案

### 兼容性
| 系统版本 | 防录屏效果 | 备注 |
|----------|------------|------|
| Windows 10 Build 19041+ | ✅ 完全黑屏 | 推荐系统 |
| Windows 10 (旧版本) | ✅ 兼容防录屏 | 基本防护 |
| Windows 8/8.1 | ✅ 兼容防录屏 | 基本防护 |
| Windows 7 | ⚠️ 有限防护 | 部分防护 |

## ⚠️ 使用须知

### 使用建议
- 建议以**管理员身份运行**程序以获得最佳效果
- 确保系统为Windows 10 Build 19041+以获得完全防录屏效果
- 首次使用请配置AI服务密钥
- 定期更新题库内容以获得更好的搜索体验

### 注意事项
1. **合规使用**：本工具仅供学习交流使用，请遵守相关法律法规和考试规则
2. **API配额**：腾讯云免费API有配额限制，请合理使用
3. **数据安全**：请妥善保管你的API密钥，不要分享给他人
4. **兼容性**：目前仅支持Windows系统

## 🛠️ 调试工具

### 窗口选择器功能
- **窗口列表**：显示所有系统窗口信息
- **选择IME**：自动选择输入法相关窗口
- **应用保护**：对选中窗口应用防录屏
- **移除保护**：移除窗口的防录屏保护
- **API测试**：测试防录屏API是否正常

### 实时检测
- **30秒监控**：实时检测新出现的窗口
- **智能识别**：自动识别各种输入法窗口
- **即时保护**：发现目标窗口立即应用保护

## 🔧 常见问题

### Q: AI服务无法使用？
A: 请检查：
- 网络连接是否正常
- API密钥是否正确配置
- AI服务类型是否设置正确

### Q: 题库搜索不到内容？
A: 请确认：
- `tiku.txt` 文件是否存在
- 文件编码是否为UTF-8
- 搜索关键词是否正确

### Q: 窗口无法拖拽？
A: 请使用 `Ctrl + 左键` 组合键进行拖拽

### Q: 防录屏功能无效？
A: 请检查：
- 操作系统是否为 Windows 7 及以上版本
- 是否以管理员权限运行程序（推荐）
- 使用 F4 键手动切换防录屏状态
- 不同系统防录屏效果可能不同：
  * Windows 10 2004+：最强防录屏
  * Windows 7-10：兼容模式防录屏
  * Windows XP/Vista：不支持或效果有限

### Q: 如何查看防录屏状态？
A: 
- 按 F4 键，程序会显示详细的防录屏状态信息
- 包括系统版本、防录屏是否启用等信息

## 🖼️ 界面介绍

### 主界面
- **搜索框**：输入关键词搜索题库
- **题库显示区**：显示题库内容和搜索结果
- **输入框**：输入要自动输入的文本
- **AI按钮**：切换到AI问答界面
- **调试按钮**：打开窗口选择器

### AI界面
- **问题输入框**：输入要询问AI的问题
- **AI回答区**：显示AI的回答内容
- **答案输入框**：将AI答案输入到目标窗口
- **返回按钮**：返回主界面

## ⚙️ 配置文件详解

### AI_set 配置
```yaml
AI_set:
  type: 2  # 0=未设置, 1=Spark, 2=Deepseek
```

### 腾讯云Deepseek配置（推荐）
```yaml
deepseek:
  api_key: "sk-xxx"  # 腾讯云API密钥
  model: "deepseek-chat"
  api_url: "https://api.lkeap.cloud.tencent.com/v1"
```

**获取腾讯云API密钥：**
1. 注册腾讯云账号
2. 开通大模型知识引擎服务
3. 创建API密钥
4. 将密钥填入配置文件

### 讯飞星火配置
```yaml
SPARK:
  Spark_url: "wss://spark-api.xf-yun.com/v3.5/chat"
  api_key: "xxx"
  api_secret: "xxx"
  appid: "xxx"
  domain: "generalv3.5"
```

### 题库路径配置
```yaml
PATH:
  path: "tiku.txt"  # 题库文件路径
```

## 📄 开源许可

本项目基于MIT许可证开源，详细信息请查看 [LICENSE](LICENSE) 文件。

## 🤝 贡献指南

欢迎提交Bug报告、功能请求或代码贡献：

1. Fork本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建Pull Request

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 📧 邮箱：[您的邮箱]
- 🐛 Issues：[GitHub Issues](https://github.com/maoqijie/CX_EXAM_python/issues)

---

**⚠️ 免责声明**：本工具仅供学习和研究使用，请遵守相关法律法规和考试规定。使用者应自行承担使用本工具的风险和责任。 