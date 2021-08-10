package internal_test

import (
	"context"
	"fmt"

	servicebus "github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
)

func ExampleMessageIterator() {
	subject := servicebus.AsMessageSliceIterator([]*servicebus.Message{
		servicebus.NewMessageFromString("hello"),
		servicebus.NewMessageFromString("world"),
	})

	for !subject.Done() {
		cursor, err := subject.Next(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(cursor.Data))
	}

	// Output:
	// hello
	// world
}
