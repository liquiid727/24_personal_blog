// 通用工具：提供字符串 slug 生成
package utils

import (
    "regexp"
    "strings"
)

// Slugify 将任意字符串转换为 URL 友好的短标识（小写、连字符）
func Slugify(s string) string {
    s = strings.ToLower(s)
    s = strings.TrimSpace(s)
    s = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(s, "")
    s = regexp.MustCompile(`\s+`).ReplaceAllString(s, "-")
    s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
    return strings.Trim(s, "-")
}