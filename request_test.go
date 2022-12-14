package httperror_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/johnwarden/httperror"

	"github.com/stretchr/testify/assert"
)

func testRequest(h http.Handler, path string) (int, string) {
	r, _ := http.NewRequest("GET", path, strings.NewReader(url.Values{}.Encode())) // URL-encoded payload

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, r)
	resp := rr.Result()
	defer resp.Body.Close()
	// io.Copy(os.Stdout, res.Body)
	body, _ := io.ReadAll(resp.Body)

	return resp.StatusCode, string(body)
}

func TestRequest(t *testing.T) {
	{
		s, _ := testRequest(okHandler, "/")
		assert.Equal(t, 200, s, "got 200 OK response")
	}

	{
		s, _ := testRequest(notFoundHandler, "/foo")

		assert.Equal(t, 404, s, "got 404 Not Found response")
	}
}

func TestCustomErrorHandler(t *testing.T) {
	s, m := testRequest(httperror.WrapHandlerFunc(helloHandler, customErrorHandler), "/")
	assert.Equal(t, 400, s, "got 400 Bad request response")
	assert.Equal(t, "400 Sorry, we couldn't parse your request: missing 'name' parameter\n", m, "got custom error message")
}

func TestPanic(t *testing.T) {
	{
		h := getMeOuttaHere
		h = httperror.PanicMiddleware(h)

		var e error
		errorHandler := func(w http.ResponseWriter, err error) {
			e = err
			httperror.DefaultErrorHandler(w, err)
		}

		s, m := testRequest(httperror.WrapHandlerFunc(h, errorHandler), "/")
		assert.Equal(t, 500, s, "got 500 status code")
		assert.Equal(t, "500 Internal Server Error\n", m, "got 500 text/plain response")
		assert.True(t, errors.Is(e, httperror.Panic))
		assert.Equal(t, "panic: Get me outta here!", e.Error())
	}

	{
		h := fail
		h = httperror.PanicMiddleware(h)

		var e error
		errorHandler := func(w http.ResponseWriter, err error) {
			e = err
			httperror.DefaultErrorHandler(w, err)
		}

		s, m := testRequest(httperror.WrapHandlerFunc(h, errorHandler), "/")
		assert.Equal(t, 500, s, "got 500 status code")
		assert.Equal(t, "500 Internal Server Error\n", m, "got 500 text/plain response")
		assert.True(t, errors.Is(e, httperror.Panic))
		assert.True(t, errors.Is(e, sentinalError))
		assert.Equal(t, "panic: SOME_ERROR", e.Error())

	}
}

func TestApplyStandardMiddleware(t *testing.T) {
	{
		h := httperror.ApplyStandardMiddleware(okHandler, myMiddleware)
		s, _ := testRequest(h, "/")
		assert.Equal(t, 200, s)
	}

	{
		h := httperror.ApplyStandardMiddleware(notFoundHandler, myMiddleware)
		s, m := testRequest(h, "/")
		assert.Equal(t, 404, s)
		assert.Equal(t, "404 Not Found\n", m, "got correct response status")
	}

	{
		inner := httperror.XApplyStandardMiddleware[string](nameHandler, myMiddleware)

		h := httperror.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			return inner(w, r, "Bill")
		})

		s, m := testRequest(h, "/")
		assert.Equal(t, 200, s)
		assert.Equal(t, "Hello, Bill\n", m, "got middleware output")
	}
}

var sentinalError = fmt.Errorf("SOME_ERROR")

var getMeOuttaHere = httperror.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/plain")
	panic("Get me outta here!")
})

var fail = httperror.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/plain")
	panic(sentinalError)
})

var okHandler = httperror.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(200)
	_, _ = w.Write([]byte("OK\n"))

	return nil
})

// notFoundHandler is a HandlerFunc that does nothing but return NotFound. This
// should cause an appropriate Not Found error page to be served, since any
// error returned by an httperror.HandlerFuncs that is not handled by middleware
// will by handled by the default WriteErrorResponse method.
var notFoundHandler = httperror.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) error {
	w.Header().Set("Content-Type", "text/plain")
	return httperror.NotFound
})

var nameHandler = httperror.XHandlerFunc[string](func(w http.ResponseWriter, r *http.Request, name string) error {
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "Hello, %s\n", name)

	return nil
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

func myMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Foo", "Bar")
		h.ServeHTTP(w, r)
	})
}
