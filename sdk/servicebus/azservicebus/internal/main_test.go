package internal

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("../.env file does not exist")
	}

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
