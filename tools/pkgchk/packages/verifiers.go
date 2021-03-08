package packages

// VerifyFunc ...
type VerifyFunc func(p Package) error

// Verifier could verify a SDK package
type Verifier interface {
	// Verify verifies the given package
	Verify(pkg Package) []error
}
