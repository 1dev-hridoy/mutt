package consts

import "github.com/dishan1223/mutt/internal/config"

func GetPort() string {
	return ":" + config.MustGetEnv("PORT")
}

const HASH_COST = 10
