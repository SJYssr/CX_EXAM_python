using System;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Input;
using System.Windows.Controls;
using System.Windows.Media;
using System.Windows.Documents;
using CxExamClient.Services;

namespace CxExamClient
{
    public partial class MainWindow : Window
    {
        private readonly SearchHelper _searchHelper;
        private readonly AiService _aiService;
        private double _currentOpacity = 0.5;
        private bool _isSmall = false;
        private int _textSize = 10;
        private Point _mousePosition;

        public MainWindow()
        {
            try
            {
                InitializeComponent();
                
                _searchHelper = new SearchHelper();
                _aiService = new AiService();

                InitializeWindow();
                LoadTikuContent();
                SetupEventHandlers();
                // 暂时注释掉ShowAiStatus，可能导致问题
                // ShowAiStatus();
            }
            catch (Exception ex)
            {
                System.IO.File.WriteAllText("startup_error.log", 
                    $"Startup Error: {ex.Message}\nStackTrace: {ex.StackTrace}");
                throw;
            }
        }

        private void InitializeWindow()
        {
            try
            {
                // 配置窗口（包括置顶和防录屏）
                WindowManager.ConfigureWindow(this);
                
                // 设置初始透明度
                Opacity = _currentOpacity;
            }
            catch (Exception ex)
            {
                System.IO.File.WriteAllText("window_init_error.log", 
                    $"Window Init Error: {ex.Message}\nStackTrace: {ex.StackTrace}");
                // 不重新抛出异常，继续执行
            }
        }

        private void LoadTikuContent()
        {
            var content = TikuManager.LoadTikuContent();
            var document = new FlowDocument();
            var paragraph = new Paragraph();
            paragraph.Inlines.Add(new Run(content));
            document.Blocks.Add(paragraph);
            TextBox.Document = document;
        }

        private void SetupEventHandlers()
        {
            // 鼠标事件
            MouseLeftButtonDown += MainWindow_MouseLeftButtonDown;
            MouseMove += MainWindow_MouseMove;
            
            // 键盘事件
            KeyDown += MainWindow_KeyDown;
            
            // 鼠标滚轮事件
            MouseWheel += MainWindow_MouseWheel;
            
            // 右键事件
            MouseRightButtonDown += MainWindow_MouseRightButtonDown;
        }

        private void ShowAiStatus()
        {
            var config = ConfigManager.Config;
            string message = config.AI_set.type switch
            {
                0 => "当前未设置AI",
                1 => "正在使用SparkAI",
                2 => "正在使用DeepseekAI",
                _ => "请检查config.yaml文件中的AI_set配置项"
            };
            
            MessageBox.Show(message, "AI设置", MessageBoxButton.OK, MessageBoxImage.Information);
        }

        #region 界面切换
        private void AiButton_Click(object sender, RoutedEventArgs e)
        {
            SwitchToAiSearch();
        }

        private void BackButton_Click(object sender, RoutedEventArgs e)
        {
            SwitchToMain();
        }

        private void SwitchToAiSearch()
        {
            MainFrame.Visibility = Visibility.Collapsed;
            AiFrame.Visibility = Visibility.Visible;
            AiTextBox.IsReadOnly = true;
        }

        private void SwitchToMain()
        {
            AiFrame.Visibility = Visibility.Collapsed;
            MainFrame.Visibility = Visibility.Visible;
        }
        #endregion

        #region 搜索功能
        private void SearchButton_Click(object sender, RoutedEventArgs e)
        {
            HighlightSearch();
        }

        private void SearchEntry_KeyDown(object sender, KeyEventArgs e)
        {
            if (e.Key == Key.Enter)
            {
                HighlightSearch();
            }
        }

        private void HighlightSearch()
        {
            var searchTerm = SearchEntry.Text;
            if (string.IsNullOrEmpty(searchTerm) || searchTerm == "输入关键词搜索...")
                return;

            _searchHelper.HighlightText(TextBox, searchTerm);
        }

        private void NextSearchResult()
        {
            _searchHelper.NextResult(TextBox);
        }
        #endregion

        #region 输入功能
        private void SubmitButton_Click(object sender, RoutedEventArgs e)
        {
            var text = InputEntry.Text;
            if (!string.IsNullOrEmpty(text) && text != "输入内容后点击输入按钮...")
            {
                _ = InputService.SendTextWithDelay(text);
                InputEntry.Clear();
                InputEntry.Text = "输入内容后点击输入按钮...";
            }
        }

        private void InputEntry_KeyDown(object sender, KeyEventArgs e)
        {
            if (e.Key == Key.Enter)
            {
                SubmitButton_Click(sender, e);
            }
        }

        private void AiSubmitButton_Click(object sender, RoutedEventArgs e)
        {
            var text = AiInputEntry.Text;
            if (!string.IsNullOrEmpty(text) && text != "输入AI答案...")
            {
                _ = InputService.SendTextWithDelay(text);
                AiInputEntry.Clear();
                AiInputEntry.Text = "输入AI答案...";
            }
        }

        private void AiInputEntry_KeyDown(object sender, KeyEventArgs e)
        {
            if (e.Key == Key.Enter)
            {
                AiSubmitButton_Click(sender, e);
            }
        }
        #endregion

        #region AI功能
        private async void AiSearchButton_Click(object sender, RoutedEventArgs e)
        {
            var question = AiSearchEntry.Text;
            if (string.IsNullOrEmpty(question) || question == "输入问题进行AI搜索...")
                return;

                         AiSearchButton.IsEnabled = false;
             AiTextBox.IsReadOnly = false;
             var thinkingDocument = new FlowDocument();
             var thinkingParagraph = new Paragraph();
             thinkingParagraph.Inlines.Add(new Run("正在思考中，请稍候..."));
             thinkingDocument.Blocks.Add(thinkingParagraph);
             AiTextBox.Document = thinkingDocument;
             AiTextBox.IsReadOnly = true;

            try
            {
                                 var response = await _aiService.AskAiAsync(question, (content) =>
                 {
                     Dispatcher.Invoke(() =>
                     {
                         AiTextBox.IsReadOnly = false;
                         var range = new TextRange(AiTextBox.Document.ContentEnd, AiTextBox.Document.ContentEnd);
                         range.Text = content;
                         AiTextBox.IsReadOnly = true;
                         AiTextBox.ScrollToEnd();
                     });
                 });

                 AiTextBox.IsReadOnly = false;
                 var responseDocument = new FlowDocument();
                 var responseParagraph = new Paragraph();
                 responseParagraph.Inlines.Add(new Run(response));
                 responseDocument.Blocks.Add(responseParagraph);
                 AiTextBox.Document = responseDocument;
                 AiTextBox.IsReadOnly = true;
            }
            finally
            {
                AiSearchButton.IsEnabled = true;
            }
        }

        private void AiSearchEntry_KeyDown(object sender, KeyEventArgs e)
        {
            if (e.Key == Key.Enter)
            {
                AiSearchButton_Click(sender, e);
            }
        }
        #endregion

        #region 调试功能
        private void DebugButton_Click(object sender, RoutedEventArgs e)
        {
            ShowWindowSelector();
        }

