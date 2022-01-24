package main

import "github.com/Azure/azure-sdk-for-go/sdk/internal/perf"

func main() {
	perf.Run([]perf.NewPerfTest{
		NewNoOpTest,
		NewSleepTest,
	})
}
