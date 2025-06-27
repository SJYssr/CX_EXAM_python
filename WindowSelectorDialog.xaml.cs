using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.ComponentModel;
using System.Linq;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Media;
using CxExamClient.Services;

namespace CxExamClient
{
    public partial class WindowSelectorDialog : Window
    {
        // Windows API 声明
        [DllImport("user32.dll")]
        private static extern bool EnumWindows(EnumWindowsProc enumProc, IntPtr lParam);

        [DllImport("user32.dll")]
        private static extern int GetWindowText(IntPtr hWnd, StringBuilder lpString, int nMaxCount);

        [DllImport("user32.dll")]
        private static extern int GetClassName(IntPtr hWnd, StringBuilder lpClassName, int nMaxCount);

        [DllImport("user32.dll")]
        private static extern bool IsWindowVisible(IntPtr hWnd);

        [DllImport("user32.dll", SetLastError = true)]
        private static extern bool SetWindowDisplayAffinity(IntPtr hWnd, uint dwAffinity);

        [DllImport("kernel32.dll")]
        private static extern uint GetLastError();

        [DllImport("user32.dll")]
        private static extern uint GetWindowThreadProcessId(IntPtr hWnd, out uint lpdwProcessId);

        [DllImport("kernel32.dll")]
        private static extern IntPtr OpenProcess(uint dwDesiredAccess, bool bInheritHandle, uint dwProcessId);

        [DllImport("kernel32.dll")]
        private static extern bool CloseHandle(IntPtr hObject);

        [DllImport("psapi.dll")]
        private static extern uint GetModuleBaseName(IntPtr hProcess, IntPtr hModule, StringBuilder lpBaseName, uint nSize);

        private delegate bool EnumWindowsProc(IntPtr hWnd, IntPtr lParam);

        private const uint WDA_EXCLUDEFROMCAPTURE = 0x00000011;
        private const uint PROCESS_QUERY_INFORMATION = 0x0400;
        private const uint PROCESS_VM_READ = 0x0010;

        // 数据模型
        public class WindowInfo : INotifyPropertyChanged
        {
            private bool _isSelected;
            private string _status = "未保护";
            private Brush _statusColor = Brushes.Black;

            public IntPtr Handle { get; set; }
            public string Title { get; set; } = "";
            public string ClassName { get; set; } = "";
            public string ProcessName { get; set; } = "";
            public bool IsVisible { get; set; }
            
            public bool IsSelected
            {
                get => _isSelected;
                set
                {
                    _isSelected = value;
                    OnPropertyChanged(nameof(IsSelected));
                }
            }

            public string Status
            {
                get => _status;
                set
                {
                    _status = value;
                    StatusColor = value switch
                    {
                        "已保护" => Brushes.Green,
                        "保护失败" => Brushes.Red,
                        "已移除" => Brushes.Orange,
                        _ => Brushes.Black
                    };
                    OnPropertyChanged(nameof(Status));
                }
            }

            public Brush StatusColor
            {
                get => _statusColor;
                set
                {
                    _statusColor = value;
                    OnPropertyChanged(nameof(StatusColor));
                }
            }

            public event PropertyChangedEventHandler? PropertyChanged;
            protected virtual void OnPropertyChanged(string propertyName)
            {
                PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(propertyName));
            }
        }

        private readonly ObservableCollection<WindowInfo> _windows = new ObservableCollection<WindowInfo>();
        private readonly HashSet<IntPtr> _protectedWindows = new HashSet<IntPtr>();

        public WindowSelectorDialog()
        {
            InitializeComponent();
            WindowListView.ItemsSource = _windows;
            RefreshWindowList();
        }

        private void RefreshButton_Click(object sender, RoutedEventArgs e)
        {
            RefreshWindowList();
        }

        private void SelectImeOnlyButton_Click(object sender, RoutedEventArgs e)
        {
            // 先取消所有选择
            foreach (var window in _windows)
            {
                window.IsSelected = false;
            }
            
            // 只选择输入法相关窗口
            foreach (var window in _windows)
            {
                if (IsPossibleImeWindow(window))
                {
                    window.IsSelected = true;
                }
            }
            
            int selectedCount = _windows.Count(w => w.IsSelected);
            StatusText.Text = $"已选择 {selectedCount} 个输入法相关窗口";
        }

        private void UnselectAllButton_Click(object sender, RoutedEventArgs e)
        {
            foreach (var window in _windows)
            {
                window.IsSelected = false;
            }
        }

        private void ApplyButton_Click(object sender, RoutedEventArgs e)
        {
            ApplyProtection();
        }

        private void RemoveButton_Click(object sender, RoutedEventArgs e)
        {
            RemoveProtection();
        }

        private void TestApiButton_Click(object sender, RoutedEventArgs e)
        {
            TestAPI();
        }

        private void RealTimeImeButton_Click(object sender, RoutedEventArgs e)
        {
            StartRealTimeImeDetection();
        }

        private void CloseButton_Click(object sender, RoutedEventArgs e)
        {
            Close();
        }

        private void RefreshWindowList()
        {
            try
            {
                StatusText.Text = "正在刷新窗口列表...";
                _windows.Clear();

                // 枚举所有窗口
                EnumWindows(EnumWindowCallback, IntPtr.Zero);

                // 按是否可见和标题排序
                var sortedWindows = _windows.OrderByDescending(w => w.IsVisible)
                                           .ThenBy(w => string.IsNullOrEmpty(w.Title))
                                           .ThenBy(w => w.Title)
                                           .ToList();

                _windows.Clear();
                foreach (var window in sortedWindows)
                {
                    _windows.Add(window);
                }

                StatusText.Text = $"找到 {_windows.Count} 个窗口";
            }
            catch (Exception ex)
            {
                StatusText.Text = $"刷新失败: {ex.Message}";
                MessageBox.Show($"刷新窗口列表失败:\n{ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }

        private bool EnumWindowCallback(IntPtr hWnd, IntPtr lParam)
        {
            try
            {
                // 获取窗口信息
                var titleBuilder = new StringBuilder(256);
                GetWindowText(hWnd, titleBuilder, titleBuilder.Capacity);
                string title = titleBuilder.ToString();

                var classBuilder = new StringBuilder(256);
                GetClassName(hWnd, classBuilder, classBuilder.Capacity);
                string className = classBuilder.ToString();

                bool isVisible = IsWindowVisible(hWnd);

                // 获取进程名
                string processName = GetProcessName(hWnd);

                // 创建窗口信息
                var windowInfo = new WindowInfo
                {
                    Handle = hWnd,
                    Title = string.IsNullOrEmpty(title) ? "[无标题]" : title,
                    ClassName = className,
                    ProcessName = processName,
                    IsVisible = isVisible,
                    Status = _protectedWindows.Contains(hWnd) ? "已保护" : "未保护"
                };

                // 高亮显示可能的输入法窗口
                if (IsPossibleImeWindow(windowInfo))
                {
                    windowInfo.IsSelected = true;
                }

                _windows.Add(windowInfo);
                return true;
            }
            catch
            {
                return true; // 继续枚举
            }
        }

        private string GetProcessName(IntPtr hWnd)
        {
            try
            {
                uint processId;
                GetWindowThreadProcessId(hWnd, out processId);

                IntPtr processHandle = OpenProcess(PROCESS_QUERY_INFORMATION | PROCESS_VM_READ, false, processId);
                if (processHandle != IntPtr.Zero)
                {
                    var processName = new StringBuilder(256);
                    if (GetModuleBaseName(processHandle, IntPtr.Zero, processName, (uint)processName.Capacity) > 0)
                    {
                        CloseHandle(processHandle);
                        return processName.ToString();
                    }
                    CloseHandle(processHandle);
                }
                return $"PID:{processId}";
            }
            catch
            {
                return "未知";
            }
        }

        private bool IsPossibleImeWindow(WindowInfo window)
        {
            // 更精确的输入法窗口识别 - 只选择最可能的候选词窗口

            // 高置信度的输入法类名（直接匹配）
            string[] highConfidenceImeClasses = {
                "MSCTFIME UI",                    // 微软输入法UI
                "MSCandUIWindow_Candidate",       // 微软候选词窗口
                "MSCandUIWindow_Composition",     // 微软编辑窗口
                "MSCTF_CANDIDATEWINDOW",          // 微软CTF候选窗口
                "QQPYCandWnd",                    // QQ拼音候选词
                "SogouPYCandidate",               // 搜狗候选词
                "WeChatCandWnd",                  // 微信候选词
                "BaiduPYCandWnd",                 // 百度候选词
                "XunfeiInputCandWnd",             // 讯飞候选词
                "TouchKeyboard",                  // 触摸键盘
                "IPTip_Main_Window"               // 触摸键盘主窗口
            };

            // 检查高置信度类名（精确匹配）
            if (highConfidenceImeClasses.Any(className => 
                string.Equals(window.ClassName, className, StringComparison.OrdinalIgnoreCase)))
            {
                return true;
            }

            // 检查包含特定输入法关键词的类名
            if (window.ClassName.Contains("MSCTFIME", StringComparison.OrdinalIgnoreCase) ||
                window.ClassName.Contains("CandWnd", StringComparison.OrdinalIgnoreCase) ||
                window.ClassName.Contains("Candidate", StringComparison.OrdinalIgnoreCase))
            {
                return true;
            }

            // 检查明确的输入法标题关键词
            string[] explicitImeKeywords = {
                "候选", "拼音", "输入法", "Candidate", "IME"
            };

            if (explicitImeKeywords.Any(keyword => 
                window.Title.Contains(keyword, StringComparison.OrdinalIgnoreCase)))
            {
                return true;
            }

            // 检查知名输入法品牌
            string[] imeBrands = { "搜狗", "QQ", "百度", "讯飞", "微信" };
            if (imeBrands.Any(brand => 
                window.Title.Contains(brand, StringComparison.OrdinalIgnoreCase) ||
                window.ClassName.Contains(brand, StringComparison.OrdinalIgnoreCase)))
            {
                return true;
            }

            // 不再自动选择所有空标题窗口，避免误选
                         return false;
        }

        private void TestAPI()
        {
            try
            {
                StatusText.Text = "正在测试API...";
                
                // 获取系统信息
                var osVersion = Environment.OSVersion.Version;
                bool isWin10_2004Plus = osVersion.Major > 10 || (osVersion.Major == 10 && osVersion.Build >= 19041);
                
                // 找到主程序窗口
                IntPtr mainWindowHandle = IntPtr.Zero;
                string mainWindowTitle = "";
                
                foreach (var window in _windows)
                {
                    if (window.ProcessName.Contains("CxExamClient", StringComparison.OrdinalIgnoreCase) && 
                        window.IsVisible && 
                        !string.IsNullOrEmpty(window.Title) && 
                        window.Title != "[无标题]")
                    {
                        mainWindowHandle = window.Handle;
                        mainWindowTitle = window.Title;
                        break;
                    }
                }
                
                if (mainWindowHandle == IntPtr.Zero)
                {
                    StatusText.Text = "测试失败: 找不到主程序窗口";
                    MessageBox.Show("无法找到主程序窗口进行测试！", "测试失败", MessageBoxButton.OK, MessageBoxImage.Warning);
                    return;
                }
                
                // 测试API
                GetLastError(); // 清除之前的错误
                bool testResult = SetWindowDisplayAffinity(mainWindowHandle, WDA_EXCLUDEFROMCAPTURE);
                uint errorCode = GetLastError();
                
                string testMessage = $"API测试结果:\n\n" +
                                   $"🔧 系统信息:\n" +
                                   $"• Windows版本: {osVersion}\n" +
                                   $"• 支持WDA_EXCLUDEFROMCAPTURE: {(isWin10_2004Plus ? "是" : "否")}\n" +
                                   $"• 管理员权限: {(Environment.UserName.Contains("Administrator") ? "是" : "检测中...")}\n\n" +
                                   $"🎯 测试目标:\n" +
                                   $"• 窗口句柄: {mainWindowHandle}\n" +
                                   $"• 窗口标题: {mainWindowTitle}\n" +
                                   $"• 进程名: CxExamClient.exe\n\n" +
                                   $"📊 测试结果:\n" +
                                   $"• API调用: {(testResult ? "成功 ✅" : "失败 ❌")}\n" +
                                   $"• 错误码: {errorCode}\n";
                
                if (!testResult)
                {
                    string errorReason = errorCode switch
                    {
                        5 => "访问被拒绝 - 需要更高权限",
                        6 => "句柄无效",
                        87 => "参数无效",
                        1400 => "窗口句柄无效",
                        1008 => "尝试访问无效地址",
                        0 => "API调用失败(可能系统不支持)",
                        _ => $"未知错误"
                    };
                    
                    testMessage += $"• 失败原因: {errorReason}\n\n";
                    
                    testMessage += $"💡 建议:\n";
                    if (errorCode == 5)
                    {
                        testMessage += $"• 尝试以管理员身份重新运行程序\n";
                    }
                    else if (errorCode == 0 && !isWin10_2004Plus)
                    {
                        testMessage += $"• 系统版本过低，建议升级到Windows 10 Build 19041+\n";
                    }
                    else if (errorCode == 6 || errorCode == 1400)
                    {
                        testMessage += $"• 窗口句柄问题，程序可能正在关闭\n";
                    }
                    else
                    {
                        testMessage += $"• 系统可能不支持此API或存在其他限制\n";
                    }
                }
                else
                {
                    testMessage += $"• API工作正常！ ✅\n\n";
                    testMessage += $"💡 API测试成功，问题可能在于:\n";
                    testMessage += $"• 选择的窗口类型不支持防录屏\n";
                    testMessage += $"• 某些特定进程有保护机制\n";
                    testMessage += $"• 输入法窗口的特殊属性\n";
                    
                    // 恢复窗口状态
                    SetWindowDisplayAffinity(mainWindowHandle, 0);
                }
                
                StatusText.Text = testResult ? "API测试成功" : "API测试失败";
                MessageBox.Show(testMessage, "API测试结果", MessageBoxButton.OK, 
                    testResult ? MessageBoxImage.Information : MessageBoxImage.Warning);
            }
            catch (Exception ex)
            {
                StatusText.Text = "测试异常";
                                 MessageBox.Show($"测试过程中发生异常:\n{ex.Message}", "测试异常", MessageBoxButton.OK, MessageBoxImage.Error);
             }
         }

         private void StartRealTimeImeDetection()
         {
             var result = MessageBox.Show(
                 "🔍 实时输入法检测器\n\n" +
                 "准备开始实时检测输入法窗口！\n\n" +
                 "📋 操作步骤:\n" +
                 "1. 点击确定后，立即打开记事本或任何文本编辑器\n" +
                 "2. 切换到中文输入法（微软拼音）\n" +
                 "3. 开始输入拼音，让候选词窗口出现\n" +
                 "4. 保持候选词窗口显示状态\n" +
                 "5. 程序将自动检测并尝试保护输入法窗口\n\n" +
                 "⏱️ 检测将持续30秒\n\n" +
                 "准备好了吗？",
                 "实时输入法检测", 
                 MessageBoxButton.OKCancel, 
                 MessageBoxImage.Information);

             if (result == MessageBoxResult.OK)
             {
                 StartRealTimeDetection();
             }
         }

         private async void StartRealTimeDetection()
         {
             try
             {
                 StatusText.Text = "实时检测中... (30秒)";
                 
                 int detectionCount = 0;
                 int protectedCount = 0;
                 var detectedImeWindows = new List<WindowInfo>();
                 var processedWindows = new HashSet<IntPtr>(); // 跟踪已处理的窗口

                 // 首先检测所有现有窗口
                 StatusText.Text = "检测现有窗口...";

                 // 然后实时检测新窗口30秒
                 for (int i = 0; i < 30; i++)
                 {
                     await Task.Delay(1000); // 每秒检测一次
                     
                     StatusText.Text = $"实时检测中... ({30-i}秒) - 已检测:{detectionCount} 已保护:{protectedCount}";
                     
                     // 枚举当前所有窗口，包括新窗口和现有窗口
                     var currentWindows = new List<IntPtr>();
                     EnumWindows((hWnd, lParam) => { currentWindows.Add(hWnd); return true; }, IntPtr.Zero);
                     
                     // 检查所有窗口（包括之前未处理的）
                     foreach (var hWnd in currentWindows)
                     {
                         if (!processedWindows.Contains(hWnd) && IsWindowVisible(hWnd))
                         {
                             // 获取窗口信息
                             var titleBuilder = new StringBuilder(256);
                             GetWindowText(hWnd, titleBuilder, titleBuilder.Capacity);
                             string title = titleBuilder.ToString();

                             var classBuilder = new StringBuilder(256);
                             GetClassName(hWnd, classBuilder, classBuilder.Capacity);
                             string className = classBuilder.ToString();

                             string processName = GetProcessName(hWnd);

                             var windowInfo = new WindowInfo
                             {
                                 Handle = hWnd,
                                 Title = string.IsNullOrEmpty(title) ? "[无标题]" : title,
                                 ClassName = className,
                                 ProcessName = processName,
                                 IsVisible = true
                             };

                             // 使用扩展的输入法检测
                             if (IsLikelyImeWindow(windowInfo))
                             {
                                 detectionCount++;
                                 detectedImeWindows.Add(windowInfo);
                                 
                                 // 立即尝试保护
                                 GetLastError();
                                 bool success = SetWindowDisplayAffinity(hWnd, WDA_EXCLUDEFROMCAPTURE);
                                 
                                 if (success)
                                 {
                                     windowInfo.Status = "已保护";
                                     protectedCount++;
                                     _protectedWindows.Add(hWnd);
                                 }
                                 else
                                 {
                                     uint errorCode = GetLastError();
                                     windowInfo.Status = $"保护失败({errorCode})";
                                 }
                                 
                                 // 添加到窗口列表
                                 _windows.Insert(0, windowInfo); // 插入到顶部
                                 
                                 StatusText.Text = $"发现输入法! 类名:{className} - 已保护:{protectedCount}/{detectionCount}";
                             }
                             
                             // 标记为已处理
                             processedWindows.Add(hWnd);
                         }
                     }
                 }
                 
                 // 显示检测结果
                 string resultMessage = $"🔍 实时输入法检测完成！\n\n" +
                                      $"📊 检测结果:\n" +
                                      $"• 检测到输入法窗口: {detectionCount} 个\n" +
                                      $"• 成功保护: {protectedCount} 个\n" +
                                      $"• 检测时长: 30 秒\n\n";

                 if (detectedImeWindows.Count > 0)
                 {
                     resultMessage += $"🎯 检测到的输入法窗口:\n";
                     foreach (var window in detectedImeWindows.Take(5)) // 只显示前5个
                     {
                         resultMessage += $"• {window.ClassName} - {window.Status}\n";
                     }
                     
                     if (protectedCount > 0)
                     {
                         resultMessage += $"\n✅ 成功保护了 {protectedCount} 个输入法窗口！\n" +
                                        $"现在这些窗口在录屏时应该显示为黑屏。";
                     }
                     else
                     {
                         resultMessage += $"\n❌ 虽然检测到输入法窗口，但保护失败。\n" +
                                        $"可能这些窗口类型不支持防录屏API。";
                     }
                 }
                 else
                 {
                     resultMessage += $"❌ 未检测到输入法窗口。\n\n" +
                                    $"💡 建议:\n" +
                                    $"• 确保输入法候选词窗口确实出现了\n" +
                                    $"• 尝试不同的输入法（搜狗、QQ等）\n" +
                                    $"• 输入更多文字让候选词持续显示";
                 }
                 
                 StatusText.Text = $"检测完成: 保护 {protectedCount}/{detectionCount} 个窗口";
                 MessageBox.Show(resultMessage, "实时检测结果", MessageBoxButton.OK, 
                     protectedCount > 0 ? MessageBoxImage.Information : MessageBoxImage.Warning);
             }
             catch (Exception ex)
             {
                 StatusText.Text = "检测异常";
                 MessageBox.Show($"实时检测过程中发生异常:\n{ex.Message}", "检测异常", MessageBoxButton.OK, MessageBoxImage.Error);
             }
         }

         private bool IsLikelyImeWindow(WindowInfo window)
         {
             // 大幅扩展的输入法窗口检测 - 包含调试信息
             
             // 1. 检查所有可能的微软输入法类名
             string[] msImeClasses = {
                 "MSCTFIME UI",
                 "MSCandUIWindow_Candidate", 
                 "MSCandUIWindow_Composition",
                 "MSCTF_CANDIDATEWINDOW",
                 "MSCTFIME Composition",
                 "MSCTFIME",
                 "CTF_FLOATWND",
                 "CiceroUIWndFrame",
                 "CicMarshalWnd",
                 "Windows.UI.Core.CoreWindow",
                 "ApplicationFrameWindow",
                 "ApplicationFrameInputSink",
                 "TextInputHost",
                 "InputMethod.CandidateWindow.Host",
                 "Microsoft.InputMethod.CandidateWindow"
             };

             // 2. 第三方输入法类名
             string[] thirdPartyImeClasses = {
                 "QQPYCandWnd", "QQPinyin", "TXGuiFoundation",
                 "SogouPYCandidate", "SogouPY", "SogouPYWnd", "SOGOU_UNICODE",
                 "WeChatCandWnd", "WeChatInput", "WeChatIME",
                 "BaiduPYCandWnd", "BaiduType", "BaiduPinyin",
                 "XunfeiInputCandWnd", "iFlyInputWnd", "iFlyIME"
             };

             // 3. 通用输入法相关类名
             string[] genericImeClasses = {
                 "CandidateWindow", "CompositionWindow", "InputWindow",
                 "TouchKeyboard", "IPTip_Main_Window", "IMEWindow"
             };

             var allImeClasses = msImeClasses.Concat(thirdPartyImeClasses).Concat(genericImeClasses);

             // 精确匹配
             if (allImeClasses.Any(cls => 
                 string.Equals(window.ClassName, cls, StringComparison.OrdinalIgnoreCase)))
             {
                 return true;
             }

             // 部分匹配 - 更宽松的检测
             string[] partialMatches = {
                 "Cand", "MSCTFIME", "CTF", "IME", "Input", "Composition", 
                 "Sogou", "QQ", "WeChat", "Baidu", "iFly"
             };

             if (partialMatches.Any(partial => 
                 window.ClassName.Contains(partial, StringComparison.OrdinalIgnoreCase)))
             {
                 return true;
             }

             // 检查标题关键字
             if (!string.IsNullOrEmpty(window.Title) && window.Title != "[无标题]")
             {
                 string[] titleKeywords = { 
                     "候选", "输入法", "拼音", "五笔", "Candidate", "IME", "Input",
                     "搜狗", "QQ", "百度", "讯飞", "微信", "微软"
                 };
                 if (titleKeywords.Any(keyword => 
                     window.Title.Contains(keyword, StringComparison.OrdinalIgnoreCase)))
                 {
                     return true;
                 }
             }

             // 检查进程名
             string[] imeProcesses = { 
                 "sogou", "qq", "baidu", "iflyime", "wechat", "tencent",
                 "microsoft", "msctf", "input", "ime"
             };
             if (imeProcesses.Any(proc => 
                 window.ProcessName.Contains(proc, StringComparison.OrdinalIgnoreCase)))
             {
                 return true;
             }

             // 特殊检测：小窗口且无标题（可能是候选词窗口）
             if (window.IsVisible && 
                 (string.IsNullOrEmpty(window.Title) || window.Title == "[无标题]") &&
                 !window.ClassName.Contains("Button") &&
                 !window.ClassName.Contains("Static") &&
                 !window.ClassName.Contains("Edit"))
             {
                 // 这可能是候选词窗口，但需要进一步验证
                 // 排除明显不是输入法的窗口类名
                 string[] excludeClasses = {
                     "Chrome_", "Firefox", "Edge", "Explorer", "Notepad",
                     "Shell_", "Progman", "WorkerW", "DV2ControlHost"
                 };

                 if (!excludeClasses.Any(exclude => 
                     window.ClassName.Contains(exclude, StringComparison.OrdinalIgnoreCase)))
                 {
                     return true; // 可能是候选词窗口
                 }
             }

             return false;
         }

        private void ApplyProtection()
        {
            try
            {
                var selectedWindows = _windows.Where(w => w.IsSelected).ToList();
                if (!selectedWindows.Any())
                {
                    MessageBox.Show("请先选择要保护的窗口！", "提示", MessageBoxButton.OK, MessageBoxImage.Information);
                    return;
                }

                // 检查Windows版本兼容性
                var osVersion = Environment.OSVersion.Version;
                bool isWindows10Plus = osVersion.Major >= 10;
                bool isWin10_2004Plus = osVersion.Major > 10 || (osVersion.Major == 10 && osVersion.Build >= 19041);

                if (!isWindows10Plus)
                {
                    var result = MessageBox.Show(
                        $"检测到Windows版本: {osVersion}\n\n" +
                        $"WDA_EXCLUDEFROMCAPTURE API需要Windows 10或更高版本。\n" +
                        $"在当前系统上可能无法正常工作。\n\n" +
                        $"是否仍要尝试？", 
                        "兼容性警告", 
                        MessageBoxButton.YesNo, 
                        MessageBoxImage.Warning);
                    
                    if (result == MessageBoxResult.No)
                        return;
                }

                int successCount = 0;
                int failCount = 0;

                StatusText.Text = $"正在应用保护... (Windows {osVersion})";

                                 foreach (var window in selectedWindows)
                 {
                     try
                     {
                         // 检查窗口是否仍然有效
                         if (!IsWindowVisible(window.Handle))
                         {
                             window.Status = "窗口已关闭";
                             failCount++;
                             continue;
                         }

                         // 先测试一个简单的窗口操作
                         uint processId;
                         GetWindowThreadProcessId(window.Handle, out processId);
                         
                         if (processId == 0)
                         {
                             window.Status = "无效句柄";
                             failCount++;
                             continue;
                         }

                         // 清除上一次的错误
                         GetLastError();
                         
                         bool success = SetWindowDisplayAffinity(window.Handle, WDA_EXCLUDEFROMCAPTURE);
                         if (success)
                         {
                             window.Status = "已保护";
                             _protectedWindows.Add(window.Handle);
                             successCount++;
                         }
                         else
                         {
                             // 获取详细的错误码
                             uint errorCode = GetLastError();
                             
                             string errorReason = errorCode switch
                             {
                                 5 => "访问被拒绝",
                                 6 => "句柄无效",
                                 87 => "参数无效", 
                                 1400 => "窗口句柄无效",
                                 1008 => "尝试访问无效地址",
                                 0 => "API调用失败(无错误码)",
                                 _ => $"错误码:{errorCode}"
                             };
                             
                             // 进一步诊断
                             if (window.ProcessName.ToLower().Contains("system") || 
                                window.ProcessName.ToLower().Contains("winlogon") ||
                                window.ProcessName.ToLower().Contains("csrss"))
                             {
                                 window.Status = $"系统进程({errorReason})";
                             }
                             else if (errorCode == 5)
                             {
                                 window.Status = $"权限不足({errorReason})";
                             }
                             else if (errorCode == 6 || errorCode == 1400)
                             {
                                 window.Status = $"句柄无效({errorReason})";
                             }
                             else
                             {
                                 window.Status = $"失败({errorReason})";
                             }
                             failCount++;
                         }
                     }
                     catch (Exception ex)
                     {
                         window.Status = $"异常: {ex.Message.Substring(0, Math.Min(10, ex.Message.Length))}";
                         failCount++;
                     }
                 }

                StatusText.Text = $"保护完成: 成功 {successCount} 个，失败 {failCount} 个";

                // 添加详细的系统信息和诊断
                string systemInfo = $"系统信息:\n" +
                                   $"• Windows版本: {osVersion}\n" +
                                   $"• 是否支持WDA_EXCLUDEFROMCAPTURE: {(isWin10_2004Plus ? "是" : "否")}\n" +
                                   $"• 总窗口数: {_windows.Count}\n" +
                                   $"• 选中窗口数: {selectedWindows.Count}\n\n";

                string message = $"防录屏保护应用完成！\n\n" +
                               $"✅ 成功保护: {successCount} 个窗口\n" +
                               $"❌ 保护失败: {failCount} 个窗口\n\n" +
                               systemInfo +
                               $"💡 提示:\n" +
                               $"• 受保护的窗口在录屏时会显示为黑屏\n" +
                               $"• Windows 10 Build 19041+效果最佳\n" +
                               $"• 确保输入法候选词窗口处于显示状态\n" +
                               $"• 某些系统窗口无法被保护是正常现象";

                // 如果全部失败，提供额外的诊断信息
                if (successCount == 0 && failCount > 0)
                {
                    message += $"\n\n🔧 全部失败的可能原因:\n" +
                              $"• 系统版本不支持此API\n" +
                              $"• 需要管理员权限运行程序\n" +
                              $"• 所选窗口类型不支持防录屏\n" +
                              $"• 请尝试以管理员身份运行程序";
                }

                MessageBox.Show(message, "保护结果", MessageBoxButton.OK, MessageBoxImage.Information);
            }
            catch (Exception ex)
            {
                StatusText.Text = $"应用保护失败: {ex.Message}";
                MessageBox.Show($"应用保护失败:\n{ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }

        private void RemoveProtection()
        {
            try
            {
                var selectedWindows = _windows.Where(w => w.IsSelected).ToList();
                if (!selectedWindows.Any())
                {
                    MessageBox.Show("请先选择要移除保护的窗口！", "提示", MessageBoxButton.OK, MessageBoxImage.Information);
                    return;
                }

                int successCount = 0;
                int failCount = 0;

                StatusText.Text = "正在移除保护...";

                foreach (var window in selectedWindows)
                {
                    try
                    {
                        bool success = SetWindowDisplayAffinity(window.Handle, 0);
                        if (success)
                        {
                            window.Status = "已移除";
                            _protectedWindows.Remove(window.Handle);
                            successCount++;
                        }
                        else
                        {
                            window.Status = "移除失败";
                            failCount++;
                        }
                    }
                    catch
                    {
                        window.Status = "移除失败";
                        failCount++;
                    }
                }

                StatusText.Text = $"移除完成: 成功 {successCount} 个，失败 {failCount} 个";

                string message = $"防录屏保护移除完成！\n\n" +
                               $"✅ 成功移除: {successCount} 个窗口\n" +
                               $"❌ 移除失败: {failCount} 个窗口";

                MessageBox.Show(message, "移除结果", MessageBoxButton.OK, MessageBoxImage.Information);
            }
            catch (Exception ex)
            {
                StatusText.Text = $"移除保护失败: {ex.Message}";
                MessageBox.Show($"移除保护失败:\n{ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }
    }
} 