package httperror

import (
	"net/http"
)

// Handler is like the standard [http.Handler] interface type, but it also
// implements the Serve method which returns an error. When used as a standard
// [http.Handler], any errors will be handled by the default error handler [DefaultErrorHandler].
// But code that understands the httperror.Handler interface and can deal with
// returned errors can call the Serve method.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Serve(w http.ResponseWriter, r *http.Request) error
}

// XHandler is a generic version of [httperror.Handler]. The Serve method
// which accepts a third generic parameter.
type XHandler[P any] interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Serve(w http.ResponseWriter, r *http.Request, p P) error
}

// HandlerFunc is like the standard [http.HandlerFunc] type, but it returns an error.
// HandlerFunc implements both the [httperror.Handler] and the [http.Handler] interface.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// XHandlerFunc is a generic version of [httperror.HandlerFunc], that accepts a third generic
// parameter.
type XHandlerFunc[P any] func(w http.ResponseWriter, r *http.Request, p P) error

// ServeHTTP makes httperror.HandlerFunc implement the standard [http.Handler] interface.
// Any errors will be handled by the default error handler [DefaultErrorHandler].
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		DefaultErrorHandler(w, err)
	}
}

// ServeHTTP makes httperror.XHandlerFunc implement the standard [http.Handler] interface.
// Any errors will be handled by the default error handler [DefaultErrorHandler].
func (h XHandlerFunc[P]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var zeroValue P
	err := h(w, r, zeroValue)
	if err != nil {
		DefaultErrorHandler(w, err)
	}
}

// Serve makes [httperror.HandlerFunc] implement the [httperror.Handler] interface
func (h HandlerFunc) Serve(w http.ResponseWriter, r *http.Request) error {
	return h(w, r)
}

// Serve makes [httperror.XHandlerFunc] implement the [httperror.Handler] interface
func (h XHandlerFunc[P]) Serve(w http.ResponseWriter, r *http.Request, p P) error {
	return h(w, r, p)
}

// WrapHandlerFunc constructs an httperror.HandlerFunc with a custom error handler.
// Return an http.HandlerFunc.
func WrapHandlerFunc(h func(w http.ResponseWriter, r *http.Request) error, eh ErrorHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			eh(w, err)
		}
	})
}

// WrapXHandlerFunc constructs an httperror.XHandlerFunc with a custom error handler.
// Returns an function with the same signature but without the error return value.
func WrapXHandlerFunc[P any](h func(w http.ResponseWriter, r *http.Request, p P) error, eh ErrorHandler) func(w http.ResponseWriter, r *http.Request, p P) {
	return func(w http.ResponseWriter, r *http.Request, p P) {
		err := h(w, r, p)
		if err != nil {
			eh(w, err)
		}
	}
}
