package httperror

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

// PublicErrorf returns a new public error with the given public error message
// and status code. The resulting value implements the
// the [httperror.Public] interface, meaning [httperror.PublicMessage] can be
// called on it to return the public error message for display to the user.
func PublicErrorf(status int, format string, args ...interface{}) error {
	return publicError{fmt.Sprintf(format, args...), statusError{status}}
}


type publicError struct {
	message string
	statusError
}

// Error returns the text corresponding to this HTTP error status code.
func (e publicError) Error() string {
	var b bytes.Buffer

	b.WriteString(strconv.Itoa(e.status))
	b.WriteString(" ")
	b.WriteString(http.StatusText(e.status))

	if e.message != "" {
		b.WriteString(": ")
		b.WriteString(e.message)
	}
	return b.String()
}

func (e publicError) PublicMessage() string {
	return e.message
}
