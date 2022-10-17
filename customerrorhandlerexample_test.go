package httperror_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/johnwarden/httperror"
)

// The following example extends the basic example from the introduction by adding custom error handler that
// also logs errors.
func Example_customErrorHandler() {
	// This is the same helloHandler as the introduction. Add a custom error handler.
	h := httperror.WrapHandlerFunc(helloHandler, customErrorHandler)

	_, o := testRequest(h, "/hello")
	fmt.Println(o)
	// Output: 400 Sorry, we couldn't parse your request: missing 'name' parameter
}

func helloHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/plain")

	name, ok := r.URL.Query()["name"]
	if !ok {
		return httperror.NewPublic(http.StatusBadRequest, "missing 'name' parameter")
	}

	fmt.Fprintf(w, "Hello, %s\n", name[0])

	return nil
}

func customErrorHandler(w http.ResponseWriter, err error) {

	s := httperror.StatusCode(err)
	w.WriteHeader(s)

	if errors.Is(err, httperror.BadRequest) {
		// Handle 400 Bad Request errors by showing a user-friendly message.

		var m bytes.Buffer
		m.Write([]byte("Sorry, we couldn't parse your request: "))
		m.Write([]byte(httperror.PublicMessage(err)))

		httperror.WriteResponse(w, httperror.StatusCode(err), m.Bytes())

	} else {
		// Else use the default error handler, or customize it if you want something fancier.
		httperror.DefaultErrorHandler(w, err)
	}
}
