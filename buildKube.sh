# This script tries to build Kubernetes, and find possible
# breaking changes regarding the Azure SDK for Go

set -x

# This should only run on cronjobs
if [ "pull_request" != $TRAVIS_EVENT_TYPE ]; then
   exit 0
fi

# Only meant to run on latest go version
if [ "go version go1.8 linux/amd64" != "$(go version)" ]; then
    exit 0
fi

go get github.com/tools/godep

# This need to run in a different GOPATH
TGOPATH=$GOPATH

# get kubernetes
export KPATH=$HOME/code/kubernetes
export GOPATH=$KPATH
mkdir -p $KPATH/src
go get k8s.io/kubernetes
cd $KPATH/src/k8s.io/kubernetes

# update SDK (https://github.com/kubernetes/community/blob/master/contributors/devel/godep.md)
deps=(
    github.com/Azure/azure-sdk-for-go
    github.com/Azure/go-autorest
)

for dep in ${deps[*]}; do
    go get -u $dep
    godep update $dep/...
    ./hack/update-bazel.sh
done

./hack/update-godep-licenses.sh
git status

# try to build
EXITCODE=0
test -z "$(make 2> >(grep 'azure' | tee /dev/stderr))"
EXITCODE=$?

export GOPATH=$TGOPATH
exit $EXITCODE
