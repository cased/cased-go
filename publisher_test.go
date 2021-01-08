package cased

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCurrentPublisher(t *testing.T) {
	cp := CurrentPublisher()

	assert.NotNil(t, cp)
}

func TestSetPublisher(t *testing.T) {
	cp := CurrentPublisher()
	nc := NewPublisher(
		WithDebug(true),
	)
	SetPublisher(nc)
	defer func(oc Publisher) {
		SetPublisher(oc)
	}(cp)

	assert.False(t, cp.Options().Debug)
	assert.True(t, CurrentPublisher().Options().Debug)
}

func TestPublishURLFromEnvironment(t *testing.T) {
	defer restoreEnv("CASED_PUBLISH_URL")()
	os.Setenv("CASED_PUBLISH_URL", "https://publish.internal.host")

	c := NewPublisher()

	assert.Equal(t, "https://publish.internal.host", c.Options().PublishURL)
}

func TestPublishKeyFromEnvironment(t *testing.T) {
	defer restoreEnv("CASED_PUBLISH_KEY")()
	os.Setenv("CASED_PUBLISH_KEY", "publish_live_1mXiiYCvpmdKuKbV1tukt5gJwYO")

	c := NewPublisher()

	assert.Equal(t, "publish_live_1mXiiYCvpmdKuKbV1tukt5gJwYO", c.Options().PublishKey)
}

func TestDebugFromEnvironment(t *testing.T) {
	defer restoreEnv("CASED_DEBUG")()
	os.Setenv("CASED_DEBUG", "true")

	c := NewPublisher()

	assert.True(t, c.Options().Debug)
}

func TestSilenceFromEnvironment(t *testing.T) {
	defer restoreEnv("CASED_SILENCE")()
	os.Setenv("CASED_SILENCE", "true")

	c := NewPublisher()

	assert.True(t, c.Options().Silence, c.Options())
}

func TestHTTPTimeoutFromEnvironment(t *testing.T) {
	defer restoreEnv("CASED_HTTP_TIMEOUT")()
	os.Setenv("CASED_HTTP_TIMEOUT", "3s")

	c := NewPublisher()

	assert.Equal(t, 3*time.Second, c.Options().HTTPTimeout)
}

func TestWithPublishURL(t *testing.T) {
	nc := NewPublisher(
		WithPublishURL("https://publish.internal.host"),
	)

	assert.Equal(t, "https://publish.internal.host", nc.Options().PublishURL)
}

func TestWithPublishKey(t *testing.T) {
	nc := NewPublisher(
		WithPublishKey("publish_live_1mXgDV0ZeF1MQ3ge3GpGD00jDjC"),
	)

	assert.Equal(t, "publish_live_1mXgDV0ZeF1MQ3ge3GpGD00jDjC", nc.Options().PublishKey)
}

func TestWithDebug(t *testing.T) {
	nc := NewPublisher(
		WithDebug(true),
	)

	assert.True(t, nc.Options().Debug)
}

func TestWithSilence(t *testing.T) {
	nc := NewPublisher(
		WithSilence(true),
	)

	assert.True(t, nc.Options().Silence)
}

func TestWithHTTPClient(t *testing.T) {
	hc := &http.Client{
		Timeout: 1 * time.Second,
	}
	nc := NewPublisher(
		WithHTTPClient(hc),
	)

	assert.Same(t, hc, nc.Options().HTTPClient)
}

func TestWithHTTPTransport(t *testing.T) {
	ht := &http.Transport{}
	nc := NewPublisher(
		WithHTTPTransport(ht),
	)

	assert.Same(t, ht, nc.Options().HTTPTransport)
}

func TestWithHTTPTimeout(t *testing.T) {
	nc := NewPublisher(
		WithHTTPTimeout(30 * time.Second),
	)

	assert.Equal(t, 30*time.Second, nc.Options().HTTPTimeout)
}

func TestWithTransport(t *testing.T) {
	transport := NewNoopHTTPTransport()
	nc := NewPublisher(
		WithTransport(transport),
	)

	assert.Same(t, transport, nc.Options().Transport)
}

func restoreEnv(key string) func() {
	v := os.Getenv(key)
	os.Clearenv()
	return func() {
		os.Setenv(key, v)
	}
}
