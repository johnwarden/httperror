package httperror

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

// New constructs an error with an embedded an HTTP status code. The status
// code can be extracted using [httperror.StatusCode].
func New(status int, message string) error {
	return Wrap(fmt.Errorf(message), status)
}

// Errorf works like fmt.Errorf but it also embeds an HTTP status code. The
// status code can be extracted using [httperror.StatusCode].
func Errorf(status int, format string, args ...interface{}) error {
	return Wrap(fmt.Errorf(format, args...), status)
}

// Wrap wraps an error and embeds an HTTP status code that can be extracted
// using [httperror.StatusCode]
func Wrap(err error, status int) error {
	return wrappedError{err, statusError{status}}
}

type wrappedError struct {
	inner error
	statusError
}

// Error returns the HTTP status text corresponding to this error status code.
func (e wrappedError) Error() string {
	var b bytes.Buffer

	b.WriteString(strconv.Itoa(e.status))
	b.WriteString(" ")
	b.WriteString(http.StatusText(e.status))
	b.Write([]byte(": "))
	b.Write([]byte(e.inner.Error()))

	return b.String()
}

// Unwrap returns the inner error of a wrappedError
func (e wrappedError) Unwrap() error {
	return e.inner
}
