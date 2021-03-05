package packages

type VerifyFunc func(p Package) error

type Verifier interface {
	Verify(pkg Package) []error
}
