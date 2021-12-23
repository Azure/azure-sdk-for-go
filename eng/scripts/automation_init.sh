#!/bin/bash
set -x
echo $GOPATH
echo $GOROOT
PATH=$PATH:$GOPATH/bin
pwd
export GO111MODULE=on
cd eng/tools/generator && go build && cp generator $GOPATH/bin && cd ../../..
# ln -s /usr/bin/pwsh $GOPATH/bin/pwsh.exe
ls $GOPATH/bin/
ls $GOPATH/bin/generator.exe
pwd
ln -s $GOPATH/bin/generator.exe generator
ls generator