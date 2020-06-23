// +build modhack

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package to

// Necessary for safely adding multi-module repo.
// See: https://github.com/golang/go/wiki/Modules#is-it-possible-to-add-a-module-to-a-multi-module-repository
import _ "github.com/Azure/azure-sdk-for-go"
