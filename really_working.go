package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"unicode/utf8"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Windows剪贴板API
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

func main() {
	fmt.Println("================================================")
	fmt.Println("        超星学习通辅助工具 - 真正有效版")
	fmt.Println("================================================")
	fmt.Println()

	// 加载题库
	fmt.Print("正在加载题库文件 tiku.txt ... ")
	
	content, err := smartReadFile("tiku.txt")
	if err != nil {
		fmt.Printf("\n[错误] 无法读取题库文件: %v\n", err)
		fmt.Println("请确保 tiku.txt 文件在程序同一目录下")
		fmt.Printf("\n按回车键退出...")
		fmt.Scanln()
		return
	}

	// 解析题库
	var questions []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			questions = append(questions, line)
		}
	}

	fmt.Printf("成功! 共加载 %d 条题目\n", len(questions))
	fmt.Println()

	// 显示前3题验证加载正确
	fmt.Println("题库内容验证 (前3题):")
	for i := 0; i < 3 && i < len(questions); i++ {
		fmt.Printf("  %d. %s\n", i+1, questions[i])
	}
	fmt.Println()

	// 搜索功能
	searchFunc := func(keyword string) []string {
		keyword = strings.ToLower(strings.TrimSpace(keyword))
		if keyword == "" {
			return nil
		}

		var results []string
		for i, q := range questions {
			if strings.Contains(strings.ToLower(q), keyword) {
				results = append(results, fmt.Sprintf("[%d] %s", i+1, q))
				if len(results) >= 15 { // 限制15条结果
					results = append(results, "\n... (显示前15条结果)")
					break
				}
			}
		}
		return results
	}

	// 主循环
	fmt.Println("================================================")
	fmt.Println("                  使用说明")
	fmt.Println("================================================")
	fmt.Println("1. 输入关键词搜索题目")
	fmt.Println("2. 找到答案后会自动复制到剪贴板")
	fmt.Println("3. 在考试页面按 Ctrl+V 粘贴")
	fmt.Println("4. 输入 'exit' 或 'quit' 退出程序")
	fmt.Println("================================================")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print("请输入搜索关键词: ")
		input, _ := reader.ReadString('\n')
		keyword := strings.TrimSpace(input)

		if keyword == "" {
			continue
		}

		if keyword == "exit" || keyword == "quit" || keyword == "退出" {
			fmt.Println("感谢使用，再见!")
			break
		}

		// 特殊命令
		if keyword == "test" || keyword == "测试" {
			testText := "测试成功! 中文显示正常:\n计算机网络、数据库系统、操作系统、数据结构\n软件工程、算法设计、人工智能、机器学习"
			fmt.Println("\n=== 测试内容 ===")
			fmt.Println(testText)
			
			if err := setClipboardText(testText); err != nil {
				fmt.Printf("[错误] 复制失败: %v\n", err)
			} else {
				fmt.Println("\n[成功] 测试内容已复制到剪贴板!")
				fmt.Println("请在任意地方按 Ctrl+V 验证功能")
			}
			fmt.Println()
			continue
		}

		if keyword == "help" || keyword == "帮助" {
			fmt.Println("\n=== 帮助信息 ===")
			fmt.Println("- 输入关键词搜索题目")
			fmt.Println("- 输入 'test' 测试剪贴板功能")
			fmt.Println("- 输入 'list' 显示前10题")
			fmt.Println("- 输入 'exit' 退出程序")
			fmt.Println()
			continue
		}

		if keyword == "list" || keyword == "列表" {
			fmt.Println("\n=== 题库预览 (前10题) ===")
			count := 10
			if len(questions) < 10 {
				count = len(questions)
			}
			for i := 0; i < count; i++ {
				fmt.Printf("%d. %s\n", i+1, questions[i])
			}
			fmt.Printf("\n共有 %d 道题目\n\n", len(questions))
			continue
		}

		// 搜索题目
		fmt.Printf("\n正在搜索包含 '%s' 的题目...\n", keyword)
		
		results := searchFunc(keyword)
		
		if len(results) == 0 {
			fmt.Println("[结果] 未找到相关题目")
			fmt.Println("建议: 尝试使用更短的关键词或检查拼写")
		} else {
			fmt.Printf("[结果] 找到 %d 条相关题目:\n\n", len(results))
			
			resultText := strings.Join(results, "\n\n")
			fmt.Println(resultText)
			
			// 自动复制到剪贴板
			if err := setClipboardText(resultText); err != nil {
				fmt.Printf("\n[错误] 复制到剪贴板失败: %v\n", err)
			} else {
				fmt.Println("\n[成功] 搜索结果已自动复制到剪贴板!")
				fmt.Println("请在考试页面按 Ctrl+V 粘贴答案")
			}
		}
		
		fmt.Println("\n================================================")
	}
} 