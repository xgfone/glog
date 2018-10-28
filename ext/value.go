package ext

import (
	"fmt"

	"github.com/go-stack/stack"
	"github.com/xgfone/miss"
)

// Caller is the same as miss.Caller(true), but removing the GOPATH prefix.
func Caller(format ...string) miss.Valuer {
	return func(depth int, level miss.Level) (interface{}, error) {
		return fmt.Sprintf("%+v", stack.Caller(depth+1)), nil
	}
}

// CallerStack returns a Valuer returning the caller stack without runtime.
//
// The default is using "%+s:%d:%n" as the format. See github.com/go-stack/stack
func CallerStack(format ...string) miss.Valuer {
	return func(depth int, level miss.Level) (interface{}, error) {
		s := stack.Trace().TrimBelow(stack.Caller(depth + 1)).TrimRuntime()
		if len(s) > 0 {
			return fmt.Sprintf("%+v", s), nil
		}
		return "", nil
	}
}
