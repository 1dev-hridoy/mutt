package config

import (
	"os"

	"github.com/joho/godotenv"
)

func MustLoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("PANIC :: Failed to load env.")
	}
}

func MustGetEnv(e string) string {
	d := os.Getenv(e)
	if d == "" {
		panic("PANIC :: " + e + " is not set.")
	}
	return d
}
