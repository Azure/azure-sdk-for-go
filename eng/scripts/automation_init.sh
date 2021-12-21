#!/bin/bash
echo $GOPATH
echo $GOROOT
PATH=$PATH:$GOPATH/bin
export GO111MODULE=on
cd eng/tools/generator && go build && cp generator $GOPATH/bin && cd ../../..
# ln -s /usr/bin/pwsh $GOPATH/bin/pwsh.exe