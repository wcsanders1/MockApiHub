package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/log"
	"github.com/wcsanders1/MockApiHub/wrapper"

	"github.com/stretchr/testify/assert"
)

func TestStartAPI_ReturnsNil_WhenStartNoTLSSuccessful(t *testing.T) {
	fakeServer := wrapper.NewFakeServerOps()
	fakeServer.On("ListenAndServe").Return(nil)
	httpConfig := config.HTTP{
		UseTLS: false,
	}
	creator := creator{
		log: log.GetFakeLogger(),
	}

	err := creator.startAPI("defaultCert", "defaultKey", fakeServer, httpConfig)
	fakeServer.WaitForListenAndServe()
	fakeServer.AssertCalled(t, "ListenAndServe")
	fakeServer.AssertNotCalled(t, "ListenAndServeTLS")
	assert.NoError(t, err)
}

func TestStartAPI_ReturnsNil_WhenListenAndServeFails(t *testing.T) {
	fakeServer := wrapper.NewFakeServerOps()
	fakeServer.On("ListenAndServe").Return(errors.New(""))
	httpConfig := config.HTTP{
		UseTLS: false,
	}
	creator := creator{
		log: log.GetFakeLogger(),
	}

	err := creator.startAPI("defaultCert", "defaultKey", fakeServer, httpConfig)
	fakeServer.WaitForListenAndServe()
	fakeServer.AssertCalled(t, "ListenAndServe")
	fakeServer.AssertNotCalled(t, "ListenAndServeTLS")
	assert.NoError(t, err)
}

func TestStartAPI_ReturnsNil_WhenListenAndServeTLSSSuccessful(t *testing.T) {
	defaultCert := "defaultCert"
	defaultKey := "defaultKey"
	fakeServer := wrapper.NewFakeServerOps()
	fakeServer.On("ListenAndServeTLS", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	httpConfig := config.HTTP{
		UseTLS: true,
	}
	creator := creator{
		log: log.GetFakeLogger(),
	}

	err := creator.startAPI(defaultCert, defaultKey, fakeServer, httpConfig)
	fakeServer.WaitForListenAndServe()
	fakeServer.AssertNotCalled(t, "ListenAndServe")
	fakeServer.AssertCalled(t, "ListenAndServeTLS", defaultCert, defaultKey)
	assert.NoError(t, err)
}

func TestStartAPI_ReturnsNil_WhenListenAndServeTLSFails(t *testing.T) {
	defaultCert := "defaultCert"
	defaultKey := "defaultKey"
	fakeServer := wrapper.NewFakeServerOps()
	fakeServer.On("ListenAndServeTLS", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New(""))
	httpConfig := config.HTTP{
		UseTLS: true,
	}
	creator := creator{
		log: log.GetFakeLogger(),
	}

	err := creator.startAPI(defaultCert, defaultKey, fakeServer, httpConfig)
	fakeServer.WaitForListenAndServe()
	fakeServer.AssertNotCalled(t, "ListenAndServe")
	fakeServer.AssertCalled(t, "ListenAndServeTLS", defaultCert, defaultKey)
	assert.NoError(t, err)
}

func TestStartAPI_ReturnsError_WhenGetCertAndKeyFails(t *testing.T) {
	fakeServer := wrapper.NewFakeServerOps()
	fakeServer.On("ListenAndServeTLS", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New(""))
	httpConfig := config.HTTP{
		UseTLS:   true,
		CertFile: "testCert",
	}
	creator := creator{
		log: log.GetFakeLogger(),
	}

	err := creator.startAPI("defaultCert", "defaultKey", fakeServer, httpConfig)
	fakeServer.AssertNotCalled(t, "ListenAndServe")
	fakeServer.AssertNotCalled(t, "ListenAndServeTLS", mock.Anything, mock.Anything)
	assert.Error(t, err)
}

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
