package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestRespondWithNotReadyStatus(t *testing.T) {
	// Arrange.
	logger := new(zerolog.Logger)
	reason := "test reason"
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		RespondWithNotReadyStatus(logger, w, reason)
	}

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)

	// Act.
	handler.ServeHTTP(rr, request)

	// Assert.
	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
	assert.Equal(t, reason, rr.Body.String())
}

func TestRespondWithNotAuthorizedError(t *testing.T) {
	// Arrange.
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		RespondWithNotAuthorizedError(w)
	}

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)

	// Act.
	handler.ServeHTTP(rr, request)

	// Assert.
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Unauthorized", strings.TrimSpace(rr.Body.String()))
}

func TestRespondWithBadRequestError(t *testing.T) {
	// Arrange.
	reason := "test reason"
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		RespondWithBadRequestError(w, reason)
	}

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)

	// Act.
	handler.ServeHTTP(rr, request)

	// Assert.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, reason, strings.TrimSpace(rr.Body.String()))
}

func TestRespondWithInternalServerError(t *testing.T) {
	// Arrange.
	reason := "test reason"
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		RespondWithInternalServerError(w, reason)
	}

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)

	// Act.
	handler.ServeHTTP(rr, request)

	// Assert.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, reason, strings.TrimSpace(rr.Body.String()))
}

func TestRespondWithJsonObject(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Arrange.
	logger := new(zerolog.Logger)
	object := &Person{
		Name: "John",
		Age:  123,
	}

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		RespondWithJsonObject(logger, w, object)
	}

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)

	// Act.
	handler.ServeHTTP(rr, request)

	// Assert.
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, MimeTypeApplicationJson, rr.Header().Get(HttpHeaderContentType))
	assert.Equal(t, `{"name":"John","age":123}`, rr.Body.String())
}
