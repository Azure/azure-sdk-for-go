package main

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/eng/tools/snippet-generator/cmd"
)

func main() {
	if err := cmd.Command().Execute(); err != nil {
		os.Exit(1)
	}
}
