if [ -z $1 ]; then
    echo "Please input outputfile"
    exit 1
fi
echo $1

outputFile=$1

if [ "$2" ]; then
  echo $2
  outputFile=$2
fi

set -e
outputFile="$(realpath $outputFile)"
echo "output json file: $outputFile"

TMPDIR="/tmp"
if [ ! "$(go version | awk '{print $3}' | cut -c 3-6)" = "1.22" ]
then
  wget -q https://golang.org/dl/go1.22.0.linux-amd64.tar.gz
  tar -C $TMPDIR -xzf go1.22.0.linux-amd64.tar.gz
  export GOROOT=$TMPDIR/go
  export PATH=$GOROOT/bin:$PATH
fi

DIRECTORY=$(cd `dirname $0` && pwd)

if [ "$GOPATH" = "" ]; then
  WORKFOLDER="$(realpath $DIRECTORY/../../../)"
  echo $WORKFOLDER
  export GOPATH=$WORKFOLDER/gofolder
fi

if [ ! -d "$GOPATH/bin" ]; then
  echo "create gopath folder"
  mkdir -p $GOPATH/bin
fi
echo $GOPATH

export GO111MODULE=on

generatorDirectory="$(realpath $DIRECTORY/../tools/generator)"
cd $generatorDirectory
go build 2>&1

cp generator $GOPATH/bin/
export PATH=$GOPATH/bin:$PATH
cd $DIRECTORY

cat > $outputFile << EOF
{
  "envs": {
    "PATH": "$GOPATH:$PATH",
    "GOPATH": "$GOPATH",
    "GOROOT": "$GOROOT"
  }
}
EOF

echo Install tsp-client
sudo npm install -g @azure-tools/typespec-client-generator-cli@latest 2>&1