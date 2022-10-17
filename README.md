Package httperror is for writing HTTP handlers that return errors instead of handling them directly.

The basic idea is described in Preslav Rachev's blog post [I Don't Like Go's Default HTTP Hanlers](https://preslav.me/2022/08/09/i-dont-like-golang-default-http-handlers/)

## Basic Example

	func helloHandler(w http.ResponseWriter, r *http.Request) error {

		w.Header().Set("Content-Type", "text/plain")

		name, ok := r.URL.Query()["name"];
		if !ok {
			return httperror.New("missing 'name' parameter", http.StatusBadRequest)
		}

		fmt.Fprintf(w, "Hello, %s\n", name[0])

		return nil;
	}

	func main() {

		h := httperror.HandlerFunc(helloHandler)

		http.Handle("/hello", h)

		http.ListenAndServe(":8080", nil)
	}


Unlike a standard HTTP handler function, helloHandler can return an
error. If the required URL parameter is missing, it uses `PublicErrorf` to
create and return an error with an embedded 400 Bad Request code and a public
error message.

Notice there is no explicit error handling code in this example. But if you
fetch http://localhost:8080/hello without a `name` URL parameter, a 400 Bad
Request page will be appropriately served.

This is because helloHandler is then converted into a `httperror.HandlerFunc`. This type has a
ServeHTTP method which makes it implements the standard library's
`http.Handler` interface. The ServeHTTP method handles any error returned by
the handler function using a default error handler, which extracts any HTTP status code and
public error message from the error value, checks the content type header,
and serves an appropriate error page with that content type, status code, and message.

## Advantages to the Error-Returning pattern

- more idiomatic Go (most go code uses the error-returning pattern)
- less boilerplate error handling code
- middleware can inspect errors, extract status codes, add context, and appropriate log and handle errors

For example, structured logging middleware can extract error status codes from
error values to include in logs, and then return the error to be
inspected or handled by other middleware.

## Extracting, Embedding, and Comparing HTTP Status Codes

	// Pre-Defined Errors
	e0 := httperror.NotFound

	// Extracting Status
	httperror.StatusCode(e1) // 404

	// Constructing Errors
	e1 := httperror.New("no such product ID", http.StatusNotFound)

	// Comparing Errors
	errors.Is(e1, httperror.NotFound) // true

	// Wrapping Errors
	var ErrNoSuchProductID = fmt.Errorf("no such product ID")
	e2 := httperror.Wrap(ErrNoSuchProductID, http.StatusNotFound)

	// Comparing Wrapped Errors
	errors.Is(e2, ErrNoSuchProductID) // true
	errors.Is(e2, httperror.NotFound) // also true!

## Public Error Messages

The default error handler, [httperror.WriteResponse] will only show a HTTP status code and the corresponding generic status text
(e.g. "404 Not Found"). It will not show detailed error messages to users, as this could like program implementation details to the public.

However, if the error value has a an embedded public error message, this will be displayed to the user. To embed a public error message,
create an error using [httperror.PublicErrorf] instead of [httperror.New]:

	e := httperror.PublicErrorf("Sorry, we can't find a product with this ID", 404)

Public error messages are extracted by [httperror.PublicMessage]:

	m := httperror.PublicMessage(e)

If your custom error type defines a `PublicMessage() string` method, then [httperror.PublicMessage] will call and
return the value from that method.

## Generic Handler and HandlerFunc Types

This package defines generic version of [httperror.Handler] and
[httperror.HandlerFunc]: [httperror.XHandler] and
[httperror.XHandlerFunc]. The latter allow your http handlers to accept a
third parameter of a generic type -- a common pattern for Go HTTP handlers
and middleware. See the httprouter example below.

## Use with Other Routers/Middleware/etc. Packages

Changing the signature of HTTP handler functions can effect almost all HTTP
handlers, routers, and middleware in your application. However these changes
are rather straightforward and should tend to simplify code.

This package is compatible with many other frameworks, routers, and middleware
in the Go ecosystem, because it is not a "framework": it just some some
types, default error handling code, and example patterns. Using any of these
types should not tightly couple your application code to this package. Even
the definitions of [httperror.Handler] and [httperror.HandlerFunc] are just a
few lines of code which can be copied into your codebase and customized.

Below we include an example of using this package with a
[github.com/julienschmidt/httprouter]
(https://github.com/julienschmidt/httprouter).



## In This Package

  - error values for all HTTP error status code (e.g. httperror.NotFound)
  - methods for embedding and extracting status codes in error values
  - default error handling functions
  - alternative Handler and HandlerFunc types and generic XHandler and XHandlerFunc types for writing handlers that return errors


## Similar Packages

- [github.com/caarlos0/httperr](https://github.com/caarlos0/httperr) uses a very similar, and I have modified the design of this package to minimize the differences.

Some differences include:
  - httperror's default error handler supports different content types 
  - httperror has functions for extracting error status codes and comparing error values
  - httperror.Handler and httperror.HandlerFunc implements http.Handler. There is no need for NewF. 



## Examples

The examples tests demonstrate some of the advantages of this approach.

