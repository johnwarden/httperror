package httperror_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"bytes"
	"fmt"
	"errors"
	"strings"
	"testing"

	"github.com/johnwarden/httperror"

	"github.com/stretchr/testify/assert"
)

func okHandler(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(200)
	_, _ = w.Write([]byte(`OK`))

	return nil
}

func testRequest(h http.Handler, path string) (int, string) {
	r, _ := http.NewRequest("GET", path, strings.NewReader(url.Values{}.Encode())) // URL-encoded payload

	{
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, r)
		resp := rr.Result()
		defer resp.Body.Close()
		// io.Copy(os.Stdout, res.Body)
		body, _ := ioutil.ReadAll(resp.Body)

		return resp.StatusCode, string(body)

	}
}

func TestRequest(t *testing.T) {
	{
		s, _ := testRequest(httperror.HandlerFunc(okHandler), "/")
		assert.Equal(t, 200, s, "got 200 OK response")
	}

	{
		s, _ := testRequest(notFoundHandler, "/foo")

		assert.Equal(t, 404, s, "got 404 Not Found response")
	}
}

func TestCustomErrorHandler(t *testing.T) {

	{
		s, m := testRequest(httperror.WrapHandlerFunc(helloHandler, customErrorHandler), "/")
		assert.Equal(t, 400, s, "got 400 Bad request response")
		assert.Equal(t, "400 Sorry, we couldn't parse your request: missing 'name' parameter\n", m, "got custom error message")
	}
}

// notFoundHandler is a HandlerFunc that does nothing but return NotFound. This
// should cause an appropriate Not Found error page to be served, since any
// error returned by an httperror.HandlerFuncs that is not handled by middleware
// will by handled by the default WriteErrorResponse method.
var notFoundHandler = httperror.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) error {

	w.Header().Set("Content-Type", "text/plain")
	return httperror.NotFound
})



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