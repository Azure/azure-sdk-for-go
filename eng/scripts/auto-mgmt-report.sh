#!/bin/bash
set -ex

export PATH=$PATH:$HOME/go/bin
export AZURE_DEVOPS_PERSONAL_ACCESS_TOKEN=$1
export AZURE_STORAGE_PRIMARY_ACCOUNT_KEY=$2

git clone https://github.com/Azure/azure-sdk-for-go.git

cd ./azure-sdk-for-go
sdkpath=`pwd`

cd ./eng/tools/mgmtreport
go run . -sdkpath $sdkpath -storageaccount $3