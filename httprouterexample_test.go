package httperror_test

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

// This handler func looks like the standard httprouter.Handle, but it can return an error.
func helloRouterHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	var name string
	if name = ps.ByName("name"); name == "" {
		return httperror.PublicErrorf(http.StatusBadRequest, "missing 'name' parameter")
	}

	fmt.Fprintf(w, "Hello, %s\n", name)

	return nil
}

// routerHandler converts a handler function of type httperror.XHandlerFunc
// [httprouter.Params] into a httprouter.Handle. Any errors returned by the
// wrapped handler function are handled by the default error handler
// httperror.DefaultErrorHandler.
func routerHandler(h httperror.XHandlerFunc[httprouter.Params]) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := h(w, r, ps)
		if err != nil {
			httperror.DefaultErrorHandler(w, err)
		}
	}
}
