package main

import "github.com/Azure/azure-sdk-for-go/sdk/internal/perf"

func main() {
	perf.RegisterArguments(RegisterArguments)
	perf.Run([]perf.NewPerfTest{
		NewDownloadTest,
		NewListTest,
		NewUploadTest,
	})
}
