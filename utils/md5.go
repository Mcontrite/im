package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(str string) string {
	md := md5.New()
	md.Write([]byte(str))
	return hex.EncodeToString(md.Sum(nil))
}
