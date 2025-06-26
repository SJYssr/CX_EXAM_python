package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config 配置结构体
type Config struct {
	AIType        int           `json:"ai_type"`
	SparkConfig   *SparkConfig  `json:"spark_config,omitempty"`
	DeepseekConfig *DeepseekConfig `json:"deepseek_config,omitempty"`
	configPath    string
	encryptionKey []byte
}

// SparkConfig 讯飞星火配置
type SparkConfig struct {
	AppID     string `json:"appid"`
	APISecret string `json:"api_secret"`
	APIKey    string `json:"api_key"`
	SparkURL  string `json:"spark_url"`
	Domain    string `json:"domain"`
}

// DeepseekConfig Deepseek配置
type DeepseekConfig struct {
	APIKey string `json:"api_key"`
	Model  string `json:"model"`
}

// NewConfig 创建新的配置实例
func NewConfig() (*Config, error) {
	configDir := filepath.Join(".", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %v", err)
	}
	
	config := &Config{
		configPath: filepath.Join(configDir, "config.json.enc"),
		AIType:     2, // 默认使用Deepseek
	}
	
	// 生成或加载加密密钥
	if err := config.initEncryptionKey(); err != nil {
		return nil, fmt.Errorf("初始化加密密钥失败: %v", err)
	}
	
	// 尝试加载现有配置
	if err := config.LoadConfig(); err != nil {
		// 如果加载失败，创建默认配置
		config.setDefaultConfig()
		if err := config.SaveConfig(); err != nil {
			return nil, fmt.Errorf("保存默认配置失败: %v", err)
		}
	}
	
	return config, nil
}

// initEncryptionKey 初始化加密密钥
func (c *Config) initEncryptionKey() error {
	keyPath := filepath.Join(filepath.Dir(c.configPath), "key.dat")
	
	// 尝试读取现有密钥
	if keyData, err := ioutil.ReadFile(keyPath); err == nil {
		c.encryptionKey = keyData
		return nil
	}
	
	// 生成新密钥
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		return fmt.Errorf("生成随机密钥失败: %v", err)
	}
	
	// 保存密钥
	if err := ioutil.WriteFile(keyPath, key, 0600); err != nil {
		return fmt.Errorf("保存密钥失败: %v", err)
	}
	
	c.encryptionKey = key
	return nil
}

// encrypt 使用AES-GCM加密数据
func (c *Config) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return nil, err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypt 使用AES-GCM解密数据
func (c *Config) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return nil, err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("密文太短")
	}
	
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	
	return plaintext, nil
}

// LoadConfig 加载加密配置
func (c *Config) LoadConfig() error {
	// 读取加密文件
	encryptedData, err := ioutil.ReadFile(c.configPath)
	if err != nil {
		return err
	}
	
	// 解密数据
	decryptedData, err := c.decrypt(encryptedData)
	if err != nil {
		return fmt.Errorf("解密配置失败: %v", err)
	}
	
	// 解析JSON
	if err := json.Unmarshal(decryptedData, c); err != nil {
		return fmt.Errorf("解析配置失败: %v", err)
	}
	
	return nil
}

// SaveConfig 保存加密配置
func (c *Config) SaveConfig() error {
	// 序列化为JSON
	jsonData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}
	
	// 加密数据
	encryptedData, err := c.encrypt(jsonData)
	if err != nil {
		return fmt.Errorf("加密配置失败: %v", err)
	}
	
	// 写入文件
	if err := ioutil.WriteFile(c.configPath, encryptedData, 0600); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	
	return nil
}

// setDefaultConfig 设置默认配置
func (c *Config) setDefaultConfig() {
	c.AIType = 2
	c.DeepseekConfig = &DeepseekConfig{
		APIKey: "", // 需要用户配置
		Model:  "deepseek-chat",
	}
	c.SparkConfig = &SparkConfig{
		AppID:     "", // 需要用户配置
		APISecret: "", // 需要用户配置
		APIKey:    "", // 需要用户配置
		SparkURL:  "wss://spark-api.xf-yun.com/v3.5/chat",
		Domain:    "generalv3.5",
	}
}

// GetAIType 获取AI类型
func (c *Config) GetAIType() int {
	return c.AIType
}

// GetSparkConfig 获取星火配置
func (c *Config) GetSparkConfig() *SparkConfig {
	return c.SparkConfig
}

// GetDeepseekConfig 获取Deepseek配置
func (c *Config) GetDeepseekConfig() *DeepseekConfig {
	return c.DeepseekConfig
}

// SetAIType 设置AI类型
func (c *Config) SetAIType(aiType int) error {
	c.AIType = aiType
	return c.SaveConfig()
}

// UpdateSparkConfig 更新星火配置
func (c *Config) UpdateSparkConfig(config *SparkConfig) error {
	c.SparkConfig = config
	return c.SaveConfig()
}

// UpdateDeepseekConfig 更新Deepseek配置
func (c *Config) UpdateDeepseekConfig(config *DeepseekConfig) error {
	c.DeepseekConfig = config
	return c.SaveConfig()
}

// EncryptString 加密字符串并返回Base64编码
func (c *Config) EncryptString(plaintext string) (string, error) {
	encrypted, err := c.encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptString 解密Base64编码的字符串
func (c *Config) DecryptString(encryptedBase64 string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", err
	}
	
	decrypted, err := c.decrypt(encrypted)
	if err != nil {
		return "", err
	}
	
	return string(decrypted), nil
}

// HashPassword 使用SHA-256哈希密码
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(hash[:])
} 