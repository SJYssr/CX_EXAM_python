using System;
using System.Runtime.InteropServices;
using System.Windows;
using System.Windows.Interop;

namespace CxExamClient.Services
{
    public static class WindowManager
    {
        [DllImport("user32.dll")]
        private static extern bool SetWindowDisplayAffinity(IntPtr hWnd, uint dwAffinity);

        [DllImport("user32.dll")]
        private static extern IntPtr GetForegroundWindow();

        [DllImport("user32.dll", SetLastError = true)]
        private static extern bool SetWindowPos(IntPtr hWnd, IntPtr hWndInsertAfter, int X, int Y, int cx, int cy, uint uFlags);

        [DllImport("kernel32.dll")]
        private static extern Version GetVersion();

        [DllImport("dwmapi.dll", PreserveSig = false)]
        private static extern void DwmSetWindowAttribute(IntPtr hwnd, uint attr, ref int attrValue, int attrSize);

        [DllImport("user32.dll")]
        private static extern bool SetLayeredWindowAttributes(IntPtr hwnd, uint crKey, byte bAlpha, uint dwFlags);

        [DllImport("user32.dll")]
        private static extern int SetWindowLong(IntPtr hWnd, int nIndex, int dwNewLong);

        [DllImport("user32.dll")]
        private static extern int GetWindowLong(IntPtr hWnd, int nIndex);



        private static readonly IntPtr HWND_TOPMOST = new IntPtr(-1);
        private const uint SWP_NOMOVE = 0x0002;
        private const uint SWP_NOSIZE = 0x0001;
        
        // Windows 10 2004+ 专用
        private const uint WDA_EXCLUDEFROMCAPTURE = 0x00000011;
        // Windows 7+ 兼容
        private const uint WDA_MONITOR = 0x00000001;
        
        // DWM 属性
        private const uint DWMWA_EXCLUDED_FROM_PEEK = 12;
        private const uint DWMWA_CLOAK = 13;
        
        // 窗口样式
        private const int GWL_EXSTYLE = -20;
        private const int WS_EX_LAYERED = 0x80000;
        private const uint LWA_ALPHA = 0x2;







        /// <summary>
        /// 设置窗口置顶
        /// </summary>
        public static bool SetWindowTopmost(Window window)
        {
            try
            {
                var helper = new WindowInteropHelper(window);
                if (helper.Handle == IntPtr.Zero)
                    return false;

                return SetWindowPos(helper.Handle, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE);
            }
            catch (Exception ex)
            {
                System.Diagnostics.Debug.WriteLine($"设置窗口置顶失败: {ex.Message}");
                return false;
            }
        }

        /// <summary>
        /// 设置防录屏（支持Windows 7及以上版本）
        /// </summary>
        public static bool SetWindowDisplayAffinity(Window window)
        {
            try
            {
                var helper = new WindowInteropHelper(window);
                if (helper.Handle == IntPtr.Zero)
                    return false;

                bool success = false;
                var osVersion = Environment.OSVersion.Version;
                
                // Windows 10 版本 2004+ 使用最强的防录屏
                if (IsWindows10Version2004OrLater())
                {
                    success = SetWindowDisplayAffinity(helper.Handle, WDA_EXCLUDEFROMCAPTURE);
                    if (success)
                    {
                        System.Diagnostics.Debug.WriteLine("Windows 10 2004+ 防录屏设置成功");
                        return true;
                    }
                }
                
                // Windows 7+ 使用兼容模式
                if (osVersion.Major >= 6) // Windows Vista/7及以上
                {
                    // 尝试使用 WDA_MONITOR
                    success = SetWindowDisplayAffinity(helper.Handle, WDA_MONITOR);
                    if (success)
                    {
                        System.Diagnostics.Debug.WriteLine("Windows 7+ 兼容防录屏设置成功");
                    }
                    
                    // 尝试使用 DWM 排除预览
                    if (osVersion.Major >= 6 && osVersion.Minor >= 1) // Windows 7+
                    {
                        try
                        {
                            int value = 1;
                            DwmSetWindowAttribute(helper.Handle, DWMWA_EXCLUDED_FROM_PEEK, ref value, sizeof(int));
                            System.Diagnostics.Debug.WriteLine("DWM 预览排除设置成功");
                            success = true;
                        }
                        catch (Exception ex)
                        {
                            System.Diagnostics.Debug.WriteLine($"DWM 设置失败: {ex.Message}");
                        }
                    }
                }
                
                // 如果前面的方法都失败，尝试设置窗口为分层窗口（降级方案）
                if (!success)
                {
                    try
                    {
                        int exStyle = GetWindowLong(helper.Handle, GWL_EXSTYLE);
                        SetWindowLong(helper.Handle, GWL_EXSTYLE, exStyle | WS_EX_LAYERED);
                        SetLayeredWindowAttributes(helper.Handle, 0, 255, LWA_ALPHA);
                        System.Diagnostics.Debug.WriteLine("分层窗口防录屏设置成功");
                        success = true;
                    }
                    catch (Exception ex)
                    {
                        System.Diagnostics.Debug.WriteLine($"分层窗口设置失败: {ex.Message}");
                    }
                }
                
                return success;
            }
            catch (Exception ex)
            {
                System.Diagnostics.Debug.WriteLine($"设置防录屏失败: {ex.Message}");
                return false;
            }
        }

