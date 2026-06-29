package consts

import "github.com/dishan1223/mutt/internal/config"

func GetPort() string {
	return ":" + config.MustGetEnv("PORT")
}

func GetHashCost() string {
	return config.MustGetEnv("HASH_COST")
}
