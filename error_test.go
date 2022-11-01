package httperror_test

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"

	"github.com/johnwarden/httperror/v2"
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
		e := httperror.New(http.StatusBadRequest, "missing parameter 'foo'")

		assert.Equal(t, "400 Bad Request: missing parameter 'foo'", e.Error())

		assert.True(t, errors.Is(e, httperror.BadRequest))
		assert.True(t, errors.Is(e, e))
		assert.False(t, errors.Is(httperror.BadRequest, e))
		assert.False(t, errors.Is(e, httperror.NotFound))

		assert.Equal(t, http.StatusBadRequest, httperror.StatusCode(e))

		assert.Equal(t, "", httperror.PublicMessage(e))

	}

	{
		e := httperror.NewPublic(http.StatusBadRequest, "missing parameter 'foo'")

		assert.Equal(t, "400 Bad Request: missing parameter 'foo'", e.Error())

		assert.True(t, errors.Is(e, httperror.BadRequest))
		assert.True(t, errors.Is(e, e))
		assert.False(t, errors.Is(httperror.BadRequest, e))
		assert.False(t, errors.Is(e, httperror.NotFound))

		assert.Equal(t, http.StatusBadRequest, httperror.StatusCode(e))
		assert.Equal(t, "missing parameter 'foo'", httperror.PublicMessage(e))
	}

	{
		e := httperror.New(http.StatusBadRequest, "")
		assert.Equal(t, e, httperror.BadRequest, "New with empty string")

		e = httperror.Errorf(http.StatusBadRequest, "%s", "")
		assert.Equal(t, e, httperror.BadRequest, "Errorf with empty string")
		assert.NotEqual(t, e, httperror.NotFound)

		e = httperror.New(http.StatusBadRequest, "What??")
		assert.NotEqual(t, e, httperror.BadRequest)
		assert.True(t, errors.Is(e, httperror.BadRequest))
	}
}
