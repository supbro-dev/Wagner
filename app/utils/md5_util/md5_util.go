/*
* @Author: supbro
* @Date:   2025/6/12 15:43
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/12 15:43
 */
package md5_util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(content string) string {
	hashBytes := md5.Sum([]byte(content))
	hashString := hex.EncodeToString(hashBytes[:]) // 转换为十六进制字符串
	return hashString
}
