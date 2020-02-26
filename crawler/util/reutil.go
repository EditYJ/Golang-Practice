package util

import "regexp"

// 得到正则过滤的值
func FilterString(content []byte, re *regexp.Regexp) string {
	res := re.FindSubmatch(content)
	if len(res) >= 2 {
		return string(res[1])
	}
	return ""
}