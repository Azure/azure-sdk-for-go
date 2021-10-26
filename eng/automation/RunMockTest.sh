#!/bin/bash

cd $CODEPATH
if [ -z $1 ]; then
    echo "Please input code path"
fi
CODE_PATH=$1
cd ${CODE_PATH}

echo $GOPATH
echo $GOROOT

go get gotest.tools/gotestsum
ls -l $GOPATH/bin
# go get github.com/Azure/azure-sdk-for-go/sdk/azidentity@latest && go get github.com/Azure/azure-sdk-for-go/sdk/armcore@latest && go get github.com/Azure/azure-sdk-for-go/sdk/azcore@latest && go get github.com/Azure/azure-sdk-for-go/sdk/to@latest && go mod tidy
go mod tidy
# go get github.com/Azure/azure-sdk-for-go/sdk/$(ResourceProvider)/arm$(ResourceProvider)
$GOPATH/bin/gotestsum --format testname -- -coverprofile cover.out