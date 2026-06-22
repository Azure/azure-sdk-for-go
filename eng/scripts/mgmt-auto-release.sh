#!/bin/bash
set -ex

if [ -z "${GH_TOKEN}" ]; then
  echo "GH_TOKEN environment variable is required"
  exit 1
fi

export PATH=$PATH:$HOME/go/bin
git config --global user.email "ReleaseHelper"
git config --global user.name "ReleaseHelper"

cd ../
git clone https://github.com/Azure/azure-sdk-for-go.git
git clone https://github.com/Azure/azure-rest-api-specs.git

cd azure-sdk-for-go
git remote add fork https://Azure:"${GH_TOKEN}"@github.com/Azure/azure-sdk-for-go.git
cd ../

go install github.com/Azure/azure-sdk-for-go/eng/tools/generator@latest

generator issue -t "${GH_TOKEN}" > sdk-release.json
cat sdk-release.json

file_size=`du -b ./sdk-release.json |awk '{print $1}'`
echo "sdk-release.json file size:" ${file_size}
if [ ${file_size} -le 70 ]; then
  echo "There are no services that need to be released"
else
  echo "run generator release-v2..."
  generator release-v2 ./azure-sdk-for-go ./azure-rest-api-specs ./sdk-release.json -t "${GH_TOKEN}"
fi