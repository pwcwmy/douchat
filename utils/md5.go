package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 小写 Md5
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

// 大写 MD5
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 加密 salt 随机数
func MakePassWord(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 解密 salt 随机数
func ValidPassWord(plainpwd, salt string, password string) bool {
	return Md5Encode(plainpwd + salt) == password
}