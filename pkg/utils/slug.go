// File: pkg/utils/slug.go
// Purpose: String slug generation utilities.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Normalizes strings into URL-friendly slugs.
package utils

import (
	"regexp"
	"strings"
)

// Slugify 生成 URL 友好标识
// Input: 任意字符串；会进行小写、修剪、过滤非字母数字、压缩空格与连字符
// Output: 仅包含小写字母数字与连字符的短标识
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, "-")
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
