package utils

import (
	"encoding/base64"
	"crypto/rand"
	"crypto/sha256"
)

const pkceCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func GenerateCodeVerifier(length int) (string, error) {
	// Create code verifier
	buf := make([]byte, length)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	// Map each byte from the slice to a random char in the pkceCharset using a modulo bias
	for i, _ := range buf {
		buf[i] = pkceCharset[int(buf[i]) % len(pkceCharset)]
	}

	return string(buf), nil
}

func GenerateCodeChallenge(codeVerifier string) string {
	// Send code verifier through SHA256 hashing algorithm
	sum := sha256.Sum256([]byte(codeVerifier))

	// Base64 URL-safe encoding
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
