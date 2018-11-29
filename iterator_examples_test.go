package servicebus_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-service-bus-go"
	"github.com/opentracing/opentracing-go/log"
)

func ExampleMessageIterator() {
	subject := servicebus.AsMessageSliceIterator([]*servicebus.Message{
		servicebus.NewMessageFromString("hello"),
		servicebus.NewMessageFromString("world"),
	})

	for !subject.Done() {
		cursor, err := subject.Next(context.Background())
		if err != nil {
			log.Error(err)
		}
		fmt.Println(string(cursor.Data))
	}

	// Output:
	// hello
	// world
}
