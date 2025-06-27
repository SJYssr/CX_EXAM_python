// +build ascii

package main

import (
	"io/ioutil"
	"unicode/utf8"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"strings"
)

// EncodingFixer handles various text encoding issues
type EncodingFixer struct{}

// NewEncodingFixer creates a new encoding fixer
func NewEncodingFixer() *EncodingFixer {
	return &EncodingFixer{}
}

// FixContent attempts to properly decode Chinese text
func (ef *EncodingFixer) FixContent(content string) string {
	// If content is already valid UTF-8, return as-is
	if utf8.ValidString(content) {
		return content
	}
	
	// Try to detect and convert from GB2312/GBK to UTF-8
	if gbkContent := ef.tryGBKToUTF8(content); gbkContent != "" {
		return gbkContent
	}
	
	// Try to fix common encoding issues
	fixed := ef.fixCommonIssues(content)
	return fixed
}

// tryGBKToUTF8 attempts to convert GBK/GB2312 to UTF-8
func (ef *EncodingFixer) tryGBKToUTF8(content string) string {
	decoder := simplifiedchinese.GBK.NewDecoder()
	result, _, err := transform.String(decoder, content)
	if err == nil && utf8.ValidString(result) {
		return result
	}
	return ""
}

// fixCommonIssues fixes common display issues
func (ef *EncodingFixer) fixCommonIssues(content string) string {
	// Remove or replace common problematic characters
	replacements := map[string]string{
		"◊": "",
		"â": "",
		"€": "",
		"™": "",
		"Â": "",
		"Ã": "",
	}
	
	result := content
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	
	return result
}

// LoadTikuWithEncoding loads tiku.txt with proper encoding handling
func (ef *EncodingFixer) LoadTikuWithEncoding() (string, error) {
	// Try reading as UTF-8 first
	content, err := ioutil.ReadFile("tiku.txt")
	if err != nil {
		return "", err
	}
	
	contentStr := string(content)
	
	// If it's valid UTF-8, return as-is
	if utf8.ValidString(contentStr) {
		return contentStr, nil
	}
	
	// Try to fix encoding issues
	fixed := ef.FixContent(contentStr)
	return fixed, nil
} 