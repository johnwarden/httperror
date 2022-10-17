package httperror

import (
	"bytes"
	"encoding/json"
	"mime"
	"net/http"
	"strconv"
)

// const contentTypeHTML = "text/html"
const (
	contentTypeTextPlain = "text/plain"
	contentTypeText      = "text"
	contentTypeJSON      = "application/json"
)

// ErrorHandler handles an error.
type ErrorHandler = func(w http.ResponseWriter, err error)

// DefaultErrorHandler writes a reasonable default error response, using the status
// code from the error if it can be extracted (see [StatusCode]), or 500 by
// default, using the content type from from w.Header(), or text/html by
// default, and using any public message (see [PublicErrorf] and [Public].)
func DefaultErrorHandler(w http.ResponseWriter, e error) {
	s := StatusCode(e)
	w.WriteHeader(s)

	var b bytes.Buffer
	b.WriteString(http.StatusText(s))
	if s := PublicMessage(e); s != "" {
		b.WriteString(": ")
		b.WriteString(s)
	}

	WriteResponse(w, s, b.Bytes())
}

// WriteResponse writes a reasonable default error response given the status
// code and optional error message. The default error handler
// [DefaultErrorHandler] calls this method after extracting the status code and any
// public error message.
func WriteResponse(w http.ResponseWriter, s int, m []byte) {
	contentType := responseContentType(w)

	switch contentType {
	case contentTypeJSON:
		writeJsonErrorBody(w, s, m)
	case contentTypeTextPlain:
		writePlainTextErrorBody(w, s, m)
	case contentTypeText:
		writePlainTextErrorBody(w, s, m)
	default:
		writeHtmlErrorBody(w, s, m)
	}
}

func writeHtmlErrorBody(w http.ResponseWriter, s int, m []byte) {
	_, _ = w.Write([]byte(`<html><head><meta http-equiv="Content-Type" content="text/html; charset=UTF-8"><title>`))
	_, _ = w.Write([]byte(`Error `))
	_, _ = w.Write([]byte(strconv.Itoa(s)))
	_, _ = w.Write([]byte(`</title></head><body>`))
	_, _ = w.Write([]byte(m))
	_, _ = w.Write([]byte("</body></html>\n"))
}

func writePlainTextErrorBody(w http.ResponseWriter, s int, m []byte) {
	_, _ = w.Write([]byte(strconv.Itoa(s)))
	_, _ = w.Write([]byte(` `))
	_, _ = w.Write([]byte(m))
	_, _ = w.Write([]byte("\n"))
}

// jsonError prints an error using general guidelines from
// https://github.com/omniti-labs/jsend
func writeJsonErrorBody(w http.ResponseWriter, s int, m []byte) {
	response := jsonhttperror{Status: "error", Message: string(m), Code: s}
	json, _ := json.Marshal(response) // No error handling for error handling

	_, _ = w.Write(json)
	_, _ = w.Write([]byte("\n"))
}

type jsonhttperror struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// responseContentType extracts the content type from the response writer, if
// the Content-Type header has been set. It does *not* return the entire
// content type header -- only the media type part (e.g. "text/html" but not
// "text/html; charset=UTF-8").
func responseContentType(w http.ResponseWriter) string {
	var contentType string
	if cts, ok := w.Header()["Content-Type"]; ok {
		contentType, _, _ = mime.ParseMediaType(cts[0])
	}
	return contentType
}
