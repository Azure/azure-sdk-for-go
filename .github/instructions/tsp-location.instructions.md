---
applyTo: '**/tsp-location.yaml'
---

## Pre-requisites

1. Install the latest node LTS package (on Windows, you can use `winget install OpenJS.NodeJS.LTS`)
2. Install the TypeSpec compiler (`npm install -g @typespec/compiler`)
3. Install tsp-client (`npm install -g @azure-tools/typespec-client-generator-cli`)

## Regen

You can regen this client by running:

1. cd to the directory containing this file, for example `cd sdk/location/azlocation`
2. `go generate` in the root of the module. If there are any errors related to tsp-client, or tsp not being installed they probably follow the instructions in "Pre-requisites".
3. Create Go examples (in examples*_test.go files) for the new methods.
4. Update the CHANGELOG describing the features, bug fixes, or breaking changes in this client.

## Shortcuts

- You can browse to the URL/commit for this repo by going to https://github.com/<repo>/tree/<commit>/<directory>, using the fields from this file.
