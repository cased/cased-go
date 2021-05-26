package casedhttp

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cased/cased-go"
)

const (
	WebhookTimestampHeader = "X-Cased-Timestamp"
	WebhookSignatureHeader = "X-Cased-Signature"
)

var (
	WebhookSecret                     = os.Getenv("CASED_WEBHOOK_SECRET")
	WebhookTimestampExpiredError      = errors.New("webhook timestamp expired")
	WebhookSignatureVerificationError = errors.New("webhook computed signature does not match signature sent with webhook")
)

type VerifyWebhookSignatureParams struct {
	// Secret used to compute the HMAC signature. Optional if secret is provided
	// by the CASED_WEBHOOK_SECRET environment variable.
	Secret *string

	// TimestampExpires if provided will reject webhook requests that are
	// delivered after specified duration. Useful to prevent replay attacks. Each
	// webhook attempt delivered from Cased will provide a new timestamp.
	TimestampExpires time.Duration
}

func VerifyWebhookSignature(req *http.Request, params *VerifyWebhookSignatureParams) error {
	if req.Method != http.MethodPost {
		return errors.New("post request expected")
	}

	secret := params.Secret
	if secret == nil {
		secret = &WebhookSecret
	}

	// Check to see if timestamp expiration is configured and enforce it.
	timestamp := req.Header.Get(WebhookTimestampHeader)
	if params.TimestampExpires > 0 {
		i, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return err
		}

		// Time from header
		tm := time.Unix(i, 0)
		expires := time.Now().Add(-params.TimestampExpires)
		if tm.Before(expires) {
			return WebhookTimestampExpiredError
		}
	}

	signature := req.Header.Get(WebhookSignatureHeader)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil
	}

	basestring := fmt.Sprintf("%s.%s", timestamp, string(body))
	mac := hmac.New(sha256.New, []byte(*secret))
	if _, err = mac.Write([]byte(basestring)); err != nil {
		return err
	}
	computed := hex.EncodeToString(mac.Sum(nil))

	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Compare the computed signature with signature sent with webhook
	if signature != computed {
		return WebhookSignatureVerificationError
	}

	return nil
}

func VerifyWebhookSignatureMiddleware(next http.Handler, params *VerifyWebhookSignatureParams) http.Handler {
	if params.Secret == nil && WebhookSecret == "" {
		panic("Must set CASED_WEBHOOK_SECRET or provide VerifyWebhookSignatureParams to VerifyWebhookSignatureMiddleware")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := VerifyWebhookSignature(req, params); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, req)
	})
}

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
