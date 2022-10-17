/*
Package httperror is for writing HTTP handlers that return errors instead of handling them directly. See the documentation at https://github.com/johnwarden/httperror
*/
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
// If the error doesn't have an embedded status code, it returns InternalServerError.
// If the error is nil, returns 200 OK.
func StatusCode(err error) int {
	var httpError httpError

	if err == nil {
		return http.StatusOK
	}

	if errors.As(err, &httpError) {
		return httpError.httpStatusCode()
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
