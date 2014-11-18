package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func (c StorageClient) computeHmac256(message string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(c.accountKey)
	if err != nil {
		return "", fmt.Errorf("azure: base64 decode error: ", err)
	}

	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func currentTimeRfc1123Formatted() string {
	const dateLayout = http.TimeFormat // reuse from net/http package
	return timeRfc1123Formatted(time.Now().UTC())
}

func timeRfc1123Formatted(t time.Time) string {
	return t.Format(http.TimeFormat)
}
