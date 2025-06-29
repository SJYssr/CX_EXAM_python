using System;
using System.Runtime.InteropServices;
using System.Threading.Tasks;
using System.Windows;

namespace CxExamClient.Services
{
    public static class InputService
    {
        [DllImport("user32.dll")]
        private static extern uint SendInput(uint nInputs, INPUT[] pInputs, int cbSize);

        [DllImport("user32.dll")]
        private static extern IntPtr GetForegroundWindow();

        [DllImport("user32.dll")]
        private static extern int GetWindowText(IntPtr hWnd, System.Text.StringBuilder text, int count);

        // SendInput 结构体定义
        [StructLayout(LayoutKind.Sequential)]
        struct INPUT
        {
            public uint type;
            public INPUTUNION union;
        }

        [StructLayout(LayoutKind.Explicit)]
        struct INPUTUNION
        {
            [FieldOffset(0)]
            public MOUSEINPUT mouseInput;
            [FieldOffset(0)]
            public KEYBDINPUT keyboardInput;
            [FieldOffset(0)]
            public HARDWAREINPUT hardwareInput;
        }

        [StructLayout(LayoutKind.Sequential)]
        struct MOUSEINPUT
        {
            public int dx;
            public int dy;
            public uint mouseData;
            public uint dwFlags;
            public uint time;
            public IntPtr dwExtraInfo;
        }

        [StructLayout(LayoutKind.Sequential)]
        struct KEYBDINPUT
        {
            public ushort wVk;
            public ushort wScan;
            public uint dwFlags;
            public uint time;
            public IntPtr dwExtraInfo;
        }

        [StructLayout(LayoutKind.Sequential)]
        struct HARDWAREINPUT
        {
            public uint uMsg;
            public ushort wParamL;
            public ushort wParamH;
        }

        private const uint INPUT_KEYBOARD = 1;
        private const uint KEYEVENTF_KEYDOWN = 0x0000;
        private const uint KEYEVENTF_KEYUP = 0x0002;
        private const uint KEYEVENTF_UNICODE = 0x0004;

        private static Random _random = new Random();

        public static async Task SendTextWithDelay(string text, int delaySeconds = 5)
        {
            try
            {
                // 获取当前窗口信息
                IntPtr foregroundWindow = GetForegroundWindow();
                var windowTitle = new System.Text.StringBuilder(256);
                GetWindowText(foregroundWindow, windowTitle, windowTitle.Capacity);
                
                // 显示开始消息
                MessageBox.Show($"即将模拟打字输入: '{text}'\n" +
                              $"目标窗口: {windowTitle}\n" +
                              $"请在{delaySeconds}秒内切换到输入框", 
                              "自动打字", MessageBoxButton.OK, MessageBoxImage.Information);
                
                await Task.Delay(delaySeconds * 1000);
                
                // 逐字符输入文本
                await TypeTextLikeHuman(text);
                
                MessageBox.Show($"打字完成！已输入 {text.Length} 个字符", "完成", 
                              MessageBoxButton.OK, MessageBoxImage.Information);
            }
            catch (Exception ex)
            {
                MessageBox.Show($"自动打字异常: {ex.Message}", "错误", 
                              MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }

        private static async Task TypeTextLikeHuman(string text)
        {
            foreach (char c in text)
            {
                // 发送字符
                SendChar(c);
                
                // 随机延迟模拟人工打字速度 (50-200ms)
                int delay = _random.Next(50, 200);
                await Task.Delay(delay);
            }
        }

        private static void SendChar(char c)
        {
            INPUT[] inputs = new INPUT[1];
            
            inputs[0] = new INPUT
            {
                type = INPUT_KEYBOARD,
                union = new INPUTUNION
                {
                    keyboardInput = new KEYBDINPUT
                    {
                        wVk = 0,
                        wScan = c,
                        dwFlags = KEYEVENTF_UNICODE,
                        time = 0,
                        dwExtraInfo = IntPtr.Zero
                    }
                }
            };

            SendInput((uint)inputs.Length, inputs, Marshal.SizeOf(typeof(INPUT)));
        }
    }
} 