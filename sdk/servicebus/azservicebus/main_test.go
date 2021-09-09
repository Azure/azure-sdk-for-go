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
		log.Fatalf(".env file does not exist")
	}

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
