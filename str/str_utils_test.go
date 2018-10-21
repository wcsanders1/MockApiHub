package str

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPort(t *testing.T) {
	port := GetPort(5000)

	assert.Equal(t, ":5000", port)
}

func TestGetURLFragments_ReturnsFragments_WhenProvidedURL(t *testing.T) {
	firstFrag := "test"
	secondFrag := "url"

	result, err := GetURLFragments(fmt.Sprintf("%s/%s", firstFrag, secondFrag))

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(result)
	assert.IsType([]string{}, result)
	assert.NotEmpty(result)
	assert.Equal(2, len(result))
	assert.Equal(firstFrag, result[0])
	assert.Equal(secondFrag, result[1])
}

func TestGetURLFragmentsReturnsError_WhenProvidedNothing(t *testing.T) {
	result, err := GetURLFragments("")

	assert := assert.New(t)
	assert.Nil(result)
	assert.Error(err)
}

func TestCleanURL_ReturnsCleanedURL_WhenProvidedURL(t *testing.T) {
	result := CleanURL("/TESt/Url/")

	assert := assert.New(t)
	assert.NotNil(result)
	assert.NotEmpty(result)
	assert.Equal("test/url", result)
}

func TestCleanURL_ReturnsNothing_WhenProvidedNothing(t *testing.T) {
	assert.Empty(t, CleanURL(""))
}

func TestRemoveColonFromParam_ReturnsColonlessParam_WhenProvidedParam(t *testing.T) {
	assert.Equal(t, "id", RemoveColonFromParam(":id"))
}

func TestRemoveColonFromParam_ReturnsNothing_WhenProvidedNothing(t *testing.T) {
	assert.Empty(t, RemoveColonFromParam(""))
}

func TestIsParam_ReturnsTrue_WhenProvidedParam(t *testing.T) {
	assert.True(t, IsParam(":id"))
}

func TestIsParam_ReturnsFalse_WhenProvidedNonParam(t *testing.T) {
	assert.False(t, IsParam("id"))
}

func TestIsParam_ReturnsFalse_WhenProvidedNothing(t *testing.T) {
	assert.False(t, IsParam(""))
}
