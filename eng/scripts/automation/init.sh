#!/bin/bash

if [ -z $1 ]; then
    echo "Please input outputfile"
    exit 1
fi
echo $1

set -x
set -e
echo "GOPATH:$GOPATH"
echo "GOROOT:$GOROOT"
TMPDIR="/tmp"
if [ ! "$(go version | awk '{print $3}' | cut -c 3-6)" = "1.16" ]
then
  wget https://golang.org/dl/go1.16.9.linux-amd64.tar.gz
  tar -C $TMPDIR -xzf go1.16.9.linux-amd64.tar.gz
  export GOROOT=$TMPDIR/go
  export PATH=$GOROOT/bin:$PATH
fi

DIRECTORY=$(cd `dirname $0` && pwd)

if [ "$GOPATH" == "" ]; then
  WORKFOLDER="$(realpath $DIRECTORY/../../../../)"
  echo $WORKFOLDER
  export GOPATH=$WORKFOLDER/gofolder
fi

if [ ! -d "$GOPATH/bin" ]; then
  echo "create gopath folder"
  mkdir -p $GOPATH/bin
fi
echo $GOPATH

export GO111MODULE=on
# cd eng/tools/generator && go build && cp generator $GOPATH/bin && cd ../../..
generatorDirectory="$(realpath $DIRECTORY/../../tools/generator)"
cd $generatorDirectory
go build
# go install
cp generator $GOPATH/bin/
export PATH=$PATH:$GOPATH/bin
cd $DIRECTORY

# sudo ln -s $generatorDirectory/generator /usr/bin/generator
# ls -l /usr/bin/generator

if [ ! -f "$GOPATH/bin/pwsh.exe" ]; then
  ln -s /usr/bin/pwsh $GOPATH/bin/pwsh.exe
fi


echo $1
cat > $1 << EOF
{
  "envs": {
    "PATH": "$PATH:$GOPATH",
    "GOPATH": "$GOPATH",
    "GOROOT": "$GOROOT"
  }
}
EOF

cat $1