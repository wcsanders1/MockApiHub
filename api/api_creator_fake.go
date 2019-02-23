package api

import (
	"net/http"

	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/wrapper"

	"github.com/stretchr/testify/mock"
)

type fakeAPICreator struct {
	mock.Mock
}

func (c *fakeAPICreator) getHandler(enforceValidJSON, allowCORS bool, headers []config.Header, dir, fileName string, file wrapper.IFileOps) func(w http.ResponseWriter, r *http.Request) {
	args := c.Called(enforceValidJSON, dir, fileName, file)
	return args.Get(0).(func(w http.ResponseWriter, r *http.Request))
}

func (c *fakeAPICreator) startAPI(defaultCert, defaultKey string, server wrapper.IServerOps, httpConfig config.HTTP) error {
	args := c.Called(defaultCert, defaultKey, server, httpConfig)
	return args.Error(0)
}
