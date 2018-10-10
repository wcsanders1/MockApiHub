package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wcsanders1/MockApiHub/api"
	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/file"
	"github.com/wcsanders1/MockApiHub/log"
	"github.com/wcsanders1/MockApiHub/ref"
	"github.com/wcsanders1/MockApiHub/str"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

// Manager coordinates and controls the apis
type Manager struct {
	apis           map[string]*api.API
	config         *config.AppConfig
	server         *http.Server
	hubAPIHandlers map[string]map[string]func(http.ResponseWriter, *http.Request)
	log            *logrus.Entry
	file           file.IBasicOps
}

const (
	apiDir    = "./api/apis"
	apiDirExt = "Api"
)

// NewManager returns an instance of the Manager type
func NewManager(config *config.AppConfig) (*Manager, error) {
	mgr := &Manager{}
	mgr.log = log.NewLogger(&config.Log, "manager").WithFields(logrus.Fields{
		log.PortField:     config.HTTP.Port,
		log.UseTLSField:   config.HTTP.UseTLS,
		log.CertFileField: config.HTTP.CertFile,
		log.KeyFileField:  config.HTTP.KeyFile,
	})
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())
	contextLogger.Info("creating new manager")

	server, err := createManagerServer(config.HTTP.Port, mgr)
	if err != nil {
		contextLogger.WithError(err).Error("error creating manager")
		return nil, err
	}

	mgr.config = config
	mgr.server = server
	mgr.apis = make(map[string]*api.API)
	mgr.file = &file.BasicOps{}
	contextLogger.Info("successfully created new manager")
	return mgr, nil
}

// StartMockAPIHub registers the mock apis and serves them
func (mgr *Manager) StartMockAPIHub() error {
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())
	contextLogger.Debug("starting mock API hub")

	if err := mgr.loadMockAPIs(); err != nil {
		contextLogger.WithError(err).Error("error loading mock APIs")
		return err
	}

	mgr.registerHubAPIHandlers()
	mgr.registerMockAPIs()

	if err := mgr.startHubServer(); err != nil {
		contextLogger.WithError(err).Error("error starting hub server")
		return err
	}

	contextLogger.Debug("successfully started mock API hub")
	return nil
}

// StopMockAPIHub shuts down all mock API servers and the hub server,
// and panics on error.
func (mgr *Manager) StopMockAPIHub() {
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())
	contextLogger.Debug("stopping mock API hub")

	mgr.shutdownMockAPIs()

	if err := mgr.shutdownHubServer(); err != nil {
		contextLogger.WithError(err).Panic("error shutting down hub server")
		return
	}

	contextLogger.Debug("successfully stopped mock API hub")
}

func (mgr *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)
	path := str.CleanURL(r.URL.String())
	contextLogger := mgr.log.WithFields(logrus.Fields{
		log.MethodField: method,
		log.PathField:   path,
		log.FuncField:   ref.GetFuncName(),
	})

	if len(method) == 0 || len(path) == 0 {
		contextLogger.Warn("endpoint not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("endpoint not found"))
		return
	}

	if handler, exists := mgr.hubAPIHandlers[method][path]; exists {
		contextLogger.Debug("endpoint hit")
		handler(w, r)
		return
	}

	contextLogger.Warn("endpoint not found")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("endpoint not found"))
}

func (mgr *Manager) shutdownMockAPIs() {
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())
	contextLogger.Debug("shutting down mock APIs")

	for apiName, api := range mgr.apis {
		contextLoggerAPI := contextLogger.WithFields(logrus.Fields{
			log.PortField:    api.GetPort(),
			log.BaseURLField: api.GetBaseURL(),
			log.APINameField: apiName,
		})
		contextLoggerAPI.Info("shutting down mock API")
		if err := api.Shutdown(); err != nil {
			contextLoggerAPI.WithError(err).Error("error shutting down mock API")
			continue
		}
		contextLoggerAPI.Info("successfully shut down mock API")
	}

	contextLogger.Debug("finished shutting down mock APIs")
}

func (mgr *Manager) shutdownHubServer() error {
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())
	contextLogger.Debug("shutting down hub server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := mgr.server.Shutdown(ctx); err != nil {
		contextLogger.WithError(err).Error("error shutting down hub server")
		return err
	}

	contextLogger.Info("successfully shut down hub server")
	return nil
}

func (mgr *Manager) startHubServer() error {
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())

	if mgr.config.HTTP.UseTLS {
		contextLogger.Debug("starting hub server using TLS")
		if err := mgr.startHubServerUsingTLS(); err != nil {
			contextLogger.WithError(err).Error("hub server error; using TLS")
			return err
		}
	} else {
		contextLogger.Debug("starting hub server not using TLS")
		if err := mgr.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			contextLogger.WithError(err).Error("hub server error; not using TLS")
			return err
		}
	}

	return nil
}

func (mgr *Manager) startHubServerUsingTLS() error {
	certFile := mgr.config.HTTP.CertFile
	keyFile := mgr.config.HTTP.KeyFile
	contextLogger := mgr.log.WithFields(logrus.Fields{
		log.CertFileField: certFile,
		log.KeyFileField:  keyFile,
		log.FuncField:     ref.GetFuncName(),
	})
	contextLogger.Debug("starting hub server using TLS")

	if _, err := os.Stat(certFile); err != nil {
		contextLogger.WithError(err).Error("error starting hub server using TLS -- cert file does not exist")
		return err
	}

	if _, err := os.Stat(keyFile); err != nil {
		contextLogger.WithError(err).Error("error starting hub server using TLS -- key file does not exist")
		return err
	}

	if err := mgr.server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
		contextLogger.WithError(err).Error("error starting hub server using TLS")
		return err
	}

	return nil
}

func (mgr *Manager) registerMockAPIs() {
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())

	for dir, api := range mgr.apis {
		contextLoggerAPI := contextLogger.WithFields(logrus.Fields{
			log.BaseURLField: api.GetBaseURL(),
			log.PortField:    api.GetPort(),
		})
		contextLoggerAPI.Debug("registering mock API")
		if err := api.Register(dir, mgr.config.HTTP.CertFile, mgr.config.HTTP.KeyFile); err != nil {
			contextLoggerAPI.WithError(err).Error("error registering mock API -- moving on to next mock API")
			continue
		}
	}
}

func (mgr *Manager) loadMockAPIs() error {
	contextLogger := mgr.log.WithFields(logrus.Fields{
		log.FuncField: ref.GetFuncName(),
		"apiDir":      apiDir,
	})
	contextLogger.Debug("loading mock APIs")

	files, err := mgr.file.ReadDir(apiDir)
	if err != nil {
		contextLogger.WithError(err).Error("error reading API directory")
		return err
	}

	for _, file := range files {
		contextLoggerFile := contextLogger.WithField("file", file.Name())

		apiConfig, err := mgr.getAPIConfig(file)
		if err != nil {
			contextLoggerFile.WithError(err).Error("error getting API config from file -- moving on to next mock API file")
			continue
		}

		contextLoggerFileAPI := contextLoggerFile.WithFields(logrus.Fields{
			log.BaseURLField:  apiConfig.BaseURL,
			log.UseTLSField:   apiConfig.HTTP.UseTLS,
			log.CertFileField: apiConfig.HTTP.CertFile,
			log.KeyFileField:  apiConfig.HTTP.KeyFile,
			log.PortField:     apiConfig.HTTP.Port,
		})

		if mgr.apiByPortExists(apiConfig.HTTP.Port) {
			contextLoggerFileAPI.Warn("a mock API is already loaded on this port -- moving on to next mock API")
			continue
		}

		api, err := api.NewAPI(apiConfig)
		if err != nil {
			contextLoggerFileAPI.WithError(err).Error("error loading mock API -- moving on to next mock API")
			continue
		}

		if api != nil {
			contextLoggerFileAPI.Info("successfully loaded mock API")
			mgr.apis[file.Name()] = api
			continue
		}

		contextLoggerFileAPI.Warn("unable to load mock API")
	}
	contextLogger.Debug("finished loading mock APIs")
	return nil
}

func (mgr *Manager) apiByPortExists(port int) bool {
	for _, api := range mgr.apis {
		if api.GetPort() == port {
			return true
		}
	}
	return false
}

func (mgr *Manager) getAPIConfig(file os.FileInfo) (*config.APIConfig, error) {
	if !file.IsDir() || !isAPI(file.Name()) {
		return nil, errors.New("not a mock API directory")
	}

	dir := file.Name()
	apiConfig, err := mgr.getAPIConfigFromDir(dir)
	if err != nil {
		return nil, err
	}
	return apiConfig, nil
}

func (mgr *Manager) getAPIConfigFromDir(dir string) (*config.APIConfig, error) {
	files, _ := mgr.file.ReadDir(fmt.Sprintf("%s/%s", apiDir, dir))
	for _, file := range files {
		if isAPIConfig(file.Name()) {
			apiConfig, err := decodeAPIConfig(dir, file.Name())
			if err != nil {
				return nil, err
			}
			return apiConfig, nil
		}
	}
	return nil, nil
}

func createManagerServer(port int, mgr *Manager) (*http.Server, error) {
	if port == 0 {
		return nil, errors.New("no port provided")
	}

	server := &http.Server{
		Addr:    str.GetPort(port),
		Handler: mgr,
	}
	return server, nil
}

func decodeAPIConfig(dir string, fileName string) (*config.APIConfig, error) {
	path := fmt.Sprintf("%s/%s/%s", apiDir, dir, fileName)
	var config config.APIConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func isAPI(dir string) bool {
	return len(dir) > 3 && dir[len(dir)-3:] == apiDirExt
}

func isAPIConfig(fileName string) bool {
	ext := filepath.Ext(fileName)
	return ext == ".toml"
}
