package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// HMACSHA256 computes a HMAC-SHA256 of data given the provided key.
func HMACSHA256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func HMACSHA256String(key []byte, data string) []byte {
	return HMACSHA256(key, []byte(data))
}

func SHA256Hex(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
