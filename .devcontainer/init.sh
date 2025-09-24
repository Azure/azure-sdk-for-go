#!/bin/bash

set -ex

# a better version of gofmt - reorganizes/adds/removes imports automatically.
go install golang.org/x/tools/cmd/goimports@latest

# useful for previewing how docs will look on pkg.go.dev
go install golang.org/x/pkgsite/cmd/pkgsite@latest

# installs 'tsp' - the TypeSpec compiler
npm install -g @typespec/compiler

# installs 'tsp-client' - The TypeSpec+Azure client generator
npm install -g @azure-tools/typespec-client-generator-cli

# NOTE: you probably won't need this as most libraries should be using TypeSpec
npm install -g autorest
