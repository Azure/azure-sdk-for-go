package armnetwork_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// 1.start recording test framework
// 	gopath\github.com\Azure\azure-sdk-for-go .\eng\common\testproxy\docker-start-proxy.ps1 start D:\Projects\Go\src\github.com\Azure\azure-sdk-for-go

// 2.run test environment:
// 	AZURE_RECORD_MODE=record;PROXY_CERT=gopath\github.com\Azure\azure-sdk-for-go\eng\common\testproxy\dotnet-devcert.crt

func TestVirtualNetworksClient_BeginCreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"createVN","westus")
	rgName := *rg.Name
	defer clean()

	// create virtual network
	vnClient := armnetwork.NewVirtualNetworksClient(subscriptionID,cred,opt)
	vnName, err := createRandomName(t, "network")
	require.NoError(t, err)
	vnPoller, err := vnClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		armnetwork.VirtualNetwork{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus"),
			},
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.StringPtr("10.1.0.0/16"),
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	vnResp, err := vnPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *vnResp.Name, vnName)
}

/* error message
=== RUN   TestVirtualNetworksClient_BeginCreateOrUpdate
    virtualnetworks_client_test.go:49:
        	Error Trace:	virtualnetworks_client_test.go:49
        	Error:      	Received unexpected error:
        	Test:       	TestVirtualNetworksClient_BeginCreateOrUpdate
--- FAIL: TestVirtualNetworksClient_BeginCreateOrUpdate (27.33s)
*/


/* go env
set GO111MODULE=auto
set GOARCH=amd64
set GOBIN=
set GOCACHE=C:\Users\username\AppData\Local\go-build
set GOENV=C:\Users\username\AppData\Roaming\go\env
set GOEXE=.exe
set GOEXPERIMENT=
set GOFLAGS=
set GOHOSTARCH=amd64
set GOHOSTOS=windows
set GOINSECURE=
set GOMODCACHE=D:\Projects\Go\pkg\mod
set GONOPROXY=
set GONOSUMDB=
set GOOS=windows
set GOPATH=D:\Projects\Go
set GOPRIVATE=
set GOPROXY=https://goproxy.io,direct
set GOROOT=D:\Program Files\Go
set GOSUMDB=sum.golang.org
set GOTMPDIR=
set GOTOOLDIR=D:\Program Files\Go\pkg\tool\windows_amd64
set GOVCS=
set GOVERSION=go1.17
set GCCGO=gccgo
set AR=ar
set CC=gcc
set CXX=g++
set CGO_ENABLED=1
set GOMOD=
set CGO_CFLAGS=-g -O2
set CGO_CPPFLAGS=
set CGO_CXXFLAGS=-g -O2
set CGO_FFLAGS=-g -O2
set CGO_LDFLAGS=-g -O2
set PKG_CONFIG=pkg-config
set GOGCCFLAGS=-m64 -mthreads -fmessage-length=0 -fdebug-prefix-map=C:\Users\V-JIAH~1\AppData\Local\Temp\go-build2728918440=/tmp/go-build -gno-record-gcc-switches
*/