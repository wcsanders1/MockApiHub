package api

import (
	"net/http"
	"testing"

	"github.com/wcsanders1/MockApiHub/config"

	"github.com/stretchr/testify/assert"
)

func TestCreateAPIServer_ReturnsServer_WhenPortProvided(t *testing.T) {
	httpConfig := &config.HTTP{
		Port: 4000,
	}
	api, _ := NewAPI(&config.APIConfig{})

	result, err := createAPIServer(httpConfig, api)

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(result)
	assert.IsType(&http.Server{}, result)
}

func TestCreateAPIServer_ReturnsError_WhenPortNotProvided(t *testing.T) {
	httpConfig := &config.HTTP{
		Port: 0,
	}
	api, _ := NewAPI(&config.APIConfig{})

	result, err := createAPIServer(httpConfig, api)

	assert := assert.New(t)
	assert.Nil(result)
	assert.Error(err)
}
