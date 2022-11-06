/*
Package httperror is for writing HTTP handlers that return errors instead of handling them directly. 

Please use the v1 branch of this package at [github.com/johnwarden/httperror]. The v2 branch has been discontinued.
*/
package httperror

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
)

// Design note: the httpError type is private. There is no value in making this public.
// Errors are passed around using the standard `error` interface type. The only time
// you would want to get at the underlying httpError value would be to extract the status code, but it
// is simpler to just call httperror.StatusCode(err), and in fact it is better
// as httperror.StatusCode returns the correct status code even for errors that aren't
// httpErrors (e.g. 500).

// httpError implements errors representing specific HTTP Status codes (from 400
// to 500). This type implements the standard error interface (with error strings
// obtained from http.StatusText), as well as the httperror.Error interface.
type httpError struct {
	status int
}

// statusCode returns the integer HTTP error status code.
func (e httpError) httpStatusCode() int {
	return e.status
}

// Error returns the text corresponding to this HTTP error status code.
func (e httpError) Error() string {
	var b bytes.Buffer

	b.WriteString(strconv.Itoa(e.status))
	b.WriteString(" ")
	b.WriteString(http.StatusText(e.status))
	return b.String()
}

// Is returns true if the target error is a status error with the same HTTP
// status code. It allows comparisons of the form
// errors.Is(err, http.StatusBadRequests)
func (e httpError) Is(target error) bool {
	if se, ok := target.(httpError); ok {
		if e.httpStatusCode() == se.status {
			return true
		}
	}
	return false
}

// Design note: keep this interface private. In my initial implementation this
// was public: I wanted to allow users to create custom types that return an
// HTTP status code. But this was a problem when it came to **comparing** error values.
// For example, if e is a custom type below, it is not possible to
// make this code to behave properly:
//
//	errors.Is(e, httperror.NotFound)
//
// It both solved this problem, and greatly simplified the implementation, to
// require all errors with HTTP status codes to be created by this package.
//

type httpStatusError = interface {
	httpStatusCode() int
}

// StatusCode extracts the HTTP status code from an error created by this package.
// If the error doesn't have an embedded status code, it returns InternalServerError.
// If the error is nil, returns 200 OK.
func StatusCode(err error) int {
	var httpError httpStatusError

	if err == nil {
		return http.StatusOK
	}

	if errors.As(err, &httpError) {
		return httpError.httpStatusCode()
	}

	return http.StatusInternalServerError
}

// BadRequest represents the StatusBadRequest HTTP error.
var BadRequest = httpError{http.StatusBadRequest}

// Unauthorized represents the StatusUnauthorized HTTP error.
var Unauthorized = httpError{http.StatusUnauthorized}

// PaymentRequired represents the StatusPaymentRequired HTTP error.
var PaymentRequired = httpError{http.StatusPaymentRequired}

// Forbidden represents the StatusForbidden HTTP error.
var Forbidden = httpError{http.StatusForbidden}

// NotFound represents the StatusNotFound HTTP error.
var NotFound = httpError{http.StatusNotFound}

// MethodNotAllowed represents the StatusMethodNotAllowed HTTP error.
var MethodNotAllowed = httpError{http.StatusMethodNotAllowed}

// NotAcceptable represents the StatusNotAcceptable HTTP error.
var NotAcceptable = httpError{http.StatusNotAcceptable}

// ProxyAuthRequired represents the StatusProxyAuthRequired HTTP error.
var ProxyAuthRequired = httpError{http.StatusProxyAuthRequired}

// RequestTimeout represents the StatusRequestTimeout HTTP error.
var RequestTimeout = httpError{http.StatusRequestTimeout}

// Conflict represents the StatusConflict HTTP error.
var Conflict = httpError{http.StatusConflict}

// Gone represents the StatusGone HTTP error.
var Gone = httpError{http.StatusGone}

// LengthRequired represents the StatusLengthRequired HTTP error.
var LengthRequired = httpError{http.StatusLengthRequired}

// PreconditionFailed represents the StatusPreconditionFailed HTTP error.
var PreconditionFailed = httpError{http.StatusPreconditionFailed}

