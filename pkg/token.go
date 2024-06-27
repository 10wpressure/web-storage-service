package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

func GenerateToken(login string) string {
	hash := sha256.Sum256([]byte(login + time.Now().String()))
	return hex.EncodeToString(hash[:])
}

func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", http.ErrNoCookie
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
