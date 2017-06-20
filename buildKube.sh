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
export TGOPATH=$GOPATH

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
    rm -rf $KPATH/src/$dep
    godep get $dep/...
done

rm -rf Godeps
rm -rf vendor
./hack/godep-save.sh

git status
# sorcery so ./hack/update-bazel.sh does not fail
rm -rf vendor/github.com/docker/docker/project/CONTRIBUTORS.md
cp $KPATH/src/github.com/docker/docker/CONTRIBUTING.md vendor/github.com/docker/docker/project
mv vendor/github.com/docker/docker/project/CONTRIBUTING.md vendor/github.com/docker/docker/project/CONTRIBUTORS.md

git status
./hack/update-bazel.sh
./hack/update-godep-licenses.sh
git status
git diff Godeps/Godeps.json

# try to build
EXITCODE=0
test -z "$(make 2> >(grep 'azure' | tee /dev/stderr))"
EXITCODE=$?

export GOPATH=$TGOPATH
exit $EXITCODE
