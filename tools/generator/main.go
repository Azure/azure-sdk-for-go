package main

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd"
)

func main() {
	if err := cmd.Command().Execute(); err != nil {
		log.Printf("[ERROR] %+v", err)
		os.Exit(1)
	}
}
