package httperror_test

import (
	"fmt"
	"net/http"

	"github.com/johnwarden/httperror"
)

// The following example extends the basic example from the introduction by
// adding custom logging middleware. Note actual logging middleware would
// need to be much more complex to correctly capture information from the
// response such as the status code for successful requests.

func Example_customMiddleware() {
	// This is the same helloHandler as the introduction
	h := httperror.HandlerFunc(helloHandler)

	// But add some custom middleware to handle and log errors.
	h = customMiddleware(h)

	_, o := testRequest(h, "/hello")
	fmt.Println(o)
	// Output: HTTP Handler returned error 400 Bad Request: missing 'name' parameter
	// 400 Sorry, we couldn't parse your request: missing 'name' parameter
}

func customMiddleware(h httperror.Handler) httperror.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {

		// Handle panics so that users see an actual error page on panic,
		// instead of an empty page.
		defer HandlePanics(w, customErrorHandler)

		// TODO: custom pre-request actions such as wrapping the response writer.

		err := h.Serve(w, r)

		if err != nil {
			// TODO: insert your application's error logging code here.
			fmt.Printf("HTTP Handler returned error %s\n", err)

			customErrorHandler(w, err)
		} else {
			fmt.Printf("HTTP Handler returned ok\n")
		}

		return nil
	}
}

// HandlePanics recovers from panics, converts the panic to an error with a
// stack and calls the given error handler. It should be called in a defer
func HandlePanics(w http.ResponseWriter, eh func(w http.ResponseWriter, e error)) {
	if rec := recover(); rec != nil {

		// Convert the panic value into an error
		err, isErr := rec.(error)
		if !isErr {
			err = fmt.Errorf("%v", rec)
		}

		// Handle eeror.
		eh(w, err)
	}
}
