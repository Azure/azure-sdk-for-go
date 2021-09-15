package azservicebus

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Failed to load env file, no live tests will run: %s", err.Error())
	}

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
