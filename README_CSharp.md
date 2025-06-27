# 超星考试客户端工具 (C# 版本)

![.NET](https://img.shields.io/badge/.NET-6.0-blue?style=flat-square)
![WPF](https://img.shields.io/badge/WPF-Application-green?style=flat-square)
![License](https://img.shields.io/badge/license-GPL--3.0-red?style=flat-square)
![Platform](https://img.shields.io/badge/platform-Windows-lightgrey?style=flat-square)

---

> 🚀 **这是超星考试客户端工具的C#重构版本，使用WPF技术栈，提供更好的性能和用户体验。**

---

## 项目介绍

本项目是原Python版本的C#重构版本，采用现代WPF技术栈开发，保持了原有的所有功能特性，同时提供了更好的性能、稳定性和用户体验。

## 主要功能

### 🔍 题库管理与搜索
- **本地题库导入**：自动读取`tiku.txt`文件，支持大批量题目管理
- **智能关键词搜索**：实时高亮匹配项，支持回车键跳转下一个结果
- **题库内容保护**：只读模式防止误操作修改题库内容

### 🤖 AI 智能搜题
- **多AI平台支持**：
  - 讯飞星火AI (WebSocket实时通信)
  - Deepseek AI (HTTP API调用)
- **实时AI问答**：支持流式响应，实时显示AI回答过程
- **智能答案输入**：AI答案可一键自动输入到目标程序

### 🖥️ 界面与交互增强
- **窗口置顶与防录屏**：调用Windows API实现窗口置顶和防截屏功能
- **灵活透明度控制**：
  - 右键快速切换（0.2/0.5）
  - Ctrl+滚轮精细调节（0.1~1.0）
- **字体大小调节**：Alt+滚轮实时调整文字大小
- **智能窗口操作**：
  - F3快速隐藏/显示窗口
  - Ctrl+鼠标左键自由拖动窗口
- **快捷键支持**：ESC/F1-F12快速退出程序

### ⚡ 技术特性
- **异步处理**：所有AI调用和输入操作均为异步，界面响应流畅
- **内存优化**：智能缓存机制，减少文件重复读取
- **错误处理**：完善的异常处理机制，提供友好的错误提示
- **模块化设计**：清晰的项目结构，易于维护和扩展

## 技术栈

- **.NET 6.0**: 现代跨平台框架
- **WPF**: Windows桌面应用程序框架
- **YAML配置**: 使用YamlDotNet解析配置文件
- **HTTP客户端**: 支持Deepseek API调用
- **WebSocket**: 支持讯飞星火实时通信
- **Windows API**: 窗口管理和输入模拟

## 项目结构

```
CxExamClient/
├── Models/                 # 数据模型
│   └── ConfigModel.cs     # 配置文件模型
├── Services/              # 服务层
│   ├── AiService.cs       # AI服务
│   ├── ConfigManager.cs   # 配置管理
│   ├── InputService.cs    # 输入服务
│   ├── SearchHelper.cs    # 搜索助手
│   ├── TikuManager.cs     # 题库管理
│   └── WindowManager.cs   # 窗口管理
├── MainWindow.xaml        # 主窗口界面
├── MainWindow.xaml.cs     # 主窗口逻辑
├── App.xaml              # 应用程序资源
├── App.xaml.cs           # 应用程序入口
├── CxExamClient.csproj   # 项目文件
├── config.yaml           # 配置文件
└── tiku.txt              # 题库文件
```

## 快速开始

### 环境要求
- Windows 10/11
- .NET 6.0 Runtime 或更高版本

### 安装步骤

1. **克隆或下载项目**
   ```bash
   git clone <repository-url>
   cd CxExamClient
   ```

2. **还原NuGet包**
   ```bash
   dotnet restore
   ```

3. **配置AI服务**
   - 编辑`config.yaml`文件
   - 填写相应的AI服务配置信息

4. **准备题库文件**
   - 确保`tiku.txt`文件存在于项目根目录
   - 每行一题，支持UTF-8编码

5. **编译运行**
   ```bash
   dotnet build
   dotnet run
   ```

### 配置说明

`config.yaml`配置文件结构：

```yaml
AI_set:
  type: 1  # 0=未设置, 1=讯飞星火, 2=Deepseek

SPARK:
  Spark_url: "wss://spark-api.xf-yun.com/v1.1/chat"
  api_key: "your_api_key"
  api_secret: "your_api_secret"
  appid: "your_appid"
  domain: "lite"

deepseek:
  api_key: "your_deepseek_api_key"
  model: "deepseek-chat"
```

## 使用说明

### 基本操作
1. **题库搜索**：在搜索框输入关键词，点击搜索或按回车
2. **AI问答**：点击"AI"切换到AI界面，输入问题进行搜索
3. **自动输入**：在输入框中输入内容，点击"输入"按钮自动发送到其他程序

### 快捷键操作
- **F3**：窗口隐藏/恢复
- **Ctrl+鼠标左键**：拖动窗口
- **右键**：切换透明度(0.2/0.5)
- **Ctrl+滚轮**：调整透明度
- **Alt+滚轮**：调整字体大小
- **ESC/F1-F12**：退出程序
- **回车**：搜索下一个结果

## 开发说明

### 添加新的AI服务
1. 在`AiService.cs`中添加新的AI调用方法
2. 在`ConfigModel.cs`中添加对应的配置模型
3. 在`config.yaml`中添加配置选项
4. 更新AI类型判断逻辑

### 自定义界面
- 编辑`App.xaml`中的样式资源
- 修改`MainWindow.xaml`布局
- 调整`MainWindow.xaml.cs`中的事件处理

## 性能优化

- **异步操作**：所有耗时操作均使用async/await模式
- **内存管理**：及时释放资源，避免内存泄漏
- **UI响应**：使用Dispatcher确保UI线程安全
- **缓存机制**：配置和题库内容智能缓存

## 常见问题

**Q: 程序启动时提示缺少文件？**
A: 确保`config.yaml`和`tiku.txt`文件在程序目录下

**Q: AI功能无法使用？**
A: 检查网络连接和API密钥配置是否正确

**Q: 窗口无法置顶？**
A: 确保程序具有足够的权限，某些安全软件可能会阻止

**Q: 自动输入功能失效？**
A: 确保目标程序处于活动状态，输入法为英文模式

## 贡献指南

欢迎提交Issue和Pull Request：

1. Fork项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建Pull Request

## 许可证

本项目基于GPL-3.0许可证开源：
- ✅ 允许商业使用、修改、分发
- ✅ 必须开源并保持同样许可证
- ❌ 不提供任何担保

## 免责声明

**本工具仅供学习交流使用，请勿用于违法违规用途。**

使用本工具产生的任何后果由使用者自行承担，开发者不承担任何法律责任。 