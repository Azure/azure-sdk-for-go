// Copyright 2021 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package packages

// Package defines a SDK package
type Package interface {
	// Root returns the root directory of the sdk
	Root() string
	// Path returns the relative path to the root directory
	Path() string
	// FullPath returns the full path of this package. It should satisfy FullPath() == filepath.Join(Root(), Path())
	FullPath() string
	// Name returns the name of this package
	Name() string
	// IsARMPackage returns true if this package is a management plane package, false otherwise.
	IsARMPackage() bool
}
