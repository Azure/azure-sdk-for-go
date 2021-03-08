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
