set -x
set -e
export GOPATH="$(realpath ../../../..)"
if [ ! -d "$GOPATH/bin" ]
then
  mkdir $GOPATH/bin
fi
PATH=$PATH:$GOPATH/bin
export GO111MODULE=on
cd eng/tools/generator && go build && cp generator $GOPATH/bin && cd ../..
ln -s /usr/bin/pwsh $GOPATH/bin/pwsh.exe
cat > $2 << EOF
{
  "envs": {
    "PATH": "$PATH:$GOPATH",
    "GOPATH": "$GOPATH"
  }
}
EOF
