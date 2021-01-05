package main

import (
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd"
)

func main() {
	if err := cmd.Command().Execute(); err != nil {
		for _, line := range strings.Split(err.Error(), "\n") {
			log.Printf("[ERROR] %s", line)
		}
		os.Exit(1)
	}
}
