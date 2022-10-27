package httperror

import (
	"context"
	"net/http"
)

type contextKey string

var errPtrKey = contextKey("errPtr")
var paramsKey = contextKey("params")

// StandardMiddleware is a standard http.Handler wrapper.
type StandardMiddleware = func(http.Handler) http.Handler

// XApplyStandardMiddleware applies middleware written for a standard [http.Handler] to an [httperror.XHandler].
// It works by passing parameters and returning errors through the context.
func XApplyStandardMiddleware[P any](h XHandler[P], m StandardMiddleware) XHandlerFunc[P] {

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		errPtr := ctx.Value(errPtrKey).(*error)
		p := ctx.Value(paramsKey).(P)

		err := h.Serve(w, r, p)

		*errPtr = err
	})

	handler = m(h)

	return func(w http.ResponseWriter, r *http.Request, p P) error {

		var err error
		c := r.Context()
		c = context.WithValue(c, errPtrKey, &err)
		c = context.WithValue(c, paramsKey, p)

		handler.ServeHTTP(w, r.WithContext(c))

		return err
	}
}

// ApplyStandardMiddleware applies middleware written for a standard [http.Handler] to an [httperror.XHandler].
// It works by passing parameters and returning errors through the context.
func ApplyStandardMiddleware(h Handler, m StandardMiddleware) HandlerFunc {
	errPtrKey := contextKey("errPtr")

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		errPtr := ctx.Value(errPtrKey).(*error)

		err := h.Serve(w, r)

		*errPtr = err
	})

	handler = m(h)

	return func(w http.ResponseWriter, r *http.Request) error {

		var err error
		c := r.Context()
		c = context.WithValue(c, errPtrKey, &err)

		handler.ServeHTTP(w, r.WithContext(c))

		return err
	}
}
