package strutil

import (
	"testing"

	"github.com/eduardolat/openroutergo/internal/assert"
)

func TestRemoveLeadingSlash(t *testing.T) {
	assert.Equal(t, "test", RemoveLeadingSlash("test"))
	assert.Equal(t, "test", RemoveLeadingSlash("/test"))
	assert.Equal(t, "/test", RemoveLeadingSlash("//test"))
}

func TestRemoveLeadingSlashes(t *testing.T) {
	assert.Equal(t, "test", RemoveLeadingSlashes("test"))
	assert.Equal(t, "test", RemoveLeadingSlashes("/test"))
	assert.Equal(t, "test", RemoveLeadingSlashes("//test"))
	assert.Equal(t, "test", RemoveLeadingSlashes("////////test"))
}
