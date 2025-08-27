package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

func GenerateRandomFileName(extension string) string {
	// 获取当前时间
	now := GetTimeNow()
	// 格式化时间为 YYYYMMDDHHMMSS 格式
	timestamp := now.Format("20060102150405")

	// 生成随机的 UUID
	uuid := generateUUID()
	// 组合文件名、时间戳和扩展名
	return fmt.Sprintf("%s_%s.%s", timestamp, uuid, extension)
}

// generateUUID 生成随机的 UUID
func generateUUID() string {
	// 生成一个 16 字节的随机数
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	// 将随机数转换为 32 位的十六进制字符串
	return hex.EncodeToString(b)
}

// GetContentType 根据文件扩展名获取内容类型
func GetContentType(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	default:
		return "application/octet-stream"
	}
}

// 判断是否为图片文件
func IsImageFile(ext string) bool {
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
		".svg":  true,
	}
	return imageExtensions[ext]
}

// 判断是否为可直接展示的文件类型
func IsViewableFile(ext string) bool {
	// 如果是图片文件，可以直接展示
	if IsImageFile(ext) {
		return true
	}

	// 其他可直接展示的文件类型
	viewableExtensions := map[string]bool{
		// HTML 文件
		".html": true,
		".htm":  true,
		// 文本文件
		".txt":  true,
		".text": true,
		".log":  true,
		// 代码文件
		".js":   true,
		".css":  true,
		".json": true,
		".xml":  true,
		".yaml": true,
		".yml":  true,
		// Markdown
		".md":       true,
		".markdown": true,
	}
	return viewableExtensions[ext]
}
