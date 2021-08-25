package modules

import (
	"crypto/md5"
	"encoding/base64"
)

func HashPasswrod(password string) string {
	if len(password) == 0 {
		return ""
	}
	hashMd5 := md5.New()
	return base64.StdEncoding.EncodeToString(hashMd5.Sum([]byte(password)))
}
