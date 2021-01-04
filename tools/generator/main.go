package main

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd"
)

func main() {
	if err := cmd.Command().Execute(); err != nil {
		log.Printf("[ERROR] generation meets an error: \n%+v", err)
		os.Exit(1)
	}
}
