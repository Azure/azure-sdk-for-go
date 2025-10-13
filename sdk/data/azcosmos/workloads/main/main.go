package main

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/workloads"
)

func main() {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	workloads.Run() // or workloads.RunWorkload(ctx) if you make RunWorkload public
	log.Println("done")
}
