package httperror

import (
	"errors"
	"net/http"
)

// Design note: keep this interface private. In my initial implementation this
// was public: I wanted to allow users to create custom types that return an
// HTTP status code. But this was a problem when it came to **comparing** error values.
// For example, if e is a custom type below, it is not possible to
// make this code to behave properly:
//
//	errors.Is(e, httperror.NotFound)
//
// It both solved this problem, and greatly simplified the implementation, to
// require all errors with HTTP status codes to be created by this package.
type httpError = interface {
	httpStatusCode() int
}

// Public is an interface that requires a PublicMessage() string method.
// [httperror.PublicMessage] will extract the public error message from errors
// that implements this interface.
type Public = interface {
	PublicMessage() string
}

// StatusCode extracts the HTTP status code from an error created by this package.
// Also if the error implements a Temporary() bool function
// (see net.Error) and it returns true, then this function returns
// StatusServiceUnavailable. Otherwise it returns InternalServerError.
func StatusCode(err error) int {
	var httpError httpError

	if err == nil {
		return http.StatusOK
	}

	if errors.As(err, &httpError) {
		return httpError.httpStatusCode()
	}

	var temporaryErr interface{ Temporary() bool }
	if errors.As(err, &temporaryErr) {
		if temporaryErr.Temporary() {
			return http.StatusServiceUnavailable
		}
	}

	return http.StatusInternalServerError
}

// PublicMessage extracts the public message from errors that have a
// `PUblicMessage() string` method.
func PublicMessage(err error) string {
	var publicError Public

	if err == nil {
		return ""
	}

	if errors.As(err, &publicError) {
		return publicError.PublicMessage()
	}

	return ""
}
