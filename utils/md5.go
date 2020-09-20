package utils

import (
	"crypto/md5"
	"encoding/hex"
)

//小写的
func MD5(str string) string {
	md := md5.New()
	md.Write([]byte(str)) // 需要加密的字符串为 123456
	return hex.EncodeToString(md.Sum(nil))
}

func MD5Password(password, salt string) string {
	return MD5(password + salt)
}

func ValidatePassword(input, salt, password string) bool {
	return MD5(input+salt) == password
}
