package main

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"unicode/utf8"
)

func convertGBKToUTF8(input string) (string, error) {
	// 检测是否为 GBK 编码（Windows 中文文件名常用）
	if detectEncoding([]byte(input)) != "GBK" {
		return input, nil
	}
	// GBK → UTF-8 转换
	decoder := simplifiedchinese.GBK.NewDecoder()
	reader := transform.NewReader(bytes.NewReader([]byte(input)), decoder)
	result, err := io.ReadAll(reader)
	return string(result), err
}

func detectEncoding(data []byte) string {
	if utf8.Valid(data) {
		return "UTF-8"
	}
	decoder := simplifiedchinese.GBK.NewDecoder()
	if _, _, err := transform.Bytes(decoder, data); err == nil {
		return "GBK"
	}
	return "Unknown"
}
