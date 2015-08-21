package main

import (
	"fmt"
	"os"

	"github.com/azure/azure-sdk-for-go/arm/examples"
)

func help() {
	fmt.Println(`Select an example to run:
		Client  -- Basic storage account example
		Storage -- Create / Delete storage account example
		`)
}

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}
	args := os.Args[2:]

	e := os.Args[1]
	switch e {
	case "Client":
		examples.Client(args)
	case "Storage":
		examples.Storage(args)
	default:
		fmt.Printf("%s is not a know example\n", e)
		help()
	}
}
