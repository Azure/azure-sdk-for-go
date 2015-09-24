package autorest

import (
	"fmt"
)

// Error describes the methods implemented by autorest errors.
type Error interface {
	error

	// PackageType should return the package type of the object emitting the error. For types, the
	// value should match that produced the the '%T' format specifier of the fmt package. For other
	// elements, such as functions, it returns just the package name (e.g., "autorest").
	PackageType() string

	// Method should return the name of the method raising the error.
	Method() string

	// Message should return the error message.
	Message() string

	// String should return a formatted containing all available details (i.e., PackageType, Method,
	// Message, and original error (if any)).
	String() string

	// Original should return the original error, if any, and nil otherwise.
	Original() error
}

type baseError struct {
	packageType string
	method      string
	message     string

	original error
}

// NewError creates a new Error conforming object from the passed packageType, method, and
// message. message is treated as a format string to which the optional args apply.
func NewError(packageType string, method string, message string, args ...interface{}) Error {
	return NewErrorWithError(nil, packageType, method, message, args...)
}

// NewErrorWithError creates a new Error conforming object from the passed packageType, method,
// message, and original error. message is treated as a format string to which the optional args
// apply.
func NewErrorWithError(original error, packageType string, method string, message string, args ...interface{}) Error {
	if _, ok := original.(Error); ok {
		return original.(Error)
	}
	return baseError{
		packageType: packageType,
		method:      method,
		message:     fmt.Sprintf(message, args...),
		original:    original,
	}
}

// PackageType returns the package type of the object emitting the error. For types, the value
// matches that produced the the '%T' format specifier of the fmt package. For other elements,
// such as functions, it returns just the package name (e.g., "autorest").
func (be baseError) PackageType() string {
	return be.packageType
}

// Method returns the name of the method raising the error.
func (be baseError) Method() string {
	return be.method
}

// Message is the error message.
func (be baseError) Message() string {
	return be.message
}

// Original returns the original error, if any, and nil otherwise.
func (be baseError) Original() error {
	return be.original
}

// Error returns the same formatted string as String.
func (be baseError) Error() string {
	return be.String()
}

// String returns a formatted containing all available details (i.e., PackageType, Method,
// Message, and original error (if any)).
func (be baseError) String() string {
	if be.original == nil {
		return fmt.Sprintf("%s:%s %s", be.packageType, be.method, be.message)
	}
	return fmt.Sprintf("%s:%s %s -- Original Error: %v", be.packageType, be.method, be.message, be.original)
}
