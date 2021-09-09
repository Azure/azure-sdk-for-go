package azservicebus

import "fmt"

// implements `internal/errorinfo/NonRetriable`
type ErrClosed struct {
	link string
}

func (ec ErrClosed) NonRetriable() {}
func (ec ErrClosed) Error() string {
	return fmt.Sprintf("%s is closed and can no longer be used", ec.link)
}
