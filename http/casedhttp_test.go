package casedhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cased/cased-go"
	"github.com/stretchr/testify/assert"
)

func TestContextMiddlewareSetsContext(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", nil)
	req.Header.Add("User-Agent", "cased-test/v1")
	req.Header.Add("X-Forwarded-For", "::1")
	req.Header.Add("X-Request-ID", "1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed")
	assert.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ae := cased.GetContextFromContext(req.Context())
		expected := cased.AuditEvent{
			"location":            cased.NewSensitiveValue("::1", "ip-address"),
			"request_http_method": "POST",
			"request_id":          "1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed",
			"request_url":         "/login",
			"request_user_agent":  "cased-test/v1",
		}

		assert.Equal(t, expected, ae)
	})

	handlerToTest := ContextMiddleware(handler)

	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}

func TestContextMiddlewareDoesNotIncludeRequestIDIfNotSet(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", nil)
	req.Header.Add("User-Agent", "cased-test/v1")
	req.Header.Add("X-Forwarded-For", "::1")
	assert.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ae := cased.GetContextFromContext(req.Context())
		expected := cased.AuditEvent{
			"location":            cased.NewSensitiveValue("::1", "ip-address"),
			"request_http_method": "POST",
			"request_url":         "/login",
			"request_user_agent":  "cased-test/v1",
		}

		assert.Equal(t, expected, ae)
	})

	handlerToTest := ContextMiddleware(handler)

	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}

func TestContextMiddlewareUsesXForwardedForIfPresent(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", nil)
	req.Header.Add("User-Agent", "cased-test/v1")
	req.Header.Add("X-Forwarded-For", "1.1.1.1")
	assert.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ae := cased.GetContextFromContext(req.Context())
		expected := cased.AuditEvent{
			"location":            cased.NewSensitiveValue("1.1.1.1", "ip-address"),
			"request_http_method": "POST",
			"request_url":         "/login",
			"request_user_agent":  "cased-test/v1",
		}

		assert.NotEqual(t, req.RemoteAddr, expected["location"])
		assert.Equal(t, expected, ae)
	})

	handlerToTest := ContextMiddleware(handler)

	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}
