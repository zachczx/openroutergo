package strutil

import (
	"testing"

	"github.com/eduardolat/openroutergo/internal/assert"
)

func TestRemoveTrailingSlash(t *testing.T) {
	assert.Equal(t, "test", RemoveTrailingSlash("test"))
	assert.Equal(t, "test", RemoveTrailingSlash("test/"))
	assert.Equal(t, "test/", RemoveTrailingSlash("test//"))
}

func TestRemoveTrailingSlashes(t *testing.T) {
	assert.Equal(t, "test", RemoveTrailingSlashes("test"))
	assert.Equal(t, "test", RemoveTrailingSlashes("test/"))
	assert.Equal(t, "test", RemoveTrailingSlashes("test//"))
	assert.Equal(t, "test", RemoveTrailingSlashes("test////////"))
}
