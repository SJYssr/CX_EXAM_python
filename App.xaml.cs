using System;
using System.IO;
using System.Windows;

namespace CxExamClient
{
    public partial class App : Application
    {
        protected override void OnStartup(StartupEventArgs e)
        {
            base.OnStartup(e);
            
            // 检查必要文件
            try
            {
                ConfigManager.CheckRequiredFiles();
            }
            catch (FileNotFoundException ex)
            {
                MessageBox.Show(ex.Message, "文件检查", MessageBoxButton.OK, MessageBoxImage.Warning);
                Environment.Exit(1);
            }
        }
    }
} 