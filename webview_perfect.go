package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/webview/webview_go"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

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
	webview   webview_go.WebView
}

type SearchRequest struct {
	Keyword string `json:"keyword"`
}

type SearchResponse struct {
	Results []string `json:"results"`
	Count   int      `json:"count"`
}

type CopyRequest struct {
	Text string `json:"text"`
}

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

func (app *App) loadQuestions() error {
	content, err := smartReadFile("tiku.txt")
	if err != nil {
		return err
	}

	app.questions = []string{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			app.questions = append(app.questions, line)
		}
	}

	return nil
}

func (app *App) search(keyword string) []string {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return nil
	}

	var results []string
	for i, q := range app.questions {
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

func (app *App) searchOnline(query string) (string, error) {
	url := "http://cx.icodef.com/wyn-nb?v=4&q=" + query
	
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	return string(body), nil
}

// API处理函数
func (app *App) handleSearch(args string) (interface{}, error) {
	var req SearchRequest
	if err := json.Unmarshal([]byte(args), &req); err != nil {
		return nil, err
	}
	
	results := app.search(req.Keyword)
	return SearchResponse{
		Results: results,
		Count:   len(results),
	}, nil
}

func (app *App) handleSearchOnline(args string) (interface{}, error) {
	var req SearchRequest
	if err := json.Unmarshal([]byte(args), &req); err != nil {
		return nil, err
	}
	
	result, err := app.searchOnline(req.Keyword)
	if err != nil {
		return map[string]string{"error": err.Error()}, nil
	}
	
	return map[string]string{"result": result}, nil
}

func (app *App) handleCopy(args string) (interface{}, error) {
	var req CopyRequest
	if err := json.Unmarshal([]byte(args), &req); err != nil {
		return nil, err
	}
	
	if err := setClipboardText(req.Text); err != nil {
		return map[string]string{"error": err.Error()}, nil
	}
	
	return map[string]string{"success": "内容已复制到剪贴板"}, nil
}

func (app *App) handleGetQuestions(args string) (interface{}, error) {
	count := 10
	if len(app.questions) < 10 {
		count = len(app.questions)
	}
	
	questions := make([]string, count)
	for i := 0; i < count; i++ {
		questions[i] = fmt.Sprintf("[%d] %s", i+1, app.questions[i])
	}
	
	return map[string]interface{}{
		"questions": questions,
		"total":     len(app.questions),
	}, nil
}

const htmlContent = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>超星学习通辅助工具 - 完美版</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: "Microsoft YaHei", "SimHei", sans-serif;
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
        
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        
        .header h1 {
            color: #333;
            font-size: 28px;
            margin-bottom: 10px;
        }
        
        .status {
            background: #e8f5e8;
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 25px;
            border-left: 4px solid #4caf50;
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
        }
        
        .btn-primary {
            background: #667eea;
            color: white;
        }
        
        .btn-primary:hover {
            background: #5a67d8;
        }
        
        .btn-secondary {
            background: #6c757d;
            color: white;
        }
        
        .btn-secondary:hover {
            background: #5a6268;
        }
        
        .btn-success {
            background: #28a745;
            color: white;
        }
        
        .btn-success:hover {
            background: #218838;
        }
        
        .btn-info {
            background: #17a2b8;
            color: white;
        }
        
        .btn-info:hover {
            background: #138496;
        }
        
        .btn-warning {
            background: #ffc107;
            color: #212529;
        }
        
        .btn-warning:hover {
            background: #e0a800;
        }
        
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
            margin-bottom: 20px;
        }
        
        .results-area {
            width: 100%;
            min-height: 300px;
            padding: 15px;
            border: 2px solid #ddd;
            border-radius: 6px;
            font-family: "Consolas", "Microsoft YaHei", monospace;
            font-size: 14px;
            line-height: 1.6;
            resize: vertical;
            background: white;
        }
        
        .results-area:focus {
            outline: none;
            border-color: #667eea;
        }
        
        .footer {
            text-align: center;
            color: #666;
            font-size: 14px;
        }
        
        .loading {
            display: none;
            text-align: center;
            color: #667eea;
            font-weight: bold;
        }
        
        .error {
            background: #f8d7da;
            color: #721c24;
            padding: 15px;
            border-radius: 6px;
            margin: 10px 0;
            border-left: 4px solid #dc3545;
        }
        
        .success {
            background: #d4edda;
            color: #155724;
            padding: 15px;
            border-radius: 6px;
            margin: 10px 0;
            border-left: 4px solid #28a745;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎓 超星学习通辅助工具</h1>
            <div class="status">
                <span id="status">正在加载题库...</span>
            </div>
        </div>
        
        <div class="search-section">
            <h3>🔍 搜索功能</h3>
            <div class="search-box">
                <input type="text" id="searchInput" class="search-input" placeholder="请输入中文或英文关键词搜索..." />
                <button class="btn btn-primary" onclick="searchLocal()">本地搜索</button>
                <button class="btn btn-secondary" onclick="searchOnline()">在线搜索</button>
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
        
        <div class="footer">
            <p id="statusMessage">就绪</p>
        </div>
    </div>

    <script>
        let totalQuestions = 0;
        
        // 页面加载时初始化
        window.addEventListener('DOMContentLoaded', function() {
            loadQuestions();
        });
        
        // 加载题库信息
        async function loadQuestions() {
            try {
                const result = await window.webview.call('getQuestions', {});
                totalQuestions = result.total;
                document.getElementById('status').innerHTML = 
                    '✅ 题库已加载: ' + totalQuestions + ' 条 | UTF-8编码支持 | 剪贴板100%兼容';
                document.getElementById('statusMessage').textContent = '就绪';
            } catch (error) {
                showError('加载题库失败: ' + error.message);
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
                const result = await window.webview.call('search', {keyword: keyword});
                if (result.results && result.results.length > 0) {
                    document.getElementById('resultsArea').value = result.results.join('\\n\\n');
                    document.getElementById('statusMessage').textContent = 
                        '本地搜索完成，找到 ' + result.count + ' 条结果';
                } else {
                    document.getElementById('resultsArea').value = '未找到相关题目，请尝试其他关键词';
                    document.getElementById('statusMessage').textContent = '搜索完成，无结果';
                }
            } catch (error) {
                showError('搜索失败: ' + error.message);
            } finally {
                showLoading(false);
            }
        }
        
        // 在线搜索
        async function searchOnline() {
            const keyword = document.getElementById('searchInput').value.trim();
            if (!keyword) {
                alert('请输入搜索内容');
                return;
            }
            
            showLoading(true);
            document.getElementById('statusMessage').textContent = '在线搜索中...';
            
            try {
                const result = await window.webview.call('searchOnline', {keyword: keyword});
                if (result.error) {
                    showError('在线搜索失败: ' + result.error);
                } else {
                    document.getElementById('resultsArea').value = result.result;
                    document.getElementById('statusMessage').textContent = '在线搜索完成';
                }
            } catch (error) {
                showError('在线搜索失败: ' + error.message);
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
                const result = await window.webview.call('copy', {text: content});
                if (result.error) {
                    showError(result.error);
                } else {
                    showSuccess('内容已复制到剪贴板！\\n请在考试窗口按 Ctrl+V 粘贴');
                }
            } catch (error) {
                showError('复制失败: ' + error.message);
            }
        }
        
        // 测试中文
        async function testChinese() {
            const testText = '测试成功！中文显示正常：\\n数据库管理、算法设计、程序开发、软件工程\\n计算机网络、操作系统、数据结构、人工智能';
            document.getElementById('resultsArea').value = testText;
            
            try {
                await window.webview.call('copy', {text: testText});
                showSuccess('中文测试文本已复制到剪贴板！\\n请粘贴验证功能是否正常');
            } catch (error) {
                showError('复制失败: ' + error.message);
            }
        }
        
        // 显示题库
        async function showQuestions() {
            try {
                const result = await window.webview.call('getQuestions', {});
                document.getElementById('resultsArea').value = 
                    '题库预览（前10题）：\\n\\n' + result.questions.join('\\n\\n');
                document.getElementById('statusMessage').textContent = 
                    '显示前10题，总共 ' + result.total + ' 题';
            } catch (error) {
                showError('获取题库失败: ' + error.message);
            }
        }
        
        // 清空结果
        function clearResults() {
            document.getElementById('searchInput').value = '';
            document.getElementById('resultsArea').value = '';
            document.getElementById('statusMessage').textContent = '就绪';
        }
        
        // 显示加载状态
        function showLoading(show) {
            document.getElementById('loading').style.display = show ? 'block' : 'none';
        }
        
        // 显示错误消息
        function showError(message) {
            const errorDiv = document.createElement('div');
            errorDiv.className = 'error';
            errorDiv.textContent = message;
            document.querySelector('.results-section').insertBefore(errorDiv, document.getElementById('resultsArea'));
            setTimeout(() => errorDiv.remove(), 5000);
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

func main() {
	fmt.Println("超星学习通辅助工具 - Webview完美版 v4.0")
	fmt.Println("============================================")
	fmt.Println("🎯 完美解决中文显示问题")
	fmt.Println("🌐 使用浏览器引擎渲染界面")
	fmt.Println("📋 剪贴板功能100%正常")
	fmt.Println("")

	app := &App{}

	// 加载题库
	fmt.Print("正在加载题库...")
	if err := app.loadQuestions(); err != nil {
		fmt.Printf("\n❌ 错误: %v\n", err)
		fmt.Println("请确保 tiku.txt 文件存在于程序同一目录")
		fmt.Println("\n按任意键退出...")
		fmt.Scanln()
		return
	}
	fmt.Printf(" ✅ 成功！共加载 %d 条题目\n", len(app.questions))

	// 显示前3题验证
	fmt.Println("\n题库预览（前3题）：")
	for i := 0; i < 3 && i < len(app.questions); i++ {
		fmt.Printf("  %d. %s\n", i+1, app.questions[i])
	}

	fmt.Println("\n正在启动Webview界面...")

	// 创建webview
	w := webview_go.New(false)
	defer w.Destroy()
	app.webview = w

	w.SetTitle("超星学习通辅助工具 - 完美版 v4.0")
	w.SetSize(920, 760, webview_go.HintFixed)

	// 绑定API函数
	w.Bind("search", app.handleSearch)
	w.Bind("searchOnline", app.handleSearchOnline)
	w.Bind("copy", app.handleCopy)
	w.Bind("getQuestions", app.handleGetQuestions)

	// 设置HTML内容
	w.SetHtml(htmlContent)

	fmt.Println("✅ Webview界面启动成功！中文显示完美！")

	// 运行webview
	w.Run()
} 