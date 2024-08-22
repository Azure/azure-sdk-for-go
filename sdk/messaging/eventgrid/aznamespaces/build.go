// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

//go:generate tsp-client sync
// NOTE: it'd be nice to avoid specifying my emitter options here. It doesn't seem to be using my tspconfig.yaml.
//go:generate tsp-client generate -c ./tspconfig.yaml --emitter-options "module-version=0.0.1;module=github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces;slice-elements-byval=true"
//go:generate goimports -w ../
//go:generate go run ./internal/generate
//go:generate goimports -w .

package aznamespaces
