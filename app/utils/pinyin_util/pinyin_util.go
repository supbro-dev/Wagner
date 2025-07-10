/*
* @Author: supbro
* @Date:   2025/7/8 17:34
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/8 17:34
 */
package pinyin_util

import (
	"github.com/mozillazg/go-pinyin"
	"strings"
	"unicode"
)

// 初始化拼音转换器（全局变量提高性能）
var pinyinArgs = pinyin.NewArgs()
var pinyinFirstLetterArgs = pinyin.NewArgs()

func init() {
	// 配置拼音参数
	pinyinFirstLetterArgs.Style = pinyin.FirstLetter // 只取首字母
}

// 获取单个汉字的首字母
func getHanziInitial(hanzi rune) string {
	// 获取拼音
	pinyinSlice := pinyin.Pinyin(string(hanzi), pinyinFirstLetterArgs)

	// 确保有结果且非空
	if len(pinyinSlice) > 0 && len(pinyinSlice[0]) > 0 {
		return pinyinSlice[0][0]
	}

	return string(hanzi) // 无法转换时返回原字符
}

// 转换字符串：汉字转首字母，其他字符不变
func ConvertMixedString(input string) string {
	var result strings.Builder
	var currentHanzi strings.Builder // 缓存连续汉字

	// 遍历字符串中的每个字符
	for _, r := range input {
		if unicode.Is(unicode.Han, r) {
			// 如果是汉字，添加到缓存
			currentHanzi.WriteRune(r)
		} else {
			// 非汉字字符：先处理缓存的汉字
			if currentHanzi.Len() > 0 {
				result.WriteString(getHanziInitials(currentHanzi.String()))
				currentHanzi.Reset()
			}
			// 添加当前非汉字字符
			result.WriteRune(r)
		}
	}

	// 处理末尾可能剩余的汉字
	if currentHanzi.Len() > 0 {
		result.WriteString(getHanziInitials(currentHanzi.String()))
	}

	return result.String()
}

// 获取连续汉字的首字母（优化性能）
func getHanziInitials(hanziString string) string {
	if len(hanziString) == 0 {
		return ""
	}

	// 单字直接处理
	if len(hanziString) == 1 {
		return getHanziInitial([]rune(hanziString)[0])
	}

	// 批量处理多字
	var result strings.Builder
	for _, r := range hanziString {
		result.WriteString(getHanziInitial(r))
	}
	return result.String()
}
