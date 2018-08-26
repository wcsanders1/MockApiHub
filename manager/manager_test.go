package manager

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIsAPIConfig(t *testing.T) {
	fileName := "test.toml"
	result := isAPIConfig(fileName)

	assert := assert.New(t)
	assert.True(result)

	fileName = "test.exe"
	result = isAPIConfig(fileName)

	assert.False(result)

	fileName = ""
	result = isAPIConfig(fileName)

	assert.False(result)
}

func TestIsAPI(t *testing.T) {
	dirName := "testApi"
	result := isAPI(dirName)

	assert := assert.New(t)
	assert.True(result)

	dirName = "test"
	result = isAPI(dirName)

	assert.False(result)

	dirName = ""
	result = isAPI(dirName)

	assert.False(result)
}