package httperror

import (
	"fmt"
	"net/http"
)

// PanicMiddleware wraps a [httperror.Handler], returning a new [httperror.HandlerFunc] that
// recovers from panics and returns them as errors. The second argument is an optional 
// function that is called if there is a panic. This function can be used, for example, to
// cleanly shutdown the server. 
func PanicMiddleware(h Handler, s func()) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (err error) {

		defer func() {
			if r := recover(); r != nil {
				isErr := false
				if err, isErr = r.(error); !isErr {
					err = fmt.Errorf("%v", r)
				}
				if s != nil {
					// the shutdown function must be called in a goroutine. Otherwise, if it is used
					// to shutdown the server, we'll get a deadlock with the server shutdown function
					// waiting for this request handler to finish, and this request waiting for the
					// server shutdown function.
					go s()
				}
			}
		}()

		err = h.Serve(w, r)
		return
	}
}

// XPanicMiddleware wraps a [httperror.XHandler], returning a new [httperror.XHandlerFunc] that
// recovers from panics and returns them as errors. The second argument is an optional 
// function that is called if there is a panic. This function can be used, for example, to
// cleanly shutdown the server. 
func XPanicMiddleware[P any](h XHandler[P], s func()) XHandlerFunc[P] {
	return func(w http.ResponseWriter, r *http.Request, p P) (err error) {
		defer func() {
			if r := recover(); r != nil {
				isErr := false
				if err, isErr = r.(error); !isErr {
					err = fmt.Errorf("%v", r)
				}
				if s != nil {
					// the shutdown function must be called in a goroutine. See comment in PanicMiddleware above.
					go s()
				}
			}
		}()

		err = h.Serve(w, r, p)
		return
	}
}
