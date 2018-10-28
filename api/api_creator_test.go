package api

import (
	"testing"

	"github.com/wcsanders1/MockApiHub/config"

	"github.com/stretchr/testify/assert"
)

func TestGetCertAndKeyFileNames_ReturnsDefaultFileNames_WhenConfigHasNone(t *testing.T) {
	defaultCert := "defaultCert"
	defaultKey := "defaultKey"

	cert, key, err := getCertAndKeyFileNames(defaultCert, defaultKey, config.HTTP{})

	assert := assert.New(t)
	assert.Equal(defaultCert, cert)
	assert.Equal(defaultKey, key)
	assert.NoError(err)
}

func TestGetCertAndKeyFileNames_ReturnsConfigNames_WhenProvided(t *testing.T) {
	configCert := "configCert"
	configKey := "configKey"
	httpConfig := config.HTTP{
		CertFile: configCert,
		KeyFile:  configKey,
	}

	cert, key, err := getCertAndKeyFileNames("defaultCert", "defaultKey", httpConfig)

	assert := assert.New(t)
	assert.Equal(configCert, cert)
	assert.Equal(configKey, key)
	assert.NoError(err)
}

func TestGetCertAndKeyFileNames_ReturnsError_WhenConfigProvidesKeyButNotCert(t *testing.T) {
	configKey := "configKey"
	httpConfig := config.HTTP{
		KeyFile: configKey,
	}

	cert, key, err := getCertAndKeyFileNames("defaultCert", "defaultKey", httpConfig)

	assert := assert.New(t)
	assert.Empty(cert)
	assert.Empty(key)
	assert.Error(err)
}

func TestGetCertAndKeyFileNames_ReturnsError_WhenConfigProvidesCertButNotKey(t *testing.T) {
	configCert := "configCert"
	httpConfig := config.HTTP{
		CertFile: configCert,
	}

	cert, key, err := getCertAndKeyFileNames("defaultCert", "defaultKey", httpConfig)

	assert := assert.New(t)
	assert.Empty(cert)
	assert.Empty(key)
	assert.Error(err)
}
