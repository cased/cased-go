package casedhttp

import (
	"context"
	"net/http"

	"github.com/cased/cased-go"
)

// ContextMiddleware ...
func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		location := req.RemoteAddr
		if forwardedIP := req.Header.Get("X-Forwarded-For"); forwardedIP != "" {
			location = forwardedIP
		}

		ae := cased.AuditEvent{
			"location":            cased.NewSensitiveValue(location, "ip-address"),
			"request_url":         req.URL.String(),
			"request_http_method": req.Method,
			"request_user_agent":  req.Header.Get("User-Agent"),
		}

		if requestID := req.Header.Get("X-Request-ID"); requestID != "" {
			ae["request_id"] = requestID
		}

		ctx := context.WithValue(req.Context(), cased.ContextKey, ae)
		req = req.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
