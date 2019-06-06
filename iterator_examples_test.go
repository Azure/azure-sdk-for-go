package servicebus_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-service-bus-go"
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
