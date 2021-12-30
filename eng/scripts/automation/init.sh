#!/bin/bash
set -x
echo "GOPATH:$GOPATH"
echo "GOROOT:$GOROOT"
ls -l $GOROOT/bin

pwd
which go
TMPDIR="/tmp"
if [ ! "$(go version | awk '{print $3}' | cut -c 3-6)" = "1.16" ]
then
  wget https://golang.org/dl/go1.16.9.linux-amd64.tar.gz
  tar -C $TMPDIR -xzf go1.16.9.linux-amd64.tar.gz
  export GOROOT=$TMPDIR/go
  export PATH=$GOROOT/bin:$PATH
fi

# if [ "$GOPATH" == "" ]; then
#   DIRECTORY=$(cd `dirname $0` && pwd)
#   echo $DIRECTORY
#   export GOPATH=$DIRECTORY/../../../gofolder
# fi
DIRECTORY=$(cd `dirname $0` && pwd)
echo $DIRECTORY
export GOPATH=$DIRECTORY/../../../gofolder
# export GOPATH=$TMPDIR/gofolder
if [ ! -d "$GOPATH/bin" ]; then
  echo "create gopath folder"
  mkdir -p $GOPATH/bin
fi
echo $GOPATH

export GO111MODULE=on
# cd eng/tools/generator && go build && cp generator $GOPATH/bin && cd ../../..
cd $DIRECTORY/../tools/generator
go build
ls -l
# go install
cp generator $GOPATH/bin/
ls -l
export PATH=$PATH:$GOPATH/bin
# ln -s /usr/bin/pwsh $GOPATH/bin/pwsh.exe
pwd
ls -l $GOPATH
ls $GOPATH/bin/
ls $GOPATH/bin/generator
pwd
generator
sudo ln -s $GOPATH/bin/generator /usr/bin/generator
ln -s /usr/bin/pwsh $GOPATH/bin/pwsh.exe
ls -l /usr/bin/generator

echo $2
cat > $2 << EOF
{
  "envs": {
    "PATH": "$PATH:$GOPATH",
    "GOPATH": "$GOPATH"
  }
}
EOF
