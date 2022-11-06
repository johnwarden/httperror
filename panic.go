package httperror

import (
	"errors"
	"fmt"
	"net/http"
)

var Panic = panicError{}

type panicError struct {
	innerError error
	message    string
}

func (e panicError) Error() string {
	if e.innerError != nil {
		return "panic: " + e.innerError.Error()
	}
	return "panic: " + e.message
}

func (e panicError) Unwrap() error {
	return e.innerError
}

func (e panicError) Is(other error) bool {
	if other == Panic {
		return true
	}
	return errors.Is(e.innerError, other)
}

// PanicMiddleware wraps a [httperror.Handler], returning a new [httperror.HandlerFunc] that
// recovers from panics and returns them as errors. Panic error can be identified using
// errors.Is(err, httperror.Panic)
func PanicMiddleware(h Handler) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		defer func() {
			if r := recover(); r != nil {
				isErr := false
				if err, isErr = r.(error); !isErr {
					err = panicError{nil, fmt.Sprintf("%v", r)}
				} else {
					err = panicError{err, ""}
				}
			}
		}()

		err = h.Serve(w, r)
		return
	}
}

// XPanicMiddleware wraps a [httperror.XHandler], returning a new [httperror.XHandlerFunc] that
// recovers from panics and returns them as errors. Panic error can be identified using
// errors.Is(err, httperror.Panic)
func XPanicMiddleware[P any](h XHandler[P]) XHandlerFunc[P] {
	return func(w http.ResponseWriter, r *http.Request, p P) (err error) {
		defer func() {
			if r := recover(); r != nil {
				isErr := false
				if err, isErr = r.(error); !isErr {
					err = panicError{nil, fmt.Sprintf("%v", r)}
				} else {
					err = panicError{err, ""}
				}
			}
		}()

		err = h.Serve(w, r, p)
		return
	}
}