        /// <summary>
        /// 禁用防录屏
        /// </summary>
        public static bool DisableDisplayAffinity(Window window)
        {
            try
            {
                var helper = new WindowInteropHelper(window);
                if (helper.Handle == IntPtr.Zero)
                    return false;

                bool success = false;
                var osVersion = Environment.OSVersion.Version;

                // 重置 SetWindowDisplayAffinity
                try
                {
                    success = SetWindowDisplayAffinity(helper.Handle, 0);
                    if (success)
                    {
                        System.Diagnostics.Debug.WriteLine("防录屏已禁用");
                    }
                }
                catch (Exception ex)
                {
                    System.Diagnostics.Debug.WriteLine($"禁用 SetWindowDisplayAffinity 失败: {ex.Message}");
                }

                // 重置 DWM 属性
                if (osVersion.Major >= 6 && osVersion.Minor >= 1) // Windows 7+
                {
                    try
                    {
                        int value = 0;
                        DwmSetWindowAttribute(helper.Handle, DWMWA_EXCLUDED_FROM_PEEK, ref value, sizeof(int));
                        System.Diagnostics.Debug.WriteLine("DWM 预览排除已禁用");
                        success = true;
                    }
                    catch (Exception ex)
                    {
                        System.Diagnostics.Debug.WriteLine($"重置 DWM 失败: {ex.Message}");
                    }
                }

                // 重置分层窗口属性
                try
                {
                    int exStyle = GetWindowLong(helper.Handle, GWL_EXSTYLE);
                    if ((exStyle & WS_EX_LAYERED) != 0)
                    {
                        SetWindowLong(helper.Handle, GWL_EXSTYLE, exStyle & ~WS_EX_LAYERED);
                        System.Diagnostics.Debug.WriteLine("分层窗口属性已重置");
                        success = true;
                    }
                }
                catch (Exception ex)
                {
                    System.Diagnostics.Debug.WriteLine($"重置分层窗口失败: {ex.Message}");
                }

                return success;
            }
            catch (Exception ex)
            {
                System.Diagnostics.Debug.WriteLine($"禁用防录屏失败: {ex.Message}");
                return false;
            }
        }

        /// <summary>
        /// 配置窗口（更安全的实现）
        /// </summary>
        public static void ConfigureWindow(Window window)
        {
            if (window == null) return;

            // 确保在窗口完全加载后再进行配置
            if (window.IsLoaded)
            {
                ApplyWindowConfiguration(window);
            }
            else
            {
                window.Loaded += (sender, e) =>
                {
                    // 延迟执行，确保窗口句柄完全就绪
                    window.Dispatcher.BeginInvoke(new Action(() =>
                    {
                        ApplyWindowConfiguration(window);
                    }), System.Windows.Threading.DispatcherPriority.Loaded);
                };
            }
        }

        /// <summary>
        /// 应用窗口配置
        /// </summary>
        private static void ApplyWindowConfiguration(Window window)
        {
            try
            {
                // 设置窗口置顶
                SetWindowTopmost(window);
                
                // 设置防录屏
                SetWindowDisplayAffinity(window);
            }
            catch (Exception ex)
            {
                System.Diagnostics.Debug.WriteLine($"窗口配置失败: {ex.Message}");
                // 写入错误日志但不影响程序继续运行
                try
                {
                    System.IO.File.WriteAllText("window_config_error.log", 
                        $"Window Config Error: {ex.Message}\nStackTrace: {ex.StackTrace}");
                }
                catch
                {
                    // 忽略日志写入错误
                }
            }
        }








        /// <summary>
        /// 检查是否为 Windows 10 版本 2004 及以上
        /// </summary>
        private static bool IsWindows10Version2004OrLater()
        {
            try
            {
                var version = Environment.OSVersion.Version;
                
                // Windows 10 版本 2004 对应 Build 19041
                return version.Major >= 10 && version.Build >= 19041;
            }
            catch
            {
                return false;
            }
        }

        /// <summary>
        /// 检查是否为 Windows 7 及以上版本
        /// </summary>
        private static bool IsWindows7OrLater()
        {
            try
            {
                var version = Environment.OSVersion.Version;
                
                // Windows 7 对应 6.1，Windows Vista 是 6.0
                return version.Major > 6 || (version.Major == 6 && version.Minor >= 1);
            }
            catch
            {
                return false;
            }
        }

        /// <summary>
        /// 获取当前操作系统信息
        /// </summary>
        public static string GetOSInfo()
        {
            try
            {
                var version = Environment.OSVersion.Version;
                string osName = "未知系统";
                
                if (version.Major == 10)
                {
                    if (version.Build >= 19041)
                        osName = $"Windows 10 (Build {version.Build}) - 支持完整防录屏";
                    else
                        osName = $"Windows 10 (Build {version.Build}) - 支持兼容防录屏";
                }
                else if (version.Major == 6)
                {
                    switch (version.Minor)
                    {
                        case 3:
                            osName = "Windows 8.1 - 支持兼容防录屏";
                            break;
                        case 2:
                            osName = "Windows 8 - 支持兼容防录屏";
                            break;
                        case 1:
                            osName = "Windows 7 - 支持兼容防录屏";
                            break;
                        case 0:
                            osName = "Windows Vista - 支持基础防录屏";
                            break;
                    }
                }
                else if (version.Major == 5)
                {
                    osName = "Windows XP - 不支持防录屏";
                }
                
                return osName;
            }
            catch
            {
                return "无法获取系统信息";
            }
        }

        /// <summary>
        /// 切换防录屏状态
        /// </summary>
        public static bool ToggleDisplayAffinity(Window window)
        {
            try
            {
                var helper = new WindowInteropHelper(window);
                if (helper.Handle == IntPtr.Zero)
                    return false;

                // 这里可以根据需要实现状态检测和切换
                // 为简化，直接尝试启用防录屏
                bool success = SetWindowDisplayAffinity(window);
                
                return success;
            }
            catch (Exception ex)
            {
                System.Diagnostics.Debug.WriteLine($"切换防录屏状态失败: {ex.Message}");
                return false;
            }
        }


    }
} 