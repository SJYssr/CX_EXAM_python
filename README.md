# 超星学习通辅助工具 - Go版本

![GitHub stars](https://img.shields.io/github/stars/SJYssr/CX_EXAM_python?style=social)
![GitHub forks](https://img.shields.io/github/forks/SJYssr/CX_EXAM_python?style=social)
![License](https://img.shields.io/github/license/SJYssr/CX_EXAM_python)

> 🎓 一个功能强大的超星学习通辅助工具，采用Go+Python混合开发，具备防录屏保护、AI智能搜索和自动输入等实用功能。

## ✨ 主要特性

- 🛡️ **防录屏保护** - 强化的防录屏功能，保护程序窗口和输入法窗口
- 🔍 **智能搜索** - 支持题库内容快速搜索，高亮显示匹配结果
- 🤖 **AI智能问答** - 集成Spark AI和Deepseek AI，智能回答问题
- ⚡ **自动输入** - 延时自动输入功能，提高效率
- 📱 **动态窗口** - 自动检测超星窗口位置，智能调整辅助工具位置
- 🎨 **界面友好** - 透明度可调、大小可变的现代化界面
- ⌨️ **快捷键支持** - 丰富的快捷键操作，操作更便捷

## 🚀 快速开始

### 环境要求

- **Go语言**: 1.21+
- **Python**: 3.7+
- **操作系统**: Windows 10/11
- **依赖库**: 详见 `go.mod` 和下方Python依赖

### 安装步骤

1. **克隆仓库**
```bash
git clone https://github.com/SJYssr/CX_EXAM_python.git
cd CX_EXAM_python
```

2. **安装Go依赖**
```bash
go mod tidy
```

3. **安装Python依赖**
```bash
pip install pynput
```

4. **构建项目**
```bash
# 运行构建脚本
build.bat

# 或者手动构建
go build -o cx_exam.exe main.go
```

5. **运行程序**
```bash
# 使用批处理文件运行
run.bat

# 或者直接运行
./cx_exam.exe
```

## 📖 使用说明

### 主界面功能

1. **搜索功能**
   - 在搜索框输入关键词
   - 点击"搜索"按钮高亮显示匹配内容
   - 按回车键可跳转到下一个搜索结果

2. **AI搜索**
   - 点击"AI"按钮进入AI搜索界面
   - 输入问题后点击"AI搜索"获取智能回答
   - 支持实时流式回复（Spark AI）

3. **自动输入**
   - 在输入框输入要发送的文字
   - 点击"输入"按钮，程序将在5秒后自动输入内容

### 快捷键说明

| 快捷键 | 功能 |
|--------|------|
| `F3` | 切换窗口大小（正常/迷你模式） |
| `F4` | 手动刷新输入法保护 |
| `ESC` | 退出程序 |
| `Ctrl + 鼠标左键` | 拖动窗口 |
| `Ctrl + 滚轮` | 调整窗口透明度 |
| `Alt + 滚轮` | 调整文字大小 |
| `右键点击` | 快速切换透明度（0.2/0.5） |
| `回车键` | 跳转到下一个搜索结果 |

### 防录屏功能

程序启动后会自动启用防录屏保护：
- 保护主窗口不被录屏软件捕获
- 自动检测并保护输入法窗口
- 支持主流输入法：搜狗、百度、QQ、微软等
- 可通过F4快捷键手动刷新保护

## 🏗️ 项目架构

```
CX_EXAM_python/
├── main.go              # 主程序入口
├── config.go            # 配置管理
├── ai_service.go        # AI服务管理
├── gui_manager.go       # GUI管理器
├── screen_protection.go # 屏幕保护功能
├── file_handler.go      # 文件处理
├── gui_components.py    # Python GUI组件
├── tiku.txt            # 题库文件
├── build.bat           # 构建脚本
├── run.bat             # 运行脚本
├── go.mod              # Go模块依赖
└── LICENSE             # 许可证
```

## 🔧 配置说明

### AI服务配置

程序支持多种AI服务，可在`config.go`中配置：

- **Spark AI**: 科大讯飞星火认知大模型
- **Deepseek AI**: 深度求索AI模型

### 题库配置

- 题库内容存储在`tiku.txt`文件中
- 支持自定义题库内容
- 程序启动时自动加载题库

## 🤝 贡献指南

我们欢迎任何形式的贡献！

1. Fork 这个仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建一个 Pull Request

## 📋 更新日志

### v2.0.0 - Go版本重构
- 🔄 主程序重构为Go语言，提升性能
- 🛡️ 增强防录屏保护功能
- 🤖 集成多种AI服务
- ⚡ 优化窗口动态定位
- 🎨 改进用户界面体验

## ⚠️ 免责声明

本工具仅用于学习和研究目的，请遵守相关法律法规和学校规定。使用本工具产生的任何后果由用户自行承担。

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE) - 查看 LICENSE 文件了解详情

## 👨‍💻 作者

**SJYssr**
- GitHub: [@SJYssr](https://github.com/SJYssr)

## 🙏 鸣谢

- [lxn/walk](https://github.com/lxn/walk) - Windows GUI框架
- [go-vgo/robotgo](https://github.com/go-vgo/robotgo) - 跨平台GUI自动化
- [pynput](https://github.com/moses-palmer/pynput) - Python输入控制库

---

⭐ 如果这个项目对您有帮助，请给个星标支持一下！ 