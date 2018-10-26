package servicebus

import (
	"fmt"
	"github.com/Azure/azure-amqp-common-go/rpc"
	"reflect"
)

type (
	ErrMissingField string

	ErrIncorrectType struct {
		Key          string
		ExpectedType reflect.Type
		ActualValue  interface{}
	}

	ErrAMQP rpc.Response
)

func (e ErrMissingField) Error() string {
	return fmt.Sprintf("missing value %q", string(e))
}

// NewErrIncorrectType lets you skip using the `reflect` package. Just provide a variable of the desired type as
// 'expected'.
func newErrIncorrectType(key string, expected, actual interface{}) ErrIncorrectType {
	return ErrIncorrectType{
		Key:          key,
		ExpectedType: reflect.TypeOf(expected),
		ActualValue:  actual,
	}
}

func (e ErrIncorrectType) Error() string {
	return fmt.Sprintf(
		"value at %q was expected to be of type %q but was actually of type %q",
		e.Key,
		e.ExpectedType,
		reflect.TypeOf(e.ActualValue))
}

func (e ErrAMQP) Error() string {
	return fmt.Sprintf("server says (%d) %s", e.Code, e.Description)
}
