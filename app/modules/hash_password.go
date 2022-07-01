package modules

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashPasswrod(password string) string {
	if len(password) == 0 {
		return ""
	}
	h := sha256.New()
	h.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
