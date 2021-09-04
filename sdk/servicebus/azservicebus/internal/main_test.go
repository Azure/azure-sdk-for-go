package internal

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
