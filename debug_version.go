package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 调试用的文件读取函数
func debugReadFile(filename string) (string, error) {
	fmt.Printf("=== 开始调试文件读取: %s ===\n", filename)
	
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	fmt.Printf("原始文件大小: %d 字节\n", len(data))
	
	// 显示前50字节的十六进制和字符
	fmt.Print("前50字节(HEX): ")
	for i := 0; i < 50 && i < len(data); i++ {
		fmt.Printf("%02X ", data[i])
	}
	fmt.Println()
	
	fmt.Print("前50字节(CHAR): ")
	for i := 0; i < 50 && i < len(data); i++ {
		if data[i] >= 32 && data[i] <= 126 {
			fmt.Printf("%c", data[i])
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println()

	// 检查BOM
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		fmt.Println("检测到UTF-8 BOM，已移除")
		data = data[3:]
	} else {
		fmt.Println("未检测到BOM标记")
	}

	// 检查UTF-8有效性
	if utf8.Valid(data) {
		fmt.Println("文件编码: UTF-8 ✓")
		result := string(data)
		fmt.Printf("UTF-8转换后字符数: %d\n", utf8.RuneCountInString(result))
		return result, nil
	}

	fmt.Println("文件编码: 非UTF-8，尝试转换...")

	// 尝试GBK
	gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	utf8Data, _, err := transform.Bytes(gbkDecoder, data)
	if err == nil {
		fmt.Println("成功从GBK转换为UTF-8")
		result := string(utf8Data)
		fmt.Printf("GBK转UTF-8后字符数: %d\n", utf8.RuneCountInString(result))
		return result, nil
	}

	// 尝试GB18030
	gb18030Decoder := simplifiedchinese.GB18030.NewDecoder()
	utf8Data, _, err = transform.Bytes(gb18030Decoder, data)
	if err == nil {
		fmt.Println("成功从GB18030转换为UTF-8")
		result := string(utf8Data)
		fmt.Printf("GB18030转UTF-8后字符数: %d\n", utf8.RuneCountInString(result))
		return result, nil
	}

	// 最后的清理
	fmt.Println("使用字符清理模式")
	cleaned := strings.ToValidUTF8(string(data), "?")
	fmt.Printf("清理后字符数: %d\n", utf8.RuneCountInString(cleaned))
	return cleaned, nil
}

func main() {
	fmt.Println("=== 超星学习通辅助工具 - 调试版本 ===")
	fmt.Println("正在进行详细的编码调试分析...")
	fmt.Println()

	// 测试字符串常量
	fmt.Println("1. 测试程序内字符串常量:")
	testStr := "测试中文字符串：数据库、算法、程序设计"
	fmt.Printf("   内容: %s\n", testStr)
	fmt.Printf("   UTF-8有效: %v\n", utf8.ValidString(testStr))
	fmt.Printf("   字节长度: %d, 字符数: %d\n", len(testStr), utf8.RuneCountInString(testStr))
	fmt.Println()

	// 测试文件读取
	fmt.Println("2. 测试文件读取:")
	content, err := debugReadFile("tiku.txt")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		fmt.Println("\n按回车键退出...")
		fmt.Scanln()
		return
	}

	fmt.Println("\n3. 解析题库内容:")
	lines := strings.Split(content, "\n")
	fmt.Printf("总行数: %d\n", len(lines))

	validLines := 0
	fmt.Println("\n前10行内容分析:")
	for i := 0; i < 10 && i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			validLines++
			fmt.Printf("第%d行: %s\n", validLines, line)
			fmt.Printf("  UTF-8有效: %v\n", utf8.ValidString(line))
			fmt.Printf("  字节长度: %d, 字符数: %d\n", len(line), utf8.RuneCountInString(line))
			
			// 检查是否包含特殊字符
			hasSpecial := false
			for _, r := range line {
				if r == 0xFFFD { // Unicode替换字符
					hasSpecial = true
					break
				}
			}
			if hasSpecial {
				fmt.Printf("  ⚠️  包含Unicode替换字符(乱码)\n")
			}
			fmt.Println()
		}
	}

	fmt.Printf("有效题目行数: %d\n", validLines)

	// 测试搜索功能
	fmt.Println("\n4. 测试搜索功能:")
	fmt.Print("请输入搜索关键词 (直接回车跳过): ")
	var keyword string
	fmt.Scanln(&keyword)
	
	if keyword != "" {
		fmt.Printf("搜索关键词: %s\n", keyword)
		keyword = strings.ToLower(strings.TrimSpace(keyword))
		
		var results []string
		scanner := bufio.NewScanner(strings.NewReader(content))
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := strings.TrimSpace(scanner.Text())
			if line != "" && strings.Contains(strings.ToLower(line), keyword) {
				results = append(results, fmt.Sprintf("[%d] %s", lineNum, line))
				if len(results) >= 5 { // 最多显示5个结果
					break
				}
			}
		}
		
		if len(results) == 0 {
			fmt.Println("未找到匹配结果")
		} else {
			fmt.Printf("找到 %d 个结果:\n", len(results))
			for _, result := range results {
				fmt.Printf("  %s\n", result)
			}
		}
	}

	fmt.Println("\n=== 调试完成 ===")
	fmt.Println("按回车键退出...")
	fmt.Scanln()
} 