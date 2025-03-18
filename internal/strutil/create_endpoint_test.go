package strutil

import (
	"testing"

	"github.com/eduardolat/openroutergo/internal/assert"
)

func TestCreateEndpoint(t *testing.T) {
	assert.Equal(t, "https://example.com/test", CreateEndpoint("https://example.com/", "/test"))
	assert.Equal(t, "https://example.com/test", CreateEndpoint("https://example.com/", "test"))
	assert.Equal(t, "https://example.com/test/", CreateEndpoint("https://example.com/", "test/"))
	assert.Equal(t, "https://example.com/test//", CreateEndpoint("https://example.com/", "test//"))
	assert.Equal(t, "https://example.com/api/v1/test", CreateEndpoint("https://example.com/api/v1/", "/test"))
}