        private void ShowWindowSelector()
        {
            try
            {
                // 创建窗口选择对话框
                var windowSelector = new WindowSelectorDialog();
                windowSelector.Owner = this;
                windowSelector.ShowDialog();
            }
            catch (Exception ex)
            {
                MessageBox.Show($"打开窗口选择器失败:\n{ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }
        #endregion

        #region 窗口操作
        private void MainWindow_MouseLeftButtonDown(object sender, MouseButtonEventArgs e)
        {
            if (Keyboard.Modifiers == ModifierKeys.Control)
            {
                _mousePosition = e.GetPosition(this);
                DragMove();
            }
        }

        private void MainWindow_MouseMove(object sender, MouseEventArgs e)
        {
            if (e.LeftButton == MouseButtonState.Pressed && Keyboard.Modifiers == ModifierKeys.Control)
            {
                var currentPosition = e.GetPosition(this);
                var offset = currentPosition - _mousePosition;
                Left += offset.X;
                Top += offset.Y;
            }
        }

        private void MainWindow_MouseRightButtonDown(object sender, MouseButtonEventArgs e)
        {
            ChangeOpacity0();
        }

        private void MainWindow_MouseWheel(object sender, MouseWheelEventArgs e)
        {
            if (Keyboard.Modifiers == ModifierKeys.Control)
            {
                ChangeOpacity(e.Delta);
            }
            else if (Keyboard.Modifiers == ModifierKeys.Alt)
            {
                ChangeTextSize(e.Delta);
            }
        }

        private void MainWindow_KeyDown(object sender, KeyEventArgs e)
        {
            switch (e.Key)
            {
                case Key.F3:
                    ChangeWeight();
                    break;
                case Key.F4:
                    ToggleAntiScreenCapture();
                    break;

                case Key.Escape:
                case Key.F1:
                case Key.F2:
                case Key.F5:
                case Key.F6:
                case Key.F7:
                case Key.F8:
                case Key.F9:
                case Key.F10:
                case Key.F11:
                case Key.F12:
                    Close();
                    break;
                case Key.Enter:
                    if (MainFrame.Visibility == Visibility.Visible)
                    {
                        NextSearchResult();
                    }
                    break;
            }
        }

        private void ChangeOpacity(int delta)
        {
            if (_isSmall) return;

            if (delta > 0)
                _currentOpacity += 0.1;
            else
                _currentOpacity -= 0.1;

            _currentOpacity = Math.Max(0.1, Math.Min(_currentOpacity, 1.0));
            Opacity = _currentOpacity;
        }

        private void ChangeOpacity0()
        {
            if (_isSmall) return;

            _currentOpacity = _currentOpacity == 0.2 ? 0.5 : 0.2;
            Opacity = _currentOpacity;
        }

        private void ChangeWeight()
        {
            if (_isSmall)
            {
                Width = 300;
                Height = 533;
                Left = 0;
                Top = 380;
                Opacity = _currentOpacity;
            }
            else
            {
                _currentOpacity = Opacity;
                Width = 5;
                Height = 910;
                Left = 0;
                Top = 0;
                Opacity = 0.1;
            }
            _isSmall = !_isSmall;
        }

        private void ChangeTextSize(int delta)
        {
            if (delta > 0)
                _textSize += 1;
            else
                _textSize -= 1;

            _textSize = Math.Max(8, Math.Min(_textSize, 24));

            if (MainFrame.Visibility == Visibility.Visible)
            {
                TextBox.FontSize = _textSize;
            }
            else if (AiFrame.Visibility == Visibility.Visible)
            {
                AiTextBox.FontSize = _textSize;
            }
        }

        private void ToggleAntiScreenCapture()
        {
            try
            {
                bool success = WindowManager.ToggleDisplayAffinity(this);
                
                // 获取系统信息
                string osInfo = WindowManager.GetOSInfo();
                
                // 构建状态信息
                string mainStatus = success ? "✅ 防录屏已启用" : "❌ 防录屏已禁用";
                
                // 显示详细状态消息框
                string statusMessage = $"=== 防录屏状态 ===\n" +
                                     $"系统: {osInfo}\n" +
                                     $"主窗口: {mainStatus}\n\n" +
                                     $"💡 提示:\n" +
                                     $"• F4: 切换防录屏状态";
                
                MessageBox.Show(statusMessage, "防录屏状态", MessageBoxButton.OK, MessageBoxImage.Information);
                
                string titleMessage = success ? "防录屏已启用" : "防录屏已禁用";
                Title = $"超星考试客户端 - {titleMessage}";
                
                // 3秒后恢复原标题
                var timer = new System.Windows.Threading.DispatcherTimer
                {
                    Interval = TimeSpan.FromSeconds(3)
                };
                timer.Tick += (s, e) =>
                {
                    Title = "超星考试客户端";
                    timer.Stop();
                };
                timer.Start();
            }
            catch (Exception ex)
            {
                MessageBox.Show($"切换防录屏失败:\n{ex.Message}", "错误", MessageBoxButton.OK, MessageBoxImage.Error);
                Title = "超星考试客户端 - 防录屏功能异常";
                
                var timer = new System.Windows.Threading.DispatcherTimer
                {
                    Interval = TimeSpan.FromSeconds(3)
                };
                timer.Tick += (s, e) =>
                {
                    Title = "超星考试客户端";
                    timer.Stop();
                };
                timer.Start();
            }
        }


        #endregion

        #region 文本框焦点处理
        private void TextBox_GotFocus(object sender, RoutedEventArgs e)
        {
            var textBox = sender as TextBox;
            if (textBox?.Text == "输入关键词搜索..." || textBox?.Text == "输入内容后点击输入按钮..." || 
                textBox?.Text == "输入问题进行AI搜索..." || textBox?.Text == "输入AI答案...")
            {
                textBox.Text = "";
                textBox.Foreground = Brushes.Black;
            }
        }

        private void TextBox_LostFocus(object sender, RoutedEventArgs e)
        {
            var textBox = sender as TextBox;
            if (string.IsNullOrWhiteSpace(textBox?.Text))
            {
                textBox.Foreground = Brushes.Gray;
                if (textBox == SearchEntry)
                    textBox.Text = "输入关键词搜索...";
                else if (textBox == InputEntry)
                    textBox.Text = "输入内容后点击输入按钮...";
                else if (textBox == AiSearchEntry)
                    textBox.Text = "输入问题进行AI搜索...";
                else if (textBox == AiInputEntry)
                    textBox.Text = "输入AI答案...";
            }
        }
        #endregion

        protected override void OnClosed(EventArgs e)
        {
            try
            {
                // 清理AI服务
                _aiService?.Dispose();
                
                System.Diagnostics.Debug.WriteLine("程序资源清理完成");
            }
            catch (Exception ex)
            {
                System.Diagnostics.Debug.WriteLine($"资源清理失败: {ex.Message}");
            }
            
            base.OnClosed(e);
        }
    }
} 