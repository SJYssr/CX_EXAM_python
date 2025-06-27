using System;
using System.Runtime.InteropServices;
using System.Threading.Tasks;
using System.Windows;

namespace CxExamClient.Services
{
    public static class InputService
    {
        [DllImport("user32.dll")]
        private static extern void keybd_event(byte bVk, byte bScan, uint dwFlags, UIntPtr dwExtraInfo);

        private const uint KEYEVENTF_KEYUP = 0x0002;
        private const byte VK_CONTROL = 0x11;
        private const byte VK_V = 0x56;

        public static async Task SendTextWithDelay(string text, int delaySeconds = 5)
        {
            await Task.Delay(delaySeconds * 1000);
            
            // 保存当前剪贴板内容
            string? originalClipboard = null;
            try
            {
                if (Clipboard.ContainsText())
                {
                    originalClipboard = Clipboard.GetText();
                }
            }
            catch { }

            try
            {
                // 将文本复制到剪贴板
                Clipboard.SetText(text);
                await Task.Delay(100);

                // 发送 Ctrl+V 粘贴
                SendKey(VK_CONTROL, true);
                SendKey(VK_V, false);
                SendKey(VK_CONTROL, false);
            }
            finally
            {
                // 恢复原剪贴板内容
                await Task.Delay(500);
                try
                {
                    if (!string.IsNullOrEmpty(originalClipboard))
                    {
                        Clipboard.SetText(originalClipboard);
                    }
                    else
                    {
                        Clipboard.Clear();
                    }
                }
                catch { }
            }
        }

        private static void SendKey(byte vk, bool isDown)
        {
            if (isDown)
            {
                keybd_event(vk, 0, 0, UIntPtr.Zero);
            }
            else
            {
                keybd_event(vk, 0, KEYEVENTF_KEYUP, UIntPtr.Zero);
            }
            System.Threading.Thread.Sleep(50);
        }
    }
} 