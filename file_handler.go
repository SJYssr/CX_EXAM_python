package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileHandler 文件处理结构体
type FileHandler struct {
	tikuPath    string
	cachePath   string
	config      *Config
	fileWatcher *FileWatcher
}

// FileWatcher 文件监视器
type FileWatcher struct {
	tikuModTime time.Time
	tikuHash    string
}

// NewFileHandler 创建文件处理实例
func NewFileHandler() (*FileHandler, error) {
	fh := &FileHandler{
		tikuPath:  "tiku.txt",
		cachePath: filepath.Join("cache", "tiku_cache.txt"),
		fileWatcher: &FileWatcher{},
	}
	
	// 创建缓存目录
	cacheDir := filepath.Dir(fh.cachePath)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("创建缓存目录失败: %v", err)
	}
	
	return fh, nil
}

// LoadTikuContent 加载题库内容
func (fh *FileHandler) LoadTikuContent() (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(fh.tikuPath); os.IsNotExist(err) {
		// 如果题库文件不存在，创建默认文件
		return fh.createDefaultTiku()
	}
	
	// 检查文件是否有更新
	if fh.shouldUseCache() {
		if content, err := fh.loadFromCache(); err == nil {
			return content, nil
		}
	}
	
	// 读取原始文件
	content, err := fh.readTikuFile()
	if err != nil {
		return "", err
	}
	
	// 处理内容
	processedContent := fh.processContent(content)
	
	// 保存到缓存
	fh.saveToCache(processedContent)
	
	// 更新文件监视器
	fh.updateFileWatcher()
	
	return processedContent, nil
}

// createDefaultTiku 创建默认题库文件
func (fh *FileHandler) createDefaultTiku() (string, error) {
	defaultContent := `# 超星学习通题库

## 使用说明
1. 在此文件中添加题目和答案
2. 支持搜索功能
3. 支持AI问答

## 示例题目

问题：计算机网络的七层模型是什么？
答案：OSI七层模型包括：物理层、数据链路层、网络层、传输层、会话层、表示层、应用层

问题：什么是TCP/IP协议？
答案：TCP/IP协议是一组通信协议，包括传输控制协议(TCP)和网际协议(IP)，是互联网的基础协议

问题：数据库的ACID特性是什么？
答案：原子性(Atomicity)、一致性(Consistency)、隔离性(Isolation)、持久性(Durability)

## 快捷键说明
- F3: 切换窗口大小
- F4: 刷新输入法保护
- ESC: 退出程序
- 右键: 改变透明度
- Ctrl+滚轮: 调整透明度
- Alt+滚轮: 调整字体大小

## 注意事项
- 请将题目和答案分别放在"问题："和"答案："后面
- 支持中英文搜索
- 程序会自动检测超星学习通窗口位置
`

	// 写入默认内容
	if err := ioutil.WriteFile(fh.tikuPath, []byte(defaultContent), 0644); err != nil {
		return "", fmt.Errorf("创建默认题库文件失败: %v", err)
	}
	
	return defaultContent, nil
}

// readTikuFile 读取题库文件
func (fh *FileHandler) readTikuFile() (string, error) {
	// 尝试以UTF-8编码读取
	content, err := ioutil.ReadFile(fh.tikuPath)
	if err != nil {
		return "", fmt.Errorf("读取题库文件失败: %v", err)
	}
	
	// 检查并处理BOM
	contentStr := string(content)
	if len(contentStr) >= 3 && contentStr[:3] == "\ufeff" {
		contentStr = contentStr[3:] // 移除UTF-8 BOM
	}
	
	return contentStr, nil
}

// processContent 处理文件内容
func (fh *FileHandler) processContent(content string) string {
	lines := strings.Split(content, "\n")
	var processedLines []string
	
	for i, line := range lines {
		// 移除行尾空白
		line = strings.TrimRight(line, " \t\r")
		
		// 处理特殊字符
		line = strings.ReplaceAll(line, "\u00a0", " ") // 替换不间断空格
		
		// 添加行号（可选）
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
			processedLines = append(processedLines, fmt.Sprintf("%s", line))
		} else {
			processedLines = append(processedLines, line)
		}
		
		// 每1000行打印一次进度
		if i > 0 && i%1000 == 0 {
			fmt.Printf("已处理 %d 行...\n", i)
		}
	}
	
	return strings.Join(processedLines, "\n")
}

// shouldUseCache 检查是否应该使用缓存
func (fh *FileHandler) shouldUseCache() bool {
	// 获取原始文件信息
	fileInfo, err := os.Stat(fh.tikuPath)
	if err != nil {
		return false
	}
	
	// 检查缓存文件是否存在
	if _, err := os.Stat(fh.cachePath); os.IsNotExist(err) {
		return false
	}
	
	// 比较修改时间
	if !fileInfo.ModTime().Equal(fh.fileWatcher.tikuModTime) {
		return false
	}
	
	// 比较文件哈希
	currentHash, err := fh.calculateFileHash(fh.tikuPath)
	if err != nil || currentHash != fh.fileWatcher.tikuHash {
		return false
	}
	
	return true
}

// loadFromCache 从缓存加载内容
func (fh *FileHandler) loadFromCache() (string, error) {
	content, err := ioutil.ReadFile(fh.cachePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// saveToCache 保存内容到缓存
func (fh *FileHandler) saveToCache(content string) error {
	return ioutil.WriteFile(fh.cachePath, []byte(content), 0644)
}

// updateFileWatcher 更新文件监视器
func (fh *FileHandler) updateFileWatcher() error {
	fileInfo, err := os.Stat(fh.tikuPath)
	if err != nil {
		return err
	}
	
	hash, err := fh.calculateFileHash(fh.tikuPath)
	if err != nil {
		return err
	}
	
	fh.fileWatcher.tikuModTime = fileInfo.ModTime()
	fh.fileWatcher.tikuHash = hash
	
	return nil
}

// calculateFileHash 计算文件哈希
func (fh *FileHandler) calculateFileHash(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	
	hash := md5.Sum(content)
	return hex.EncodeToString(hash[:]), nil
}

// SearchInContent 在内容中搜索
func (fh *FileHandler) SearchInContent(content, searchTerm string) []SearchResult {
	var results []SearchResult
	
	lines := strings.Split(content, "\n")
	searchLower := strings.ToLower(searchTerm)
	
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), searchLower) {
			result := SearchResult{
				LineNumber: i + 1,
				Line:       line,
				Context:    fh.getContext(lines, i, 2), // 获取前后2行作为上下文
			}
			results = append(results, result)
		}
	}
	
	return results
}

// SearchResult 搜索结果结构
type SearchResult struct {
	LineNumber int
	Line       string
	Context    []string
}

// getContext 获取上下文行
func (fh *FileHandler) getContext(lines []string, currentIndex, contextLines int) []string {
	var context []string
	
	start := currentIndex - contextLines
	if start < 0 {
		start = 0
	}
	
	end := currentIndex + contextLines + 1
	if end > len(lines) {
		end = len(lines)
	}
	
	for i := start; i < end; i++ {
		prefix := "  "
		if i == currentIndex {
			prefix = "→ "
		}
		context = append(context, fmt.Sprintf("%s%d: %s", prefix, i+1, lines[i]))
	}
	
	return context
}

// ExportToFile 导出内容到文件
func (fh *FileHandler) ExportToFile(content, filePath string) error {
	return ioutil.WriteFile(filePath, []byte(content), 0644)
}

// BackupTiku 备份题库文件
func (fh *FileHandler) BackupTiku() error {
	backupDir := "backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("创建备份目录失败: %v", err)
	}
	
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("tiku_backup_%s.txt", timestamp))
	
	content, err := ioutil.ReadFile(fh.tikuPath)
	if err != nil {
		return fmt.Errorf("读取原文件失败: %v", err)
	}
	
	if err := ioutil.WriteFile(backupPath, content, 0644); err != nil {
		return fmt.Errorf("写入备份文件失败: %v", err)
	}
	
	fmt.Printf("题库已备份到: %s\n", backupPath)
	return nil
}

// ValidateContent 验证内容格式
func (fh *FileHandler) ValidateContent(content string) []string {
	var warnings []string
	lines := strings.Split(content, "\n")
	
	questionCount := 0
	answerCount := 0
	
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		if strings.HasPrefix(trimmedLine, "问题：") {
			questionCount++
			if len(trimmedLine) <= 3 {
				warnings = append(warnings, fmt.Sprintf("第%d行：问题内容为空", i+1))
			}
		}
		
		if strings.HasPrefix(trimmedLine, "答案：") {
			answerCount++
			if len(trimmedLine) <= 3 {
				warnings = append(warnings, fmt.Sprintf("第%d行：答案内容为空", i+1))
			}
		}
		
		// 检查行长度
		if len(line) > 200 {
			warnings = append(warnings, fmt.Sprintf("第%d行：内容过长(>200字符)", i+1))
		}
	}
	
	if questionCount != answerCount {
		warnings = append(warnings, fmt.Sprintf("问题数量(%d)与答案数量(%d)不匹配", questionCount, answerCount))
	}
	
	return warnings
}

// GetFileStats 获取文件统计信息
func (fh *FileHandler) GetFileStats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	fileInfo, err := os.Stat(fh.tikuPath)
	if err == nil {
		stats["file_size"] = fileInfo.Size()
		stats["last_modified"] = fileInfo.ModTime()
	}
	
	content, err := fh.LoadTikuContent()
	if err == nil {
		lines := strings.Split(content, "\n")
		stats["total_lines"] = len(lines)
		
		questionCount := 0
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "问题：") {
				questionCount++
			}
		}
		stats["question_count"] = questionCount
	}
	
	// 检查缓存状态
	if _, err := os.Stat(fh.cachePath); err == nil {
		stats["cache_exists"] = true
	} else {
		stats["cache_exists"] = false
	}
	
	return stats
} 