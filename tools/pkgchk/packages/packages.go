package packages

type Package interface {
	Root() string
	Path() string
	FullPath() string
	Name() string
	IsARMPackage() bool
}
