#!/bin/bash

set -x

REALEXITSTATUS=0
if [[ $MODE == 'default' ]]; then
    bash rungas.sh
    REALEXITSTATUS=$(($REALEXITSTATUS+$?))

    grep -L -r --include *.go --exclude-dir vendor -P "Copyright (\d{4}|\(c\)) Microsoft" ./ | tee /dev/stderr | test -z "$(< /dev/stdin)"
    REALEXITSTATUS=$(($REALEXITSTATUS+$?))

    test -z "$(go build $(go list ./... | grep -v vendor) | tee /dev/stderr)"
    REALEXITSTATUS=$(($REALEXITSTATUS+$?))
    
    test -z "$(go fmt $(go list ./... | grep -v vendor) | tee /dev/stderr)"
    REALEXITSTATUS=$(($REALEXITSTATUS+$?))

    go vet $(go list ./... | grep -v vendor)
    REALEXITSTATUS=$(($REALEXITSTATUS+$?))
    
    go test $(sh ./findTestedPackages.sh)
    REALEXITSTATUS=$(($REALEXITSTATUS+$?))
fi

if [[ $MODE == 'breaking' ]]; then
    go run ./tools/apidiff/main.go packages ./services FETCH_HEAD~1 FETCH_HEAD --copyrepo --breakingchanges
    REALEXITSTATUS=$(($REALEXITSTATUS+$?))
fi
exit $REALEXITSTATUS
