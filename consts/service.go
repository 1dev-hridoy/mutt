package consts

import "github.com/dishan1223/mutt/internal/config"

// All the env variables that needs to be stored in this file
// shall all be returned by a function.
// Make sure functions are named in a way that they are self explanatory.
func GetAPIKeyBytes() string {
	return config.MustGetEnv("API_KEY_BYTES")
}

func GetMaxLogSize() string {
	return config.MustGetEnv("MAX_LOG_SIZE")
}

func GetMaxStackTrace() string {
	return config.MustGetEnv("MAX_STACK_TRACE")
}
