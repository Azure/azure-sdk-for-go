// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package nonexport

type Public struct {
}

func (p *Public) PublicMethod() {

}

func (p *Public) privateMethod() {

}

func (p *Public) removePrivateMethod() {

}

type private struct {
}
