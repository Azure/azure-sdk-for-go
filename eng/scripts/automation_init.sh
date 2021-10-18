set -x
set -e
if [ ! "$(go version | awk '{print $3}' | cut -c 3-6)" = "1.16" ]
then
  wget https://golang.org/dl/go1.16.9.linux-amd64.tar.gz
  tar -C $TMPDIR -xzf go1.16.9.linux-amd64.tar.gz
  export PATH=$TMPDIR/go/bin:$PATH
fi
export GOPATH="$(realpath ../../../..)"
if [ ! -d "$GOPATH/bin" ]
then
  mkdir $GOPATH/bin
fi
PATH=$PATH:$GOPATH/bin
export GO111MODULE=on
cd eng/tools/generator && go build && cp generator $GOPATH/bin && cd ../../..
ln -s /usr/bin/pwsh $GOPATH/bin/pwsh.exe
cat > $2 << EOF
{
  "envs": {
    "PATH": "$PATH:$GOPATH",
    "GOPATH": "$GOPATH"
  }
}
EOF
