# This script tries to build Terraform related packages,
# and find possible breaking changes regarding the Azure
# SDK for Go

set -x

# This should only run on cronjobs
if [ "pull_request" != $TRAVIS_EVENT_TYPE ]; then
    exit 0
fi

go get github.com/kardianos/govendor
REALEXITSTATUS=0

packages=(github.com/hashicorp/terraform
    github.com/terraform-providers/terraform-provider-azurerm
    github.com/terraform-providers/terraform-provider-azure)

for package in ${packages[*]}; do
    go get $package
    cd $GOPATH/src/$package

    # get list of dependencies on the SDK
    govendor list -no-status | grep 'azure-sdk-for-go' | tee deps.txt

    # update list of dependencies
    while read dep; do
        govendor remove $dep
    done <deps.txt

    while read dep; do
        govendor add $dep
    done <deps.txt

    # try to build
    make
    REALEXITSTATUS=$(($REALEXITSTATUS+$?))
done

exit $REALEXITSTATUS
