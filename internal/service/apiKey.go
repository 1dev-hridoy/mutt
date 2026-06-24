package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/dishan1223/mutt/consts"
	"github.com/gofiber/fiber/v3/log"
)

func GenerateAPIKey() (string, string) {
	b := make([]byte, consts.API_KEY_BYTES)
	_, err := rand.Read(b)
	if err != nil {
		log.Error("Failed to generate API key", "error", err)
		return "", ""
	}
	raw := hex.EncodeToString(b)
	return raw, HashAPIKey(raw)
}

func HashAPIKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}
