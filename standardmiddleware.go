package httperror

import (
	"context"
	"net/http"
)

type contextKey string

var key = contextKey("key")

// StandardMiddleware is a standard http.Handler wrapper.
type StandardMiddleware = func(http.Handler) http.Handler

type standardMiddleware[P any] struct {
	params P
	err    error
}

// XApplyStandardMiddleware applies middleware written for a standard
// [http.Handler] to an [httperror.XHandler], returning an
// [httperror.XHandler]. It is possible to apply standard middleware to
// [httperror.XHandler] without using this function,because
// [httperror.XHandler] implements the standard [http.Handler] interface.
// However, the result would be an [http.Handler], not an
// [httperror.XHandler], and so parameters could not passed to it and it
// could not return an error. This function solves that problem by passing
// errors and parameters through the context.
func XApplyStandardMiddleware[P any](h XHandler[P], ms ...StandardMiddleware) XHandlerFunc[P] {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sm := ctx.Value(key).(*standardMiddleware[P])

		sm.err = h.Serve(w, r, sm.params)
	})

	for _, m := range ms {
		handler = m(handler)
	}

	return func(w http.ResponseWriter, r *http.Request, p P) error {
		sm := &standardMiddleware[P]{p, nil}
		c := r.Context()
		c = context.WithValue(c, key, sm)

		handler.ServeHTTP(w, r.WithContext(c))

		return sm.err
	}
}

// ApplyStandardMiddleware applies middleware written for a standard
// [http.Handler] to an [httperror.Handler], returning an
// [httperror.Handler]. It is possible to apply standard middleware to
// [httperror.Handler] without using this function,because
// [httperror.Handler] implements the standard [http.Handler] interface.
// However, the result would be an [http.Handler], not an
// [httperror.Handler], and so parameters could not passed to it and it
// could not return an error. This function solves that problem by passing
// errors and parameters through the context.
func ApplyStandardMiddleware(h Handler, ms ...StandardMiddleware) HandlerFunc {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sm := ctx.Value(key).(*standardMiddleware[any])

		sm.err = h.Serve(w, r)
	})

	for _, m := range ms {
		handler = m(handler)
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		sm := &standardMiddleware[any]{}
		c := r.Context()
		c = context.WithValue(c, key, sm)

		handler.ServeHTTP(w, r.WithContext(c))

		return sm.err
	}
}
