package json

import (
	"errors"
	"os"
	"testing"

	"github.com/wcsanders1/MockApiHub/wrapper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var goodJSON = []byte(`{
	"JSON": "good",
	"test": "good"
}`)

var badJSON = []byte(`{
	"JSON": "bad,
	"test": "good"
}`)

func TestGetJSON_ReturnsJSON_WhenRetrievedFromFile(t *testing.T) {
	file := "testpath"
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), nil)
	fileOps.On("ReadAll", mock.AnythingOfType("*os.File")).Return(goodJSON, nil)

	result, err := GetJSON(file, fileOps)

	assert := assert.New(t)
	assert.NotNil(result)
	assert.Error(err)
	fileOps.AssertCalled(t, "Open", file)
	fileOps.AssertCalled(t, "ReadAll", mock.AnythingOfType("*os.File"))
}

func TestGetJSON_ReturnsError_WhenOpenFileFails(t *testing.T) {
	file := "testpath"
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), errors.New(""))

	result, err := GetJSON(file, fileOps)

	assert := assert.New(t)
	assert.Error(err)
	assert.Nil(result)
	fileOps.AssertCalled(t, "Open", file)
	fileOps.AssertNotCalled(t, "ReadAll", mock.AnythingOfType("*os.File"))
}

func TestGetJSON_ReturnsError_WhenReadFileFails(t *testing.T) {
	file := "testpath"
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), nil)
	fileOps.On("ReadAll", mock.AnythingOfType("*os.File")).Return(goodJSON, errors.New(""))

	result, err := GetJSON(file, fileOps)

	assert := assert.New(t)
	assert.Error(err)
	assert.Nil(result)
	fileOps.AssertCalled(t, "Open", file)
	fileOps.AssertCalled(t, "ReadAll", mock.AnythingOfType("*os.File"))
}

func TestGetJSON_ReturnsError_WhenJSONInvalid(t *testing.T) {
	file := "testpath"
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), nil)
	fileOps.On("ReadAll", mock.AnythingOfType("*os.File")).Return(badJSON, nil)

	result, err := GetJSON(file, fileOps)

	assert := assert.New(t)
	assert.Error(err)
	assert.Nil(result)
	fileOps.AssertCalled(t, "Open", file)
	fileOps.AssertCalled(t, "ReadAll", mock.AnythingOfType("*os.File"))
}

func TestIsValidJSON_ReturnsTrue_WhenJSONValid(t *testing.T) {
	assert.True(t, isValidJSON(goodJSON))
}

func TestIsValidJSON_ReturnsFalse_WhenJSONInvalid(t *testing.T) {
	assert.False(t, isValidJSON(badJSON))
}

func TestIsValidJSON_ReturnsFalse_WhenGivenNothing(t *testing.T) {
	assert.False(t, isValidJSON([]byte("")))
}
