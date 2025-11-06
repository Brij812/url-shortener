package utils

import (
	"crypto/sha1"
	"encoding/base64"
)

func GenerateShortCode(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	hash := h.Sum(nil)
	code := base64.URLEncoding.EncodeToString(hash)
	return code[:6]
}
