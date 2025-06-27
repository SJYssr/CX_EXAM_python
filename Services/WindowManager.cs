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

        private static readonly IntPtr HWND_TOPMOST = new IntPtr(-1);
        private const uint SWP_NOMOVE = 0x0002;
        private const uint SWP_NOSIZE = 0x0001;
        private const uint WDA_EXCLUDEFROMCAPTURE = 0x00000011;

        public static void SetWindowTopmost(Window window)
        {
            var helper = new WindowInteropHelper(window);
            SetWindowPos(helper.Handle, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE);
        }

        public static void SetWindowDisplayAffinity(Window window)
        {
            var helper = new WindowInteropHelper(window);
            SetWindowDisplayAffinity(helper.Handle, WDA_EXCLUDEFROMCAPTURE);
        }

        public static void ConfigureWindow(Window window)
        {
            window.Loaded += (sender, e) =>
            {
                // 设置窗口置顶
                SetWindowTopmost(window);
                
                // 设置防录屏
                SetWindowDisplayAffinity(window);
            };
        }
    }
} 