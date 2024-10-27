package customError

import (
	"fmt"
	"runtime"
	"strings"
)

// CustomError wraps error with stack trace
type CustomError struct {
	Err   error
	Stack string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%v\nStack Trace:\n%s", e.Err, e.Stack)
}

// NewError creates a new error with stack trace
func NewError(err error) error {
	if err == nil {
		return nil
	}

	// Get the stack trace
	var sb strings.Builder
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		fmt.Fprintf(&sb, "\t%s:%d %s\n", file, line, fn.Name())
	}

	return &CustomError{
		Err:   err,
		Stack: sb.String(),
	}
}
