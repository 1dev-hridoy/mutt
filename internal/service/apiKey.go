package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	"github.com/dishan1223/mutt/consts"
	"github.com/gofiber/fiber/v3/log"
)

func GenerateAPIKey() (string, string) {
	// All env variables that are returned by a function (e.g., GetAPIKeyBytes) are
	// returning a string. So, we need to convert them to int before using them.
	// strconv.Atoi is used to convert the string to int. If the conversion fails, it returns an error.
	API_KEY_BYTES, err := strconv.Atoi(consts.GetAPIKeyBytes())
	if err != nil {
		log.Error("Invalid API key size", "error", err)
		return "", ""
	}
	b := make([]byte, API_KEY_BYTES)
	_, err = rand.Read(b)
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
