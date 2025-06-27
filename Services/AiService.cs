using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Http;
using System.Text;
using System.Security.Cryptography;
using System.Threading.Tasks;
using Newtonsoft.Json;
using WebSocketSharp;
using CxExamClient.Models;

namespace CxExamClient.Services
{
    public class AiService
    {
        private readonly HttpClient _httpClient;
        private readonly ConfigModel _config;

        public AiService()
        {
            _httpClient = new HttpClient();
            _config = ConfigManager.Config;
        }

        public async Task<string> AskAiAsync(string question, Action<string>? onMessage = null)
        {
            try
            {
                switch (_config.AI_set.type)
                {
                    case 1:
                        return await AskSparkAiAsync(question, onMessage);
                    case 2:
                        return await AskDeepseekAiAsync(question);
                    default:
                        return "未配置AI服务";
                }
            }
            catch (Exception ex)
            {
                return $"发生错误: {ex.Message}";
            }
        }

        private async Task<string> AskDeepseekAiAsync(string question)
        {
            const string url = "https://api.deepseek.com/v1/chat/completions";
            
            var requestData = new
            {
                model = _config.deepseek.model,
                messages = new[]
                {
                    new { role = "user", content = question }
                },
                temperature = 0.7
            };

            var jsonContent = JsonConvert.SerializeObject(requestData);
            var content = new StringContent(jsonContent, Encoding.UTF8, "application/json");
            
            _httpClient.DefaultRequestHeaders.Clear();
            _httpClient.DefaultRequestHeaders.Add("Authorization", $"Bearer {_config.deepseek.api_key}");

            var response = await _httpClient.PostAsync(url, content);
            var responseContent = await response.Content.ReadAsStringAsync();

            if (response.IsSuccessStatusCode)
            {
                dynamic result = JsonConvert.DeserializeObject(responseContent)!;
                return result.choices[0].message.content;
            }
            else
            {
                return $"API调用失败: {response.StatusCode} - {responseContent}";
            }
        }

        private async Task<string> AskSparkAiAsync(string question, Action<string>? onMessage = null)
        {
            var tcs = new TaskCompletionSource<string>();
            var responseBuilder = new StringBuilder();

            try
            {
                var wsParam = new SparkWebSocketParam(_config.SPARK);
                var wsUrl = wsParam.CreateUrl();

                using var ws = new WebSocket(wsUrl);
                
                ws.OnMessage += (sender, e) =>
                {
                    try
                    {
                        dynamic data = JsonConvert.DeserializeObject(e.Data)!;
                        
                        if (data.header.code != 0)
                        {
                            tcs.SetResult($"请求错误: {data.header.code}");
                            return;
                        }

                        string content = data.payload.choices.text[0].content;
                        responseBuilder.Append(content);
                        onMessage?.Invoke(content);

                        if (data.payload.choices.status == 2)
                        {
                            tcs.SetResult(responseBuilder.ToString());
                        }
                    }
                    catch (Exception ex)
                    {
                        tcs.SetException(ex);
                    }
                };

                ws.OnError += (sender, e) =>
                {
                    tcs.SetResult($"WebSocket错误: {e.Message}");
                };

                ws.OnOpen += (sender, e) =>
                {
                    var requestData = CreateSparkRequestData(_config.SPARK.appid, question, _config.SPARK.domain);
                    ws.Send(JsonConvert.SerializeObject(requestData));
                };

                ws.Connect();
                
                return await tcs.Task;
            }
            catch (Exception ex)
            {
                return $"连接失败: {ex.Message}";
            }
        }

        private object CreateSparkRequestData(string appid, string query, string domain)
        {
            return new
            {
                header = new
                {
                    app_id = appid,
                    uid = "1234"
                },
                parameter = new
                {
                    chat = new
                    {
                        domain,
                        temperature = 0.5,
                        max_tokens = 4096,
                        auditing = "default"
                    }
                },
                payload = new
                {
                    message = new
                    {
                        text = new[]
                        {
                            new { role = "user", content = query }
                        }
                    }
                }
            };
        }

        public void Dispose()
        {
            _httpClient?.Dispose();
        }
    }

    public class SparkWebSocketParam
    {
        private readonly SparkModel _config;

        public SparkWebSocketParam(SparkModel config)
        {
            _config = config;
        }

        public string CreateUrl()
        {
            var uri = new Uri(_config.Spark_url);
            var host = uri.Host;
            var path = uri.PathAndQuery;
            
            var now = DateTime.UtcNow;
            var date = now.ToString("r");

            var signatureOrigin = $"host: {host}\ndate: {date}\nGET {path} HTTP/1.1";
            
            using var hmac = new HMACSHA256(Encoding.UTF8.GetBytes(_config.api_secret));
            var signatureBytes = hmac.ComputeHash(Encoding.UTF8.GetBytes(signatureOrigin));
            var signatureBase64 = Convert.ToBase64String(signatureBytes);

            var authorizationOrigin = $"api_key=\"{_config.api_key}\", algorithm=\"hmac-sha256\", headers=\"host date request-line\", signature=\"{signatureBase64}\"";
            var authorization = Convert.ToBase64String(Encoding.UTF8.GetBytes(authorizationOrigin));

            var queryParams = new Dictionary<string, string>
            {
                ["authorization"] = authorization,
                ["date"] = date,
                ["host"] = host
            };

            var queryString = string.Join("&", queryParams.Select(kv => $"{kv.Key}={Uri.EscapeDataString(kv.Value)}"));
            
            return $"{_config.Spark_url}?{queryString}";
        }
    }
} 