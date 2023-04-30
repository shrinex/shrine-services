package verify

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsernameValid(t *testing.T) {
	assert.False(t, UsernameValid("aka"))
	assert.True(t, UsernameValid("saber"))
	assert.False(t, UsernameValid("berserker"))
}

func TestUrlValid(t *testing.T) {
	assert.True(t, UrlValid("path/to/sth"))
	assert.True(t, UrlValid("path/to/sth/a.jpg"))
	assert.True(t, UrlValid("path/to/sth/a.jpg?a=b"))
	assert.True(t, UrlValid("path/to/sth/a-b"))
}
