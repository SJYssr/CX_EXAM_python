using System;
using System.IO;
using YamlDotNet.Serialization;
using YamlDotNet.Serialization.NamingConventions;
using CxExamClient.Models;

namespace CxExamClient
{
    public static class ConfigManager
    {
        private static ConfigModel? _config;
        
        public static ConfigModel Config
        {
            get
            {
                if (_config == null)
                {
                    LoadConfig();
                }
                return _config!;
            }
        }

        public static void CheckRequiredFiles()
        {
            var configPath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "config.yaml");
            if (!File.Exists(configPath))
            {
                throw new FileNotFoundException("缺少config.yaml文件，请检查");
            }
        }

        private static void LoadConfig()
        {
            try
            {
                var configPath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "config.yaml");
                var yaml = File.ReadAllText(configPath, System.Text.Encoding.UTF8);
                
                var deserializer = new DeserializerBuilder()
                    .Build();
                
                _config = deserializer.Deserialize<ConfigModel>(yaml);
                
                if (_config == null)
                {
                    _config = new ConfigModel();
                }
            }
            catch (Exception ex)
            {
                throw new InvalidOperationException($"加载配置文件失败: {ex.Message}", ex);
            }
        }

        public static void ReloadConfig()
        {
            _config = null;
            LoadConfig();
        }
    }
} 