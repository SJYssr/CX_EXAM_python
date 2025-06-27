package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Windows剪贴板相关
var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procSetClipboard = user32.NewProc("SetClipboardData")
	procOpenClipboard = user32.NewProc("OpenClipboard")
	procCloseClipboard = user32.NewProc("CloseClipboard")
	procEmptyClipboard = user32.NewProc("EmptyClipboard")
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procGlobalAlloc  = kernel32.NewProc("GlobalAlloc")
	procGlobalLock   = kernel32.NewProc("GlobalLock")
	procGlobalUnlock = kernel32.NewProc("GlobalUnlock")
)

const (
	CF_UNICODETEXT = 13
	GMEM_MOVEABLE  = 0x0002
)

type App struct {
	questions []string
}

var app *App

// 智能读取文件
func smartReadFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// 移除BOM
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}

	// UTF-8检查
	if utf8.Valid(data) {
		return string(data), nil
	}

	// GBK转换
	gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	utf8Data, _, err := transform.Bytes(gbkDecoder, data)
	if err == nil {
		return string(utf8Data), nil
	}

	// GB18030转换
	gb18030Decoder := simplifiedchinese.GB18030.NewDecoder()
	utf8Data, _, err = transform.Bytes(gb18030Decoder, data)
	if err == nil {
		return string(utf8Data), nil
	}

	return strings.ToValidUTF8(string(data), ""), nil
}

// Windows剪贴板操作
func setClipboardText(text string) error {
	ret, _, _ := procOpenClipboard.Call(0)
	if ret == 0 {
		return fmt.Errorf("无法打开剪贴板")
	}
	defer procCloseClipboard.Call()

	procEmptyClipboard.Call()

	utf16Text := syscall.StringToUTF16(text)
	textSize := len(utf16Text) * 2

	hMem, _, _ := procGlobalAlloc.Call(GMEM_MOVEABLE, uintptr(textSize))
	if hMem == 0 {
		return fmt.Errorf("内存分配失败")
	}

	pMem, _, _ := procGlobalLock.Call(hMem)
	if pMem == 0 {
		return fmt.Errorf("内存锁定失败")
	}

	copy((*[1 << 30]uint16)(unsafe.Pointer(pMem))[:len(utf16Text)], utf16Text)
	procGlobalUnlock.Call(hMem)

	ret, _, _ = procSetClipboard.Call(CF_UNICODETEXT, hMem)
	if ret == 0 {
		return fmt.Errorf("剪贴板设置失败")
	}

	return nil
}

func (a *App) search(keyword string) []string {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return nil
	}

	var results []string
	for i, q := range a.questions {
		if strings.Contains(strings.ToLower(q), keyword) {
			results = append(results, fmt.Sprintf("[%d] %s", i+1, q))
			if len(results) >= 20 {
				results = append(results, "... (最多显示20条结果)")
				break
			}
		}
	}
	return results
}

// 启动浏览器
func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("不支持的操作系统")
	}
	if err != nil {
		fmt.Println("无法自动打开浏览器:", err)
		fmt.Printf("请手动打开浏览器访问: %s\n", url)
	}
}

// HTML页面
const htmlPage = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>超星学习通辅助工具</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: "Microsoft YaHei", "微软雅黑", "SimSun", "宋体", Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container {
            max-width: 900px;
            margin: 0 auto;
            background: white;
            border-radius: 15px;
            padding: 30px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
        }
        h1 {
            text-align: center;
            color: #333;
            margin-bottom: 20px;
            font-size: 28px;
        }
        .status {
            background: #e8f5e8;
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 25px;
            border-left: 4px solid #4caf50;
            font-weight: bold;
        }
        .search-section {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 10px;
            margin-bottom: 25px;
        }
        .search-box {
            display: flex;
            gap: 10px;
            margin-bottom: 15px;
        }
        .search-input {
            flex: 1;
            padding: 12px;
            border: 2px solid #ddd;
            border-radius: 6px;
            font-size: 16px;
            font-family: inherit;
        }
        .search-input:focus {
            outline: none;
            border-color: #667eea;
        }
        .btn {
            padding: 12px 20px;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            font-weight: bold;
            transition: all 0.3s;
            color: white;
        }
        .btn-primary { background: #667eea; }
        .btn-primary:hover { background: #5a67d8; }
        .btn-success { background: #28a745; }
        .btn-success:hover { background: #218838; }
        .btn-info { background: #17a2b8; }
        .btn-info:hover { background: #138496; }
        .btn-warning { background: #ffc107; color: #212529; }
        .btn-warning:hover { background: #e0a800; }
        .btn-secondary { background: #6c757d; }
        .btn-secondary:hover { background: #5a6268; }
        .button-group {
            display: flex;
            gap: 10px;
            margin-bottom: 25px;
            flex-wrap: wrap;
        }
        .results-section {
            background: #f8f9fa;
            border-radius: 10px;
            padding: 20px;
        }
        .results-area {
            width: 100%;
            min-height: 300px;
            padding: 15px;
            border: 2px solid #ddd;
            border-radius: 6px;
            font-family: "Consolas", "微软雅黑", monospace;
            font-size: 14px;
            line-height: 1.6;
            resize: vertical;
            background: white;
        }
        .results-area:focus {
            outline: none;
            border-color: #667eea;
        }
        .loading {
            display: none;
            text-align: center;
            color: #667eea;
            font-weight: bold;
            margin-bottom: 10px;
        }
        .success {
            background: #d4edda;
            color: #155724;
            padding: 10px;
            border-radius: 6px;
            margin-bottom: 10px;
            border-left: 4px solid #28a745;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎓 超星学习通辅助工具</h1>
        <div class="status" id="status">正在加载题库...</div>
        
        <div class="search-section">
            <h3>🔍 搜索功能</h3>
            <div class="search-box">
                <input type="text" id="searchInput" class="search-input" placeholder="请输入中文或英文关键词搜索..." />
                <button class="btn btn-primary" onclick="searchLocal()">本地搜索</button>
            </div>
        </div>
        
        <div class="button-group">
            <button class="btn btn-success" onclick="copyToClipboard()">📋 复制到剪贴板</button>
            <button class="btn btn-info" onclick="testChinese()">🧪 测试中文</button>
            <button class="btn btn-warning" onclick="showQuestions()">📄 显示题库</button>
            <button class="btn btn-secondary" onclick="clearResults()">🗑️ 清空</button>
        </div>
        
        <div class="results-section">
            <h3>📋 搜索结果</h3>
            <div class="loading" id="loading">搜索中，请稍候...</div>
            <textarea id="resultsArea" class="results-area" placeholder="搜索结果将在这里显示...

💡 提示：
1. 输入关键词进行本地搜索
2. 支持中文和英文搜索
3. 点击复制按钮将结果复制到剪贴板
4. 在考试窗口按 Ctrl+V 粘贴"></textarea>
        </div>
    </div>

    <script>
        let totalQuestions = 0;
        
        // 页面加载时初始化
        window.addEventListener('DOMContentLoaded', function() {
            loadStatus();
        });
        
        // 加载状态
        async function loadStatus() {
            try {
                const response = await fetch('/api/status');
                const result = await response.json();
                totalQuestions = result.total;
                document.getElementById('status').textContent = 
                    '✅ 题库已加载: ' + totalQuestions + ' 条 | UTF-8编码支持 | 剪贴板100%兼容';
            } catch (error) {
                document.getElementById('status').textContent = '❌ 加载失败: ' + error.message;
            }
        }
        
        // 本地搜索
        async function searchLocal() {
            const keyword = document.getElementById('searchInput').value.trim();
            if (!keyword) {
                alert('请输入搜索关键词');
                return;
            }
            
            showLoading(true);
            try {
                const response = await fetch('/api/search', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({keyword: keyword})
                });
                const result = await response.json();
                
                if (result.results && result.results.length > 0) {
                    document.getElementById('resultsArea').value = result.results.join('\\n\\n');
                } else {
                    document.getElementById('resultsArea').value = '未找到相关题目，请尝试其他关键词';
                }
            } catch (error) {
                document.getElementById('resultsArea').value = '搜索失败: ' + error.message;
            } finally {
                showLoading(false);
            }
        }
        
        // 复制到剪贴板
        async function copyToClipboard() {
            const content = document.getElementById('resultsArea').value;
            if (!content.trim()) {
                alert('没有内容可复制');
                return;
            }
            
            try {
                const response = await fetch('/api/copy', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({text: content})
                });
                const result = await response.json();
                
                if (result.success) {
                    showSuccess('内容已复制到剪贴板！请在考试窗口按 Ctrl+V 粘贴');
                } else {
                    alert('复制失败: ' + result.message);
                }
            } catch (error) {
                alert('复制失败: ' + error.message);
            }
        }
        
        // 测试中文
        async function testChinese() {
            const testText = '测试成功！中文显示正常：\\n数据库管理、算法设计、程序开发、软件工程\\n计算机网络、操作系统、数据结构、人工智能';
            document.getElementById('resultsArea').value = testText;
            
            try {
                const response = await fetch('/api/copy', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({text: testText})
                });
                const result = await response.json();
                
                if (result.success) {
                    showSuccess('中文测试文本已复制到剪贴板！请粘贴验证功能是否正常');
                } else {
                    alert('复制失败: ' + result.message);
                }
            } catch (error) {
                alert('复制失败: ' + error.message);
            }
        }
        
        // 显示题库
        async function showQuestions() {
            try {
                const response = await fetch('/api/questions');
                const result = await response.json();
                
                document.getElementById('resultsArea').value = 
                    '题库预览（前10题）：\\n\\n' + result.questions.join('\\n\\n');
            } catch (error) {
                document.getElementById('resultsArea').value = '获取题库失败: ' + error.message;
            }
        }
        
        // 清空结果
        function clearResults() {
            document.getElementById('searchInput').value = '';
            document.getElementById('resultsArea').value = '';
        }
        
        // 显示加载状态
        function showLoading(show) {
            document.getElementById('loading').style.display = show ? 'block' : 'none';
        }
        
        // 显示成功消息
        function showSuccess(message) {
            const successDiv = document.createElement('div');
            successDiv.className = 'success';
            successDiv.textContent = message;
            document.querySelector('.results-section').insertBefore(successDiv, document.getElementById('resultsArea'));
            setTimeout(() => successDiv.remove(), 3000);
        }
        
        // 回车键搜索
        document.getElementById('searchInput').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                searchLocal();
            }
        });
    </script>
</body>
</html>`

// HTTP处理函数
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlPage))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := map[string]interface{}{
		"total":   len(app.questions),
		"success": true,
	}
	json.NewEncoder(w).Encode(response)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != "POST" {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Keyword string `json:"keyword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}

	results := app.search(req.Keyword)
	response := map[string]interface{}{
		"results": results,
		"count":   len(results),
		"success": true,
	}

	json.NewEncoder(w).Encode(response)
}

func copyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != "POST" {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}

	err := setClipboardText(req.Text)
	response := map[string]interface{}{
		"success": err == nil,
		"message": "内容已复制到剪贴板",
	}
	if err != nil {
		response["message"] = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}

func questionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	count := 10
	if len(app.questions) < 10 {
		count = len(app.questions)
	}

	questions := make([]string, count)
	for i := 0; i < count; i++ {
		questions[i] = fmt.Sprintf("[%d] %s", i+1, app.questions[i])
	}

	response := map[string]interface{}{
		"questions": questions,
		"total":     len(app.questions),
		"success":   true,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	fmt.Println("超星学习通辅助工具 - 终极Webview版本")
	fmt.Println("================================================")
	fmt.Println("🎯 使用浏览器界面，彻底解决中文显示问题")
	fmt.Println("📋 剪贴板功能100%正常")
	fmt.Println("")

	app = &App{}

	// 加载题库
	fmt.Print("正在加载题库...")
	content, err := smartReadFile("tiku.txt")
	if err != nil {
		fmt.Printf("\n❌ 错误: %v\n", err)
		fmt.Println("请确保 tiku.txt 文件存在于程序同一目录")
		fmt.Println("\n按任意键退出...")
		fmt.Scanln()
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			app.questions = append(app.questions, line)
		}
	}

	fmt.Printf(" ✅ 成功！共加载 %d 条题目\n", len(app.questions))

	// 显示前3题验证
	fmt.Println("\n题库预览（前3题）：")
	for i := 0; i < 3 && i < len(app.questions); i++ {
		fmt.Printf("  %d. %s\n", i+1, app.questions[i])
	}

	// 设置HTTP路由
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/status", statusHandler)
	http.HandleFunc("/api/search", searchHandler)
	http.HandleFunc("/api/copy", copyHandler)
	http.HandleFunc("/api/questions", questionsHandler)

	// 查找可用端口
	port := 8080
	for {
		addr := ":" + strconv.Itoa(port)
		if _, err := http.Get("http://localhost" + addr); err != nil {
			break
		}
		port++
	}

	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Printf("\n🌐 启动Web服务器: %s\n", url)
	fmt.Println("✅ 中文显示绝对正常！功能100%完整！")

	// 启动浏览器
	go func() {
		time.Sleep(1 * time.Second)
		openBrowser(url)
	}()

	// 启动HTTP服务器
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("🚀 服务器运行在 %s\n", addr)
	fmt.Println("📝 如果浏览器没有自动打开，请手动访问上面的网址")
	fmt.Println("⚠️  关闭此窗口将停止服务器")
	
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("❌ 服务器启动失败: %v\n", err)
		fmt.Println("\n按任意键退出...")
		fmt.Scanln()
	}
} 