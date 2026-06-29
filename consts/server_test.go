package consts

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("PORT", "3000")
	os.Setenv("HASH_COST", "4")
	os.Setenv("API_KEY_BYTES", "32")
	os.Setenv("MAX_LOG_SIZE", "65536")
	os.Setenv("MAX_STACK_TRACE", "32768")
	os.Exit(m.Run())
}
