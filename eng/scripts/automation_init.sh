set -e
outputFile=""

if [ "$1" ]; then
  echo $1
  outputFile=$1
fi

if [ "$2" ]; then
  echo $2
  outputFile=$2
fi

if [ "$outputFile" != "" ]; then
  outputFile="$(realpath $outputFile)"
  echo "output json file: $outputFile"
fi


TMPDIR="/tmp"
if ! command -v semver &> /dev/null; then
    echo "Install semver"
    npm install -g semver
fi
current_version=$(go version | awk '{print $3}' | sed 's/go//')
if ! semver -r ">=1.23.0" "$current_version"; then
  wget -q https://golang.org/dl/go1.23.2.linux-amd64.tar.gz
  tar -C $TMPDIR -xzf go1.23.2.linux-amd64.tar.gz
  export GOROOT=$TMPDIR/go
  export PATH=$GOROOT/bin:$PATH
fi

DIRECTORY=$(cd `dirname $0` && pwd)

if [ "$GOPATH" = "" ]; then
  WORKFOLDER="$(realpath $DIRECTORY/../../../)"
  echo "WORKFOLDER: $WORKFOLDER"
  export GOPATH=$WORKFOLDER/gofolder
fi

if [ ! -d "$GOPATH/bin" ]; then
  echo "create gopath folder"
  mkdir -p $GOPATH/bin
fi
echo "GOPATH: $GOPATH"

export GO111MODULE=on

generatorDirectory="$(realpath $DIRECTORY/../tools/generator)"
cd $generatorDirectory
echo "Install generator..."
go build 2>&1

cp generator $GOPATH/bin/
rm generator
export PATH=$GOPATH/bin:$PATH
cd $DIRECTORY

if [ "$outputFile" != "" ]; then
  # Output environment information to the specified file
  cat > $outputFile << EOF
{
  "envs": {
    "PATH": "$GOPATH:$PATH",
    "GOPATH": "$GOPATH",
    "GOROOT": "$GOROOT"
  }
}
EOF
else
  # Output environment information to terminal
  echo '{
  "envs": {
    "PATH": "'"$GOPATH:$PATH"'",
    "GOPATH": "'"$GOPATH"'",
    "GOROOT": "'"$GOROOT"'"
  }
}'
fi

echo Install tsp-client
tspClientDir="$(realpath $DIRECTORY/../common/tsp-client)"
npm --prefix "$tspClientDir" ci 2>&1
