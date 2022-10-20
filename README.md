Package httperror is for writing HTTP handlers that return errors instead of handling them directly. This package defines:

	- errors with embedded HTTP status codes
	- handler types that return errors
	- default sensible error and panic handling functions

This readme introduces this package an provides example usage. See the [godoc](https://pkg.go.dev/github.com/johnwarden/httperror) for more details.

## Basic Example

	func helloHandler(w http.ResponseWriter, r *http.Request) error {

		w.Header().Set("Content-Type", "text/plain")

		name, ok := r.URL.Query()["name"];
		if !ok {
			return httperror.New(http.StatusBadRequest, "missing 'name' parameter")
		}

		fmt.Fprintf(w, "Hello, %s\n", name[0])

		return nil;
	}

	func main() {

		h := httperror.HandlerFunc(helloHandler)

		http.Handle("/hello", h)

		http.ListenAndServe(":8080", nil)
	}


Unlike a standard [http.HandlerFunc](https://pkg.go.dev/net/http#HandlerFunc), the `helloHandler` function above can
return an error.  Although there is no explicit error handling code in this
example, if you run it and fetch http://localhost:8080/hello without a
`name` URL parameter, an appropriate plain-text 400 Bad Request page will be served.

This is because `helloHandler` is converted into a [httperror.HandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#HandlerFunc), which has a
`ServeHTTP` method and thus implements the standard library's [http.Handler](https://pkg.go.dev/net/http#Handler) interface. `ServeHTTP` handles the 400
Bad Raquest error returned by `helloHandler` using a default error
handler, which serves an appropriate error page given the content type and
the status code embedded in the error.

## Advantages to Returning Errors over Handling Them Directly

- more idiomatic Go
- reduce risk of "naked returns" as described by Preslav Rachev's in [I Don't Like Go's Default HTTP Hanlers](https://preslav.me/2022/08/09/i-dont-like-golang-default-http-handlers/)
- middleware can inspect errors, extract status codes, add context, and appropriately log and handle errors


## Custom Error Handlers

Use [WrapHandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#WrapHandlerFunc) to add a custom error handler. 


	func customErrorHandler(w http.ResponseWriter, e error) { 
		s := httperror.StatusCode(e)
		w.WriteHeader(s)
		// now serve an appropriate error response
	}

	h := httperror.WrapHandlerFunc(helloHandler, customErrorHandler)

Here is a [more complete example](#example-custom-error-handler).

## Middleware

Returning errors from functions enable some new middleware patterns. Here is an example of custom middleware that [logs errors](#example-log-middleware).

[PanicMiddleware](https://pkg.go.dev/github.com/johnwarden/httperror#PanicMiddleware) 
causes panics to be returned as errors. This way, users will be
served an appropriate 500 error response (instead of an empty response), and
middleware can inspect and log errors. A further refinement might be to
trigger a graceful server shutdown on panic.

## Extracting, Embedding, and Comparing HTTP Status Codes

	// Pre-Defined Errors
	e := httperror.NotFound

	// Extracting Status
	httperror.StatusCode(e) // 404

	// Constructing Errors
	e = httperror.New(http.StatusNotFound, "no such product ID")

	// Comparing Errors
	errors.Is(e, httperror.NotFound) // true

	// Wrapping Errors
	var ErrNoSuchProductID = fmt.Errorf("no such product ID")
	e = httperror.Wrap(ErrNoSuchProductID, http.StatusNotFound)

	// Comparing Wrapped Errors
	errors.Is(e, ErrNoSuchProductID) // true
	errors.Is(e, httperror.NotFound) // also true!

## Public Error Messages

The default error handler, [DefaultErrorHandler](https://pkg.go.dev/github.com/johnwarden/httperror#DefaultErrorHandler) will
not show the full error string to users, because these often contain stack traces or other implementation details that should not be exposed to the public.

But if the error value has an embedded public error message, the error handler will display this to the user. To embed a public error message,
create an error using [NewPublic](https://pkg.go.dev/github.com/johnwarden/httperror#NewPublic) or [PublicErrorf](https://pkg.go.dev/github.com/johnwarden/httperror#PublicErrorf) instead of [New](https://pkg.go.dev/github.com/johnwarden/httperror#New) or [Errorf](https://pkg.go.dev/github.com/johnwarden/httperror#Errorf):

	e := httperror.NewPublic(404, "Sorry, we can't find a product with this ID")

Public error messages are extracted by [PublicMessage](https://pkg.go.dev/github.com/johnwarden/httperror#PublicMessage):

	m := httperror.PublicMessage(e)

If your custom error type defines a `PublicMessage() string` method, then [PublicMessage](https://pkg.go.dev/github.com/johnwarden/httperror#PublicMessage) will call and return the value from that method.

## Generic Handler and HandlerFunc Types

This package defines generic versions of [httperror.Handler](https://pkg.go.dev/github.com/johnwarden/httperror#Handler) and
[httperror.HandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#HandlerFunc) -- [httperror.XHandler](https://pkg.go.dev/github.com/johnwarden/httperror#XHandler) and [httperror.XHandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#HandlerFunc) -- that accept a
third parameter of any type. 

The third parameter can contain parsed request parameters, authorized user
IDs, and other information required by  handlers. For example, the
`helloHandler` function in the introductory example might be cleaner if it
accepted its parameters as a struct.

	type HelloParams struct {
		Name string
	}

	func helloHandler(w http.ResponseWriter, r *http.Request, ps HelloParams) error { 
		fmt.Fprintf(w, "Hello, %s\n", ps.Name)
		return nil
	}

	h = httperror.XHandlerFunc[HelloParams](helloHandler)

Our handler function can use generic middleware written for [httperror.XHandler](https://pkg.go.dev/github.com/johnwarden/httperror#XHandler)s, such as such as (https://pkg.go.dev/github.com/johnwarden/httperror#XPanicMiddleware).

	h = httperror.XPanicMiddleware[HelloParams](h)

Ultimately the handler function needs to be provided by a HelloParams struct by middleware that parses request parameters. Here is an [example](#example-httprouter) that does this using the popular [github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) package.


## Use with Other Routers/Middleware/etc. Packages

Changing the signature of HTTP handler functions can effect almost all HTTP
handlers, routers, and middleware in your application. However these changes
are rather straightforward and should tend to simplify code.

This package is compatible with many other frameworks, routers, and middleware
in the Go ecosystem, because it is not a "framework": it just some some
types, default error handling code, and example patterns. Using any of these
types should not tightly couple your application code to this package. Even
the definitions of [httperror.Handler](https://pkg.go.dev/github.com/johnwarden/httperror#Handler) and
[httperror.HandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#HandlerFunc) are just a
few lines of code which can be copied into your codebase and customized.

See [this example](#example-httprouter) of using the error-returning handler
pattern with a [github.com/julienschmidt/httprouter]
(https://github.com/julienschmidt/httprouter).


## Similar Packages

[github.com/caarlos0/httperr](https://github.com/caarlos0/httperr) uses a very similar approach, for example the definition of: [httperr.HandlerFunc](https://pkg.go.dev/github.com/caarlos0/httperr#HandlerFunc) and [httperror.HandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#HandlerFunc) are identical. I have designed this package to be mostly compatible with httperr. 

## Examples

The complete examples below demonstrate some of the advantages of this approach.

## Example: Custom Error Handler

This example extends the basic example from the introduction by adding custom
error handler.


	package httperror_test

	import (
		"bytes"
		"errors"
		"fmt"
		"net/http"

		"github.com/johnwarden/httperror"
	)

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
			// Else use the default error handler.
			httperror.DefaultErrorHandler(w, err)
		}
	}


## Example: Log Middleware


The following example extends the basic example from the introduction by
adding custom logging middleware. Actual logging middleware would probably need
to be much more complex to correctly capture information from the response
such as the status code for successful requests.

	package httperror_test

	import (
		"fmt"
		"net/http"

		"github.com/johnwarden/httperror"
	)


	func Example_logMiddleware() {
		// This is the same helloHandler as the introduction
		h := httperror.HandlerFunc(helloHandler)

		// But add some custom middleware to handle and log errors.
		h = customLogMiddleware(h)

		_, o := testRequest(h, "/hello")
		fmt.Println(o)
		// Output: HTTP Handler returned error 400 Bad Request: missing 'name' parameter
		// 400 Sorry, we couldn't parse your request: missing 'name' parameter
	}

	func customLogMiddleware(h httperror.Handler) httperror.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {

			// TODO: custom pre-request actions such as wrapping the response writer.

			err := h.Serve(w, r)

			if err != nil {
				// TODO: insert your application's error logging code here.
				fmt.Printf("HTTP Handler returned error %s\n", err)
			}

			return err
		}
	}


## Example: HTTPRouter

This example illustrates the use of the error-returning paradigm described in this document with a popular router package, [github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter). To make things more interesting, the handler function accepts its parameters as a struct instead of a value of type [httprouter.Params](https://pkg.go.dev/github.com/julienschmidt/httprouter#Params), thereby decoupling the handler from the router. 

	import (
		"fmt"
		"net/http"


		"github.com/johnwarden/httperror"

		"github.com/julienschmidt/httprouter"
	)

	func Example_httprouter() {
		router := httprouter.New()

		h := routerHandler(helloRouterHandler)

		router.GET("/hello/:name", h)

		_, o := testRequest(router, "/hello/Sunshine")
		fmt.Println(o)
		// Output: Hello, Sunshine
	}

	type HelloParams struct {
		Name string
	}

	// This helloRouterHandler func looks like the standard http Handler, but it takes
	// a third argument of type HelloParams argument and can return an error.

	func helloRouterHandler(w http.ResponseWriter, r *http.Request, ps HelloParams) error { 
		if ps.Name == "" { 
			return httperror.NewPublic(http.StatusBadRequest, "missing 'name' parameter") 
		}

		fmt.Fprintf(w, "Hello, %s\n", ps.Name)

		return nil
	}

	// routerHandler wraps a handler function of type httperror.XHandler[HelloParams]
	// and converts it into a httprouter.Handle. The resulting function
	// converts its argument of type httprouter.Params into a value of type HelloParams,
	// and passes it to the inner handler function. 

	func routerHandler(h httperror.XHandler[HelloParams]) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

			var params HelloParams
			params.Name = ps.ByName("name")

			err := h(w, r, params)
			if err != nil {
				httperror.DefaultErrorHandler(w, err)
			}
		}
	}




