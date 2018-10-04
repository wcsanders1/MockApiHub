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

	emptyURL := ""
	frags, err = GetURLFragments(emptyURL)

	assert.Nil(frags)
	assert.NotNil(err)
	assert.Error(err)
}

func TestCleanURL(t *testing.T) {
	url := "/TESt/Url/"
	result := CleanURL(url)

	assert := assert.New(t)
	assert.NotNil(result)
	assert.NotEmpty(result)
	assert.Equal("test/url", result)

	emptyURL := ""
	result = CleanURL(emptyURL)

	assert.Empty(result)
}

func TestRemoveColonFromParam(t *testing.T) {
	param := ":id"
	niceParam := RemoveColonFromParam(param)

	assert := assert.New(t)
	assert.Equal("id", niceParam)

	emptyStr := ""
	assert.Empty(RemoveColonFromParam(emptyStr))
}

func TestIsParam(t *testing.T) {
	param := ":id"
	nonParam := "id"
	emptyStr := ""

	assert := assert.New(t)
	assert.True(IsParam(param))
	assert.False(IsParam(nonParam))
	assert.False(IsParam(emptyStr))
}
