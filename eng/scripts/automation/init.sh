#!/bin/bash
set -x
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

# if [ "$GOPATH" == "" ]; then
#   DIRECTORY=$(cd `dirname $0` && pwd)
#   echo $DIRECTORY
#   export GOPATH=$DIRECTORY/../../../gofolder
# fi
DIRECTORY=$(cd `dirname $0` && pwd)
# WORKFOLDER="$(realpath $DIRECTORY/../../../../)"
# echo $WORKFOLDER
# export GOPATH=$WORKFOLDER/gofolder
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
cd $DIRECTORY/../../tools/generator
go build
# go install
cp generator $GOPATH/bin/
export PATH=$PATH:$GOPATH/bin
# ln -s /usr/bin/pwsh $GOPATH/bin/pwsh.exe
ls $GOPATH/bin/
cd $DIRECTORY
generator
# sudo ln -s $GOPATH/bin/generator /usr/bin/generator
if [ ! -f "$GOPATH/bin/pwsh.exe" ]; then
  ln -s /usr/bin/pwsh $GOPATH/bin/pwsh.exe
fi
# ls -l /usr/bin/generator

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