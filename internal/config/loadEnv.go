package config

import "github.com/joho/godotenv"

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("PANIC :: Failed to load env.")
	}
}
