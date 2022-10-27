package httperror

import (
	"fmt"
	"net/http"
)

// PanicMiddleware wraps a [httperror.Handler], returning a new [httperror.HandlerFunc] that
// recovers from panics and returns them as errors.
func PanicMiddleware(h Handler) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		defer func() {
			if r := recover(); r != nil {
				isErr := false
				if err, isErr = r.(error); !isErr {
					err = fmt.Errorf("%v", r)
				}
			}
		}()

		err = h.Serve(w, r)
		return
	}
}

// XPanicMiddleware wraps a [httperror.XHandler], returning a new [httperror.XHandlerFunc] that
// recovers from panics and returns them as errors.
func XPanicMiddleware[P any](h XHandler[P]) XHandlerFunc[P] {
	return func(w http.ResponseWriter, r *http.Request, p P) (err error) {
		defer func() {
			if r := recover(); r != nil {
				isErr := false
				if err, isErr = r.(error); !isErr {
					err = fmt.Errorf("%v", r)
				}
			}
		}()

		err = h.Serve(w, r, p)
		return
	}
}
