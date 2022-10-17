package httperror_test

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"

	"github.com/johnwarden/httperror"
	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	{
		var e error = httperror.TooManyRequests
		assert.Equal(t, http.StatusTooManyRequests, httperror.StatusCode(e), "httperror.StatusCode")
		assert.Equal(t, "429 Too Many Requests", e.Error())
	}

	{
		e := errors.New("test error")
		assert.Equal(t, httperror.StatusCode(e), http.StatusInternalServerError, "StatusCode defaults to Internal Server Error for non-http errors")
	}

	{
		errMissingParam := errors.New("missing parameter 'foo'")

		e := httperror.Wrap(errMissingParam, http.StatusBadRequest)

		assert.Equal(t, "400 Bad Request: missing parameter 'foo'", e.Error())

		assert.True(t, errors.Is(e, httperror.BadRequest))
		assert.False(t, errors.Is(httperror.BadRequest, e))
		assert.False(t, errors.Is(e, httperror.NotFound))
		assert.False(t, errors.Is(e, httperror.InternalServerError))
		assert.True(t, errors.Is(e, errMissingParam))

		assert.Equal(t, http.StatusBadRequest, httperror.StatusCode(e), "HTTPStatusCode()")

		assert.Equal(t, "", httperror.PublicMessage(e))

	}

	{
		e := httperror.Errorf(http.StatusBadRequest, "missing parameter 'foo'")

		assert.Equal(t, "400 Bad Request: missing parameter 'foo'", e.Error())

		assert.True(t, errors.Is(e, httperror.BadRequest))
		assert.True(t, errors.Is(e, e))
		assert.False(t, errors.Is(httperror.BadRequest, e))
		assert.False(t, errors.Is(e, httperror.NotFound))

		assert.Equal(t, http.StatusBadRequest, httperror.StatusCode(e))

		assert.Equal(t, "", httperror.PublicMessage(e))

	}

	{
		e := httperror.PublicErrorf(http.StatusBadRequest, "missing parameter 'foo'")

		assert.Equal(t, "400 Bad Request: missing parameter 'foo'", e.Error())

		assert.True(t, errors.Is(e, httperror.BadRequest))
		assert.True(t, errors.Is(e, e))
		assert.False(t, errors.Is(httperror.BadRequest, e))
		assert.False(t, errors.Is(e, httperror.NotFound))

		assert.Equal(t, http.StatusBadRequest, httperror.StatusCode(e))
		assert.Equal(t, "missing parameter 'foo'", httperror.PublicMessage(e))
	}
}
