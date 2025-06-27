namespace CxExamClient.Models
{
    public class ConfigModel
    {
        public AiSetModel AI_set { get; set; } = new();
        public PathModel PATH { get; set; } = new();
        public SparkModel SPARK { get; set; } = new();
        public DeepseekModel deepseek { get; set; } = new();
    }

    public class AiSetModel
    {
        public int type { get; set; } = 0;
    }

    public class PathModel
    {
        public string path { get; set; } = string.Empty;
    }

    public class SparkModel
    {
        public string Spark_url { get; set; } = string.Empty;
        public string api_key { get; set; } = string.Empty;
        public string api_secret { get; set; } = string.Empty;
        public string appid { get; set; } = string.Empty;
        public string domain { get; set; } = string.Empty;
    }

    public class DeepseekModel
    {
        public string api_key { get; set; } = string.Empty;
        public string model { get; set; } = "deepseek-chat";
    }
} 