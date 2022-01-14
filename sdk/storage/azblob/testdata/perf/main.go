package main

import "github.com/Azure/azure-sdk-for-go/sdk/internal/perf"

func main() {
	tests := []perf.PerfTest{
		&uploadPerfTest{},
		&downloadPerfTest{},
		&listBlobPerfTest{},
	}
	perf.Run(tests)
}
