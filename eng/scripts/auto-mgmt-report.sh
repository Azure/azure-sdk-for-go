#!/bin/bash
set -ex

export PATH=$PATH:$HOME/go/bin
export AZURE_DEVOPS_PERSONAL_ACCESS_TOKEN=$1
export AZURE_STORAGE_PRIMARY_ACCOUNT_KEY=$2

# git config --global user.email "ReleaseHelper"
# git config --global user.name "ReleaseHelper"

cd ../
git clone https://github.com/Azure/azure-sdk-for-go.git

go install github.com/Azure/azure-sdk-for-go/eng/tools/mgmtreport@latest

mgmtreport -sdkpath ./azure-sdk-for-go -storageaccount chenjieshistorage