Package httperror is for writing HTTP handlers that return errors instead of handling them directly.

This readme introduces this package with examples. Individual types and methods
are documented in the [godoc](https://pkg.go.dev/github.com/johnwarden/httperror)

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


Unlike a standard HTTP handler function, the `helloHandler` function above can
return an error.  Although there is no explicit error handling code in this
example, if you run this example and fetch http://localhost:8080/hello without a
`name` URL parameter, a 400 Bad Request page will be appropriately served.

This is because helloHandler is  converted into a [httperror.HandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#HandlerFunc), which has a
ServeHTTP method and thus implements the standard library's [http.Handler](https://pkg.go.dev/net/http#Handler) interface. The ServeHTTP method handles
any error returned by the handler function using a default error handler,
which serves an appropriate error page given the content type, status code.

## Advantages to Returning Errors over Handling Them Directly

- more idiomatic Go (most go code uses the error-returning pattern)
- reduce risk of "naked returns" as described by Preslav Rachev's in [I Don't Like Go's Default HTTP Hanlers](https://preslav.me/2022/08/09/i-dont-like-golang-default-http-handlers/).
- middleware can inspect errors, extract status codes, add context, and appropriate log and handle errors


## Custom Error Handlers and Middleware

Use [WrapHandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#WrapHandlerFunc) to construct a [httperror.HandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#HandlerFunc) with a custom error handler. Here is an [example](https://pkg.go.dev/github.com/johnwarden/httperror#example-package-CustomErrorHandler).

And [here](https://pkg.go.dev/github.com/johnwarden/httperror#example-package-CustomMiddleware) is an example of custom middleware that wraps a handler to log errors and treats panics as errors.


## Extracting, Embedding, and Comparing HTTP Status Codes

	// Pre-Defined Errors
	e0 := httperror.NotFound

	// Extracting Status
	httperror.StatusCode(e1) // 404

	// Constructing Errors
	e1 := httperror.New(http.StatusNotFound, "no such product ID")

	// Comparing Errors
	errors.Is(e1, httperror.NotFound) // true

	// Wrapping Errors
	var ErrNoSuchProductID = fmt.Errorf("no such product ID")
	e2 := httperror.Wrap(ErrNoSuchProductID, http.StatusNotFound)

	// Comparing Wrapped Errors
	errors.Is(e2, ErrNoSuchProductID) // true
	errors.Is(e2, httperror.NotFound) // also true!

## Public Error Messages

The default error handler, [DefaultErrorHandler](https://pkg.go.dev/github.com/johnwarden/httperror#DefaultErrorHandler) will
not show detailed error messages to users, as this could leak program implementation details to the public.

However, if the error value has an embedded public error message, this will be displayed to the user. To embed a public error message,
create an error using [NewPublic](https://pkg.go.dev/github.com/johnwarden/httperror#NewPublic) instead of [New](https://pkg.go.dev/github.com/johnwarden/httperror#New):

	e := httperror.New(404, "Sorry, we can't find a product with this ID")

Public error messages are extracted by [PublicMessage](https://pkg.go.dev/github.com/johnwarden/httperror#PublicMessage):

	m := httperror.PublicMessage(e)

If your custom error type defines a `PublicMessage() string` method, then [PublicMessage](https://pkg.go.dev/github.com/johnwarden/httperror#PublicMessage) will call and return the value from that method.

## Generic Handler and HandlerFunc Types

This package defines generic version of [httperror.Handler](https://pkg.go.dev/github.com/johnwarden/httperror#Handler) and
[httperror.HandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#HandlerFunc): [httperror.XHandler](https://pkg.go.dev/github.com/johnwarden/httperror#XHandler) and
[httperror.XHandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#HandlerFunc). The latter allow your http handlers to accept a
third parameter of a generic type -- a common pattern for Go HTTP handlers
and middleware. See the [httprouter example](https://pkg.go.dev/github.com/johnwarden/httperror#example-package-Httprouter) in the godoc.

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

Below we include an [example](https://pkg.go.dev/github.com/johnwarden/httperror#example-package-Httprouter) of
using this package with a[github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter).



## Similar Packages

[github.com/caarlos0/httperr](https://github.com/caarlos0/httperr) uses a very similar approach, for example the definition of: [httperr.HandlerFunc](https://pkg.go.dev/github.com/caarlos0/httperr#HandlerFunc) and [httperror.HandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#HandlerFunc) are identical. And I have modified the design of this package to mimic httperr where possible, but there are a few way this package differs:

  - The [default error handler](https://pkg.go.dev/github.com/johnwarden/httperror#DefaultErrorHandler) supports more content types.
  - Functions for [extracting error status codes](https://pkg.go.dev/github.com/johnwarden/httperror#DefaultErrorHandler).
  - Comparison of errors to pre-defined error values (e.g. `errors.Is(err, httperror.NotFound)`).
  - [Public error messages](#public-error-messages).
  - [httperror.Handler](https://pkg.go.dev/github.com/johnwarden/httperror#Handler) and [httperror.HandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#HandlerFunc) implement both [http.Handler](https://pkg.go.dev/net/http#Handler) and [httperror.Handler](https://pkg.go.dev/github.com/johnwarden/httperror#Handler). When used as an [http.Handler](https://pkg.go.dev/net/http#Handler), errors returned by an [httperror.HandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#HandlerFunc) will be handled by the [default error handler](https://pkg.go.dev/github.com/johnwarden/httperror#DefaultErrorHandler).
  - Generic [httperror.XHandler](https://pkg.go.dev/github.com/johnwarden/httperror#XHandler) and [httperror.XHandlerFunc](https://pkg.go.dev/github.com/johnwarden/httperror#XHandlerFunc) types.


## Examples

The [examples](https://pkg.go.dev/github.com/johnwarden/httperror#readme-examples) in the godoc demonstrate some of the advantages of this approach.

