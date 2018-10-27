package api

import (
	"fmt"
	"net/http"

	"github.com/wcsanders1/MockApiHub/constants"
	"github.com/wcsanders1/MockApiHub/json"
	"github.com/wcsanders1/MockApiHub/wrapper"

	"github.com/sirupsen/logrus"
)

type (
	// IHandlerManager provides functionality to get a handler for a mock API
	IHandlerManager interface {
		GetHandler(enforceValidJSON bool, dir, fileName string, file wrapper.IFileOps, logger *logrus.Entry) func(w http.ResponseWriter, r *http.Request)
	}

	// HandlerManager implements IAPIHandlerManager
	HandlerManager struct{}
)

// GetHandler returns an API handler
func (h HandlerManager) GetHandler(enforceValidJSON bool, dir, fileName string, file wrapper.IFileOps, logger *logrus.Entry) func(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("%s/%s/%s", constants.APIDir, dir, fileName)
	if enforceValidJSON {
		return getJSONHandler(path, file, logger)
	}
	return getGeneralHandler(path, file, logger)
}

func getJSONHandler(path string, file wrapper.IFileOps, logger *logrus.Entry) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.GetJSON(path, file)
		if err != nil {
			logger.WithError(err).Error("error serving JSON from this endpoint")
			writeError(err, w)
			return
		}
		logger.Debug("successfully retrieved JSON; serving it")
		w.Write(bytes)
	}
}

func getGeneralHandler(path string, file wrapper.IFileOps, logger *logrus.Entry) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fileInfo, err := file.Open(path)
		defer fileInfo.Close()
		if err != nil {
			logger.WithError(err).Error("error opening file")
			writeError(err, w)
			return
		}

		bytes, err := file.ReadAll(fileInfo)
		if err != nil {
			logger.WithError(err).Error("error reading file")
			writeError(err, w)
			return
		}

		logger.Debug("successfully serving data")
		w.Write(bytes)
	}
}

func writeError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
