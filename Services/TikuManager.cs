using System;
using System.IO;
using System.Text;

namespace CxExamClient.Services
{
    public static class TikuManager
    {
        private static string? _tikuContent;
        
        public static string LoadTikuContent()
        {
            if (_tikuContent != null)
            {
                return _tikuContent;
            }
            
            try
            {
                var tikuPath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "tiku.txt");
                if (File.Exists(tikuPath))
                {
                    _tikuContent = File.ReadAllText(tikuPath, Encoding.UTF8);
                }
                else
                {
                    _tikuContent = "未找到tiku.txt，请确保文件夹中有此文件";
                }
            }
            catch (Exception ex)
            {
                _tikuContent = $"读取题库文件时发生错误: {ex.Message}";
            }
            
            return _tikuContent;
        }

        public static void ReloadTiku()
        {
            _tikuContent = null;
        }
    }
} 