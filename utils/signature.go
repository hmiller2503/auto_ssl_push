package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
	"strings"
)

func GenerateSignature(params map[string]string, privateKey string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signedStr strings.Builder
	for _, k := range keys {
		signedStr.WriteString(k)
		signedStr.WriteString(params[k])
	}
	signedStr.WriteString(privateKey)

	h := sha1.New()
	h.Write([]byte(signedStr.String()))
	return hex.EncodeToString(h.Sum(nil))
}
