# This script tries to build Kubernetes, and find possible
# breaking changes regarding the Azure SDK for Go

set -x

# This should only run on cronjobs
if [ "pull_request" != $TRAVIS_EVENT_TYPE ]; then
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
./hack/godep-restore.sh
DEP='github.com/Azure/azure-sdk-for-go'
rm -rf $KPATH/src/$DEP
godep get $DEP/...
rm -rf Godeps
rm -rf vendor
./hack/godep-save.sh
./hack/update-bazel.sh
./hack/update-godeps-licenses.sh
./hack/update-staging-client-go.sh
git checkout -- $(git status -s | grep "^ D" | awk '{print $2}' | grep ^Godeps)

# try to build
EXITCODE=0
test -z "$(make 2> >(grep 'azure-sdk-for-go'))"
EXITCODE=$?

export GOPATH=$TGOPATH
exit EXITCODE
