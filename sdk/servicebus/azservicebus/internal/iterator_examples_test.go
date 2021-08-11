package internal

import (
	"context"
	"fmt"
)

func ExampleMessageIterator() {
	subject := AsMessageSliceIterator([]*Message{
		NewMessageFromString("hello"),
		NewMessageFromString("world"),
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
