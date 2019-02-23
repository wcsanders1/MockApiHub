package api

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/fake"
	"github.com/wcsanders1/MockApiHub/log"
	"github.com/wcsanders1/MockApiHub/wrapper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var goodJSON = []byte(`{
	"JSON": "good",
	"test": "good"
}`)

func TestNewCreator_ReturnsCreator_WhenCalled(t *testing.T) {
	result := newCreator(log.GetFakeLogger())

	assert := assert.New(t)
	assert.NotNil(result)
	assert.IsType(&creator{}, result)
}

func TestGetHandler_ReturnsHandler_WhenEnforceJSONFalse(t *testing.T) {
	creator := creator{
		log: log.GetFakeLogger(),
	}

	result := creator.getHandler(false, false, nil, "testDir", "testFile", &wrapper.FakeFileOps{})

	assert := assert.New(t)
	assert.NotNil(result)
	assert.IsType(func(w http.ResponseWriter, r *http.Request) {}, result)
}

func TestGetHandler_ReturnsHandler_WhenEnforceJSONTrue(t *testing.T) {
	creator := creator{
		log: log.GetFakeLogger(),
	}

	result := creator.getHandler(true, false, nil, "testDir", "testFile", &wrapper.FakeFileOps{})

	assert := assert.New(t)
	assert.NotNil(result)
	assert.IsType(func(w http.ResponseWriter, r *http.Request) {}, result)
}

func TestGetJSONHandler_ReturnsHandler_WhenCalled(t *testing.T) {
	fileOps := wrapper.FakeFileOps{}
	logger := log.GetFakeLogger()
	funcResult := getJSONHandler("test", nil, &fileOps, logger, false)

	assert.NotNil(t, funcResult)
}

func TestJSONHandler_Writes_OnSuccess(t *testing.T) {
	path := "test/path"
	fileOps := wrapper.FakeFileOps{}
	logger := log.GetFakeLogger()
	fileOps.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), nil)
	fileOps.On("ReadAll", mock.AnythingOfType("*os.File")).Return(goodJSON, nil)
	funcResult := getJSONHandler(path, nil, &fileOps, logger, false)
	w := fake.ResponseWriter{}
	w.On("WriteHeader", mock.AnythingOfType("int")).Return(1)
	w.On("Write", mock.AnythingOfType("[]uint8")).Return(1, nil)
	request, _ := http.NewRequest("GET", "test/url", nil)

	funcResult(&w, request)

	fileOps.AssertCalled(t, "Open", path)
	fileOps.AssertCalled(t, "ReadAll", mock.AnythingOfType("*os.File"))
	w.AssertCalled(t, "Write", mock.AnythingOfType("[]uint8"))
}

func TestJSONHandler_WritesError_OnFailure(t *testing.T) {
	path := "test/path"
	fileOps := wrapper.FakeFileOps{}
	logger := log.GetFakeLogger()
	fileOps.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), nil)
	fileOps.On("ReadAll", mock.AnythingOfType("*os.File")).Return([]byte{}, errors.New(""))
	funcResult := getJSONHandler(path, nil, &fileOps, logger, false)
	w := fake.ResponseWriter{}
	w.On("WriteHeader", mock.AnythingOfType("int")).Return(1)
	w.On("Write", mock.AnythingOfType("[]uint8")).Return(1, nil)
	request, _ := http.NewRequest("GET", "test/url", nil)

	funcResult(&w, request)

	fileOps.AssertCalled(t, "Open", path)
	fileOps.AssertCalled(t, "ReadAll", mock.AnythingOfType("*os.File"))
	w.AssertCalled(t, "Write", mock.AnythingOfType("[]uint8"))
	w.AssertCalled(t, "WriteHeader", http.StatusInternalServerError)
}

func TestGetGeneralHanlder_ReturnsFunc_WhenCalled(t *testing.T) {
	fileOps := wrapper.FakeFileOps{}
	logger := log.GetFakeLogger()
	funcResult := getGeneralHandler("test", nil, &fileOps, logger, false)

	assert.NotNil(t, funcResult)
}

func TestGeneralHandler_Writes_OnSuccess(t *testing.T) {
	path := "test/path"
	fileOps := wrapper.FakeFileOps{}
	logger := log.GetFakeLogger()
	fileOps.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), nil)
	fileOps.On("ReadAll", mock.AnythingOfType("*os.File")).Return(goodJSON, nil)
	funcResult := getGeneralHandler(path, nil, &fileOps, logger, false)
	w := fake.ResponseWriter{}
	w.On("WriteHeader", mock.AnythingOfType("int")).Return(1)
	w.On("Write", mock.AnythingOfType("[]uint8")).Return(1, nil)
	request, _ := http.NewRequest("GET", "test/url", nil)

	funcResult(&w, request)

	fileOps.AssertCalled(t, "Open", path)
	fileOps.AssertCalled(t, "ReadAll", mock.AnythingOfType("*os.File"))
	w.AssertCalled(t, "Write", mock.AnythingOfType("[]uint8"))
}

func TestGeneralHandler_WritesError_WhenReadFails(t *testing.T) {
	path := "test/path"
	fileOps := wrapper.FakeFileOps{}
	logger := log.GetFakeLogger()
	fileOps.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), nil)
	fileOps.On("ReadAll", mock.AnythingOfType("*os.File")).Return([]byte{}, errors.New(""))
	funcResult := getGeneralHandler(path, nil, &fileOps, logger, false)
	w := fake.ResponseWriter{}
	w.On("WriteHeader", mock.AnythingOfType("int")).Return(1)
	w.On("Write", mock.AnythingOfType("[]uint8")).Return(1, nil)
	request, _ := http.NewRequest("GET", "test/url", nil)

	funcResult(&w, request)

	fileOps.AssertCalled(t, "Open", path)
	fileOps.AssertCalled(t, "ReadAll", mock.AnythingOfType("*os.File"))
	w.AssertCalled(t, "Write", mock.AnythingOfType("[]uint8"))
	w.AssertCalled(t, "WriteHeader", http.StatusInternalServerError)
}

func TestGeneralHandler_WritesError_WhenFileOpenFails(t *testing.T) {
	path := "test/path"
	fileOps := wrapper.FakeFileOps{}
	logger := log.GetFakeLogger()
	fileOps.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), errors.New(""))
	funcResult := getGeneralHandler(path, nil, &fileOps, logger, false)
	w := fake.ResponseWriter{}
	w.On("WriteHeader", mock.AnythingOfType("int")).Return(1)
	w.On("Write", mock.AnythingOfType("[]uint8")).Return(1, nil)
	request, _ := http.NewRequest("GET", "test/url", nil)

	funcResult(&w, request)

	fileOps.AssertCalled(t, "Open", path)
	fileOps.AssertNotCalled(t, "ReadAll", mock.Anything)
	w.AssertCalled(t, "Write", mock.AnythingOfType("[]uint8"))
	w.AssertCalled(t, "WriteHeader", http.StatusInternalServerError)
}

func TestWriteError_WritesError_WhenCalled(t *testing.T) {
	err := errors.New("test error")
	w := fake.ResponseWriter{}
	w.On("WriteHeader", mock.AnythingOfType("int")).Return(1)
	w.On("Write", mock.AnythingOfType("[]uint8")).Return(1, nil)

	writeError(err, &w)

	expectedError := []byte(err.Error())
	w.AssertCalled(t, "WriteHeader", http.StatusInternalServerError)
	w.AssertCalled(t, "Write", expectedError)
}

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
