#!/bin/bash

cd $CODEPATH
if [ -z $1 ]; then
    echo "Please input code path"
fi
CODE_PATH=$1
cd ${CODE_PATH}
# go get github.com/Azure/azure-sdk-for-go/sdk/azidentity@latest && go get github.com/Azure/azure-sdk-for-go/sdk/armcore@latest && go get github.com/Azure/azure-sdk-for-go/sdk/azcore@latest && go get github.com/Azure/azure-sdk-for-go/sdk/to@latest && go mod tidy
go mod tidy
# go mod download github.com/Azure/azure-sdk-for-go/sdk/armcore
# go get github.com/Azure/azure-sdk-for-go/sdk/armcore@v0.8.0
go build