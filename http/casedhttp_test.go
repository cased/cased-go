package casedhttp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestVerifyWebhookSignatureMiddleware(t *testing.T) {
	body := strings.NewReader(`{}`)
	req, err := http.NewRequest("POST", "/webhook", body)
	req.Header.Add(WebhookTimestampHeader, "1622047545") // some fixed time in the past
	req.Header.Add(WebhookSignatureHeader, "f6becf406ecc9d45d6c09b7994204156401ec4c39027e06d1954e0b854987b3c")
	assert.NoError(t, err)
	reached := false

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reached = true
	})

	params := &VerifyWebhookSignatureParams{
		Secret: cased.String("webhook_secret_1t57jFvmVYju00z8F4fFB8veweg"),
	}
	handlerToTest := VerifyWebhookSignatureMiddleware(handler, params)

	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)

	assert.True(t, reached)
}

func TestVerifyWebhookSignatureMiddleware_ReplayAttack(t *testing.T) {
	body := strings.NewReader(`{}`)
	req, err := http.NewRequest("POST", "/webhook", body)
	req.Header.Add(WebhookTimestampHeader, "1622047000")
	req.Header.Add(WebhookSignatureHeader, "f6becf406ecc9d45d6c09b7994204156401ec4c39027e06d1954e0b854987b3c")
	assert.NoError(t, err)
	reached := false

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reached = true
	})

	params := &VerifyWebhookSignatureParams{
		Secret:           cased.String("webhook_secret_1t57jFvmVYju00z8F4fFB8veweg"),
		TimestampExpires: time.Minute * 5,
	}
	handlerToTest := VerifyWebhookSignatureMiddleware(handler, params)

	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)

	assert.False(t, reached, "expected error to be returned")
}
