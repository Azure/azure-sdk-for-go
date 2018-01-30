package servicebus

import (
	"fmt"
	"time"
)

// ptrBool takes a boolean and returns a pointer to that bool. For use in literal pointers, ptrBool(true) -> *bool
func ptrBool(toPtr bool) *bool {
	return &toPtr
}

// ptrString takes a string and returns a pointer to that string. For use in literal pointers,
// ptrString(fmt.Sprintf("..", foo)) -> *string
func ptrString(toPtr string) *string {
	return &toPtr
}

// durationTo8601Seconds takes a duration and returns a string period of whole seconds (int cast of float)
func durationTo8601Seconds(duration *time.Duration) *string {
	return ptrString(fmt.Sprintf("PT%dS", int(duration.Seconds())))
}
