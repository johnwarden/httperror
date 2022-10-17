package httperror

import (
	"bytes"
	"net/http"
	"strconv"
)

// statusError implements errors representing specific HTTP Status codes (from 400
// to 500). This type implements the standard error interface (with error strings
// obtained from http.StatusText), as well as the httperror.Error interface.
type statusError struct {
	status int
}

// httpStatusCode returns the integer HTTP error status code.
func (e statusError) httpStatusCode() int {
	return e.status
}

// Error returns the text corresponding to this HTTP error status code.
func (e statusError) Error() string {
	var b bytes.Buffer

	b.WriteString(strconv.Itoa(e.status))
	b.WriteString(" ")
	b.WriteString(http.StatusText(e.status))
	return b.String()
}

// Is returns true if the target error is a status error with the same HTTP
// status code. It allows comparisons of the form
// errors.Is(err, http.StatusBadRequests)
func (e statusError) Is(target error) bool {
	if se, ok := target.(statusError); ok {
		if e.httpStatusCode() == se.status {
			return true
		}
	}
	return false
}

// BadRequest represents the StatusBadRequest HTTP error.
var BadRequest = statusError{http.StatusBadRequest}

// Unauthorized represents the StatusUnauthorized HTTP error.
var Unauthorized = statusError{http.StatusUnauthorized}

// PaymentRequired represents the StatusPaymentRequired HTTP error.
var PaymentRequired = statusError{http.StatusPaymentRequired}

// Forbidden represents the StatusForbidden HTTP error.
var Forbidden = statusError{http.StatusForbidden}

// NotFound represents the StatusNotFound HTTP error.
var NotFound = statusError{http.StatusNotFound}

// MethodNotAllowed represents the StatusMethodNotAllowed HTTP error.
var MethodNotAllowed = statusError{http.StatusMethodNotAllowed}

// NotAcceptable represents the StatusNotAcceptable HTTP error.
var NotAcceptable = statusError{http.StatusNotAcceptable}

// ProxyAuthRequired represents the StatusProxyAuthRequired HTTP error.
var ProxyAuthRequired = statusError{http.StatusProxyAuthRequired}

// RequestTimeout represents the StatusRequestTimeout HTTP error.
var RequestTimeout = statusError{http.StatusRequestTimeout}

// Conflict represents the StatusConflict HTTP error.
var Conflict = statusError{http.StatusConflict}

// Gone represents the StatusGone HTTP error.
var Gone = statusError{http.StatusGone}

// LengthRequired represents the StatusLengthRequired HTTP error.
var LengthRequired = statusError{http.StatusLengthRequired}

// PreconditionFailed represents the StatusPreconditionFailed HTTP error.
var PreconditionFailed = statusError{http.StatusPreconditionFailed}

// RequestEntityTooLarge represents the StatusRequestEntityTooLarge HTTP error.
var RequestEntityTooLarge = statusError{http.StatusRequestEntityTooLarge}

// RequestURITooLong represents the StatusRequestURITooLong HTTP error.
var RequestURITooLong = statusError{http.StatusRequestURITooLong}

// UnsupportedMediaType represents the StatusUnsupportedMediaType HTTP error.
var UnsupportedMediaType = statusError{http.StatusUnsupportedMediaType}

// RequestedRangeNotSatisfiable represents the StatusRequestedRangeNotSatisfiable HTTP error.
var RequestedRangeNotSatisfiable = statusError{http.StatusRequestedRangeNotSatisfiable}

// ExpectationFailed represents the StatusExpectationFailed HTTP error.
var ExpectationFailed = statusError{http.StatusExpectationFailed}

// Teapot represents the StatusTeapot HTTP error.
var Teapot = statusError{http.StatusTeapot}

// MisdirectedRequest represents the StatusMisdirectedRequest HTTP error.
var MisdirectedRequest = statusError{http.StatusMisdirectedRequest}

// UnprocessableEntity represents the StatusUnprocessableEntity HTTP error.
var UnprocessableEntity = statusError{http.StatusUnprocessableEntity}

// Locked represents the StatusLocked HTTP error.
var Locked = statusError{http.StatusLocked}

// FailedDependency represents the StatusFailedDependency HTTP error.
var FailedDependency = statusError{http.StatusFailedDependency}

// TooEarly represents the StatusTooEarly HTTP error.
var TooEarly = statusError{http.StatusTooEarly}

// UpgradeRequired represents the StatusUpgradeRequired HTTP error.
var UpgradeRequired = statusError{http.StatusUpgradeRequired}

// PreconditionRequired represents the StatusPreconditionRequired HTTP error.
var PreconditionRequired = statusError{http.StatusPreconditionRequired}

// TooManyRequests represents the StatusTooManyRequests HTTP error.
var TooManyRequests = statusError{http.StatusTooManyRequests}

// RequestHeaderFieldsTooLarge represents the StatusRequestHeaderFieldsTooLarge HTTP error.
var RequestHeaderFieldsTooLarge = statusError{http.StatusRequestHeaderFieldsTooLarge}

// UnavailableForLegalReasons represents the StatusUnavailableForLegalReasons HTTP error.
var UnavailableForLegalReasons = statusError{http.StatusUnavailableForLegalReasons}

// InternalServerError represents the StatusInternalServerError HTTP error.
var InternalServerError = statusError{http.StatusInternalServerError}

// NotImplemented represents the StatusNotImplemented HTTP error.
var NotImplemented = statusError{http.StatusNotImplemented}

// BadGateway represents the StatusBadGateway HTTP error.
var BadGateway = statusError{http.StatusBadGateway}

// ServiceUnavailable represents the StatusServiceUnavailable HTTP error.
var ServiceUnavailable = statusError{http.StatusServiceUnavailable}

// GatewayTimeout represents the StatusGatewayTimeout HTTP error.
var GatewayTimeout = statusError{http.StatusGatewayTimeout}

// HTTPVersionNotSupported represents the StatusHTTPVersionNotSupported HTTP error.
var HTTPVersionNotSupported = statusError{http.StatusHTTPVersionNotSupported}

// VariantAlsoNegotiates represents the StatusVariantAlsoNegotiates HTTP error.
var VariantAlsoNegotiates = statusError{http.StatusVariantAlsoNegotiates}

// InsufficientStorage represents the StatusInsufficientStorage HTTP error.
var InsufficientStorage = statusError{http.StatusInsufficientStorage}

// LoopDetected represents the StatusLoopDetected HTTP error.
var LoopDetected = statusError{http.StatusLoopDetected}

// NotExtended represents the StatusNotExtended HTTP error.
var NotExtended = statusError{http.StatusNotExtended}

// NetworkAuthenticationRequired represents the StatusNetworkAuthenticationRequired HTTP error.
var NetworkAuthenticationRequired = statusError{http.StatusNetworkAuthenticationRequired}
