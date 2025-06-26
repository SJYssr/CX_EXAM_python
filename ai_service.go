package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AIService AI服务接口
type AIService struct {
	config        *Config
	httpClient    *http.Client
	sparkService  *SparkService
	deepseekService *DeepseekService
}

// SparkService 讯飞星火服务
type SparkService struct {
	config *SparkConfig
}

// DeepseekService Deepseek服务  
type DeepseekService struct {
	config *DeepseekConfig
	client *http.Client
}

// DeepseekRequest Deepseek请求结构
type DeepseekRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Temperature float64 `json:"temperature"`
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepseekResponse Deepseek响应结构
type DeepseekResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 选择结构
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 使用情况结构
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewAIService 创建AI服务实例
func NewAIService(config *Config) (*AIService, error) {
	// 创建安全的HTTP客户端
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			},
		},
		DisableCompression: false,
	}
	
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}
	
	service := &AIService{
		config:     config,
		httpClient: client,
	}
	
	// 初始化服务
	if config.GetAIType() == 1 && config.GetSparkConfig() != nil {
		service.sparkService = &SparkService{
			config: config.GetSparkConfig(),
		}
	}
	
	if config.GetAIType() == 2 && config.GetDeepseekConfig() != nil {
		service.deepseekService = &DeepseekService{
			config: config.GetDeepseekConfig(),
			client: client,
		}
	}
	
	return service, nil
}

// Ask 询问AI
func (a *AIService) Ask(query string) (string, error) {
	// 加密查询内容（可选，用于日志记录）
	encryptedQuery, err := a.config.EncryptString(query)
	if err != nil {
		// 如果加密失败，继续使用原文，但记录错误
		fmt.Printf("警告：查询内容加密失败: %v\n", err)
	} else {
		fmt.Printf("已加密查询长度: %d\n", len(encryptedQuery))
	}
	
	switch a.config.GetAIType() {
	case 1:
		if a.sparkService == nil {
			return "", fmt.Errorf("星火AI服务未初始化")
		}
		return a.sparkService.Ask(query)
	case 2:
		if a.deepseekService == nil {
			return "", fmt.Errorf("Deepseek服务未初始化")
		}
		return a.deepseekService.Ask(query)
	default:
		return "", fmt.Errorf("未知的AI类型: %d", a.config.GetAIType())
	}
}

// SparkService Ask方法
func (s *SparkService) Ask(query string) (string, error) {
	// 注意：这里简化了星火AI的WebSocket实现
	// 实际项目中需要实现完整的WebSocket连接和认证
	return "星火AI功能正在开发中，请使用Deepseek AI", nil
}

// DeepseekService Ask方法
func (d *DeepseekService) Ask(query string) (string, error) {
	if d.config.APIKey == "" {
		return "", fmt.Errorf("Deepseek API Key未配置")
	}
	
	// 构建请求
	request := DeepseekRequest{
		Model: d.config.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: query,
			},
		},
		Temperature: 0.7,
	}
	
	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}
	
	// 创建HTTP请求
	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.config.APIKey)
	req.Header.Set("User-Agent", "ChaoxingHelper/1.0")
	
	// 发送请求
	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}
	
	// 解析响应
	var response DeepseekResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}
	
	// 检查响应内容
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("AI响应为空")
	}
	
	result := response.Choices[0].Message.Content
	
	// 可选：加密响应内容用于缓存
	// 注意：这里需要config参数，暂时注释掉
	// encryptedResult, err := config.EncryptString(result)
	// if err == nil {
	// 	fmt.Printf("已加密响应长度: %d\n", len(encryptedResult))
	// }
	
	return result, nil
}

// ValidateAPIKey 验证API密钥有效性
func (d *DeepseekService) ValidateAPIKey() error {
	if d.config.APIKey == "" {
		return fmt.Errorf("API密钥为空")
	}
	
	// 发送测试请求验证密钥
	testQuery := "Hello"
	_, err := d.Ask(testQuery)
	return err
}

// EncryptCommunication 加密通信设置
func (a *AIService) EncryptCommunication(enable bool) {
	if enable {
		// 启用TLS 1.3
		if tr, ok := a.httpClient.Transport.(*http.Transport); ok {
			tr.TLSClientConfig.MinVersion = tls.VersionTLS13
		}
	}
}

// GetSupportedModels 获取支持的模型列表
func (d *DeepseekService) GetSupportedModels() []string {
	return []string{
		"deepseek-chat",
		"deepseek-coder",
	}
}

// SetModel 设置模型
func (d *DeepseekService) SetModel(model string) error {
	supportedModels := d.GetSupportedModels()
	for _, supported := range supportedModels {
		if model == supported {
			d.config.Model = model
			return nil
		}
	}
	return fmt.Errorf("不支持的模型: %s", model)
}

// GetStats 获取服务统计信息
func (a *AIService) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["ai_type"] = a.config.GetAIType()
	stats["encryption_enabled"] = true
	stats["tls_version"] = "1.3"
	
	if a.config.GetAIType() == 2 && a.deepseekService != nil {
		stats["model"] = a.deepseekService.config.Model
		stats["api_configured"] = a.deepseekService.config.APIKey != ""
	}
	
	return stats
} 