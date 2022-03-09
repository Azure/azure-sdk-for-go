package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

func main() {
	perf.Run(map[string]perf.PerfMethods{
		"ClientGET": {Register: clientTestRegister, New: NewClientGETTest},
	})
}