// RequestEntityTooLarge represents the StatusRequestEntityTooLarge HTTP error.
var RequestEntityTooLarge = httpError{http.StatusRequestEntityTooLarge}

// RequestURITooLong represents the StatusRequestURITooLong HTTP error.
var RequestURITooLong = httpError{http.StatusRequestURITooLong}

// UnsupportedMediaType represents the StatusUnsupportedMediaType HTTP error.
var UnsupportedMediaType = httpError{http.StatusUnsupportedMediaType}

// RequestedRangeNotSatisfiable represents the StatusRequestedRangeNotSatisfiable HTTP error.
var RequestedRangeNotSatisfiable = httpError{http.StatusRequestedRangeNotSatisfiable}

// ExpectationFailed represents the StatusExpectationFailed HTTP error.
var ExpectationFailed = httpError{http.StatusExpectationFailed}

// Teapot represents the StatusTeapot HTTP error.
var Teapot = httpError{http.StatusTeapot}

// MisdirectedRequest represents the StatusMisdirectedRequest HTTP error.
var MisdirectedRequest = httpError{http.StatusMisdirectedRequest}

// UnprocessableEntity represents the StatusUnprocessableEntity HTTP error.
var UnprocessableEntity = httpError{http.StatusUnprocessableEntity}

// Locked represents the StatusLocked HTTP error.
var Locked = httpError{http.StatusLocked}

// FailedDependency represents the StatusFailedDependency HTTP error.
var FailedDependency = httpError{http.StatusFailedDependency}

// TooEarly represents the StatusTooEarly HTTP error.
var TooEarly = httpError{http.StatusTooEarly}

// UpgradeRequired represents the StatusUpgradeRequired HTTP error.
var UpgradeRequired = httpError{http.StatusUpgradeRequired}

// PreconditionRequired represents the StatusPreconditionRequired HTTP error.
var PreconditionRequired = httpError{http.StatusPreconditionRequired}

// TooManyRequests represents the StatusTooManyRequests HTTP error.
var TooManyRequests = httpError{http.StatusTooManyRequests}

// RequestHeaderFieldsTooLarge represents the StatusRequestHeaderFieldsTooLarge HTTP error.
var RequestHeaderFieldsTooLarge = httpError{http.StatusRequestHeaderFieldsTooLarge}

// UnavailableForLegalReasons represents the StatusUnavailableForLegalReasons HTTP error.
var UnavailableForLegalReasons = httpError{http.StatusUnavailableForLegalReasons}

// InternalServerError represents the StatusInternalServerError HTTP error.
var InternalServerError = httpError{http.StatusInternalServerError}

// NotImplemented represents the StatusNotImplemented HTTP error.
var NotImplemented = httpError{http.StatusNotImplemented}

// BadGateway represents the StatusBadGateway HTTP error.
var BadGateway = httpError{http.StatusBadGateway}

// ServiceUnavailable represents the StatusServiceUnavailable HTTP error.
var ServiceUnavailable = httpError{http.StatusServiceUnavailable}

// GatewayTimeout represents the StatusGatewayTimeout HTTP error.
var GatewayTimeout = httpError{http.StatusGatewayTimeout}

// HTTPVersionNotSupported represents the StatusHTTPVersionNotSupported HTTP error.
var HTTPVersionNotSupported = httpError{http.StatusHTTPVersionNotSupported}

// VariantAlsoNegotiates represents the StatusVariantAlsoNegotiates HTTP error.
var VariantAlsoNegotiates = httpError{http.StatusVariantAlsoNegotiates}

// InsufficientStorage represents the StatusInsufficientStorage HTTP error.
var InsufficientStorage = httpError{http.StatusInsufficientStorage}

// LoopDetected represents the StatusLoopDetected HTTP error.
var LoopDetected = httpError{http.StatusLoopDetected}

// NotExtended represents the StatusNotExtended HTTP error.
var NotExtended = httpError{http.StatusNotExtended}

// NetworkAuthenticationRequired represents the StatusNetworkAuthenticationRequired HTTP error.
var NetworkAuthenticationRequired = httpError{http.StatusNetworkAuthenticationRequired}
