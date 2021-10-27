#!/bin/bash

if [ -z $1 ]; then
    echo "Please input resource Provider name"
    exit 1
fi
ResourceProvider=$1

if [ -z $2 ]; then
    echo "Please input sdk dir path"
    exit 1
fi
SDK_DIR_PATH=$2

CODE_PATH=${SDK_DIR_PATH}/sdk/${ResourceProvider}/arm${ResourceProvider}
cd ${CODE_PATH}

echo $GOPATH
echo $GOROOT

go get gotest.tools/gotestsum
ls -l $GOPATH/bin
# go get github.com/Azure/azure-sdk-for-go/sdk/azidentity@latest && go get github.com/Azure/azure-sdk-for-go/sdk/armcore@latest && go get github.com/Azure/azure-sdk-for-go/sdk/azcore@latest && go get github.com/Azure/azure-sdk-for-go/sdk/to@latest && go mod tidy
go mod tidy
# go get github.com/Azure/azure-sdk-for-go/sdk/$(ResourceProvider)/arm$(ResourceProvider)
$GOPATH/bin/gotestsum --format testname -- -coverprofile cover.out