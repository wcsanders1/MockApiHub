package manager

import (
	"testing"
	"github.com/stretchr/testify/assert"

	"MockApiHub/api"
	"MockApiHub/config"
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

func TestApiByPortExists( t *testing.T) {
	port := 4000
	mgr := Manager {
		apis: make(map[string]*api.API),
	}

	cfg := &config.APIConfig {
		HTTP: config.HTTP {
			Port: port,
		},
	}

	mgr.apis["test"], _ = api.NewAPI(cfg)
	result := mgr.apiByPortExists(port)

	assert := assert.New(t)
	assert.True(result)

	result = mgr.apiByPortExists(port + 1)
	assert.False(result)
}