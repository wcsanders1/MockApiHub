package str

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetPort(t *testing.T) {
	port := GetPort(5000)
	
	assert.Equal(t, ":5000", port)
}

func TestGetURLFragments(t *testing.T) {
	url := "test/url"
	frags, err := GetURLFragments(url)
	
	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(frags)
	assert.NotEmpty(frags)
	assert.Equal(2, len(frags))
	assert.Equal("test", frags[0])
	assert.Equal("url", frags[1])
}