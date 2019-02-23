package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/constants"
	"github.com/wcsanders1/MockApiHub/json"
	"github.com/wcsanders1/MockApiHub/log"
	"github.com/wcsanders1/MockApiHub/ref"
	"github.com/wcsanders1/MockApiHub/wrapper"

	"github.com/sirupsen/logrus"
)

type (
	iCreator interface {
		getHandler(enforceValidJSON bool, headers []config.Header, dir, fileName string, file wrapper.IFileOps) func(w http.ResponseWriter, r *http.Request)
		startAPI(defaultCert, defaultKey string, server wrapper.IServerOps, httpConfig config.HTTP) error
	}

	creator struct {
		log *logrus.Entry
	}
)

func newCreator(logger *logrus.Entry) *creator {
	return &creator{
		log: logger,
	}
}

func (c creator) getHandler(enforceValidJSON bool, headers []config.Header, dir, fileName string, file wrapper.IFileOps) func(w http.ResponseWriter, r *http.Request) {
	var path string
	if len(fileName) > 0 {
		path = fmt.Sprintf("%s/%s/%s", constants.APIDir, dir, fileName)
	} else {
		path = ""
	}

	contextLogger := c.log.WithFields(logrus.Fields{
		log.FuncField:      "handler for mock API",
		"enforceValidJSON": enforceValidJSON,
		log.PathField:      path,
	})

	if enforceValidJSON {
		return getJSONHandler(path, headers, file, contextLogger)
	}
	return getGeneralHandler(path, headers, file, contextLogger)
}

func getJSONHandler(path string, headers []config.Header, file wrapper.IFileOps, logger *logrus.Entry) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		for _, header := range headers {
			w.Header().Set(header.Key, header.Value)
		}

		if len(path) == 0 {
			return
		}

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

func getGeneralHandler(path string, headers []config.Header, file wrapper.IFileOps, logger *logrus.Entry) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, header := range headers {
			w.Header().Set(header.Key, header.Value)
		}

		if len(path) == 0 {
			return
		}

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

func (c creator) startAPI(defaultCert, defaultKey string, server wrapper.IServerOps, httpConfig config.HTTP) error {
	contextLogger := c.log.WithFields(logrus.Fields{
		log.FuncField:            ref.GetFuncName(),
		log.DefaultCertFileField: defaultCert,
		log.DefaultKeyFileField:  defaultKey,
	})

	if httpConfig.UseTLS {
		cert, key, err := getCertAndKeyFileNames(defaultCert, defaultKey, httpConfig)
		if err != nil {
			contextLogger.WithError(err).Error("error getting TLS cert and key")
			return err
		}

		go func() {
			if err := server.ListenAndServeTLS(cert, key); err != nil {
				contextLogger.WithError(err).Error("error starting mock API with TLS")
			}
		}()
		return nil
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			contextLogger.WithError(err).Error("mock API server error")
		}
	}()
	return nil
}

func getCertAndKeyFileNames(defaultCert, defaultKey string, httpConfig config.HTTP) (string, string, error) {
	if len(httpConfig.CertFile) > 0 && len(httpConfig.KeyFile) > 0 {
		return httpConfig.CertFile, httpConfig.KeyFile, nil
	}

	if len(httpConfig.CertFile) == 0 && len(httpConfig.KeyFile) > 0 {
		return "", "", errors.New("key provided without cert")
	}

	if len(httpConfig.KeyFile) == 0 && len(httpConfig.CertFile) > 0 {
		return "", "", errors.New("cert provided without key")
	}

	return defaultCert, defaultKey, nil
}
