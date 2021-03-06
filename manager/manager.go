//Package manager manages mock APIs.
package manager

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/wcsanders1/MockApiHub/api"
	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/constants"
	"github.com/wcsanders1/MockApiHub/log"
	"github.com/wcsanders1/MockApiHub/ref"
	"github.com/wcsanders1/MockApiHub/str"
	"github.com/wcsanders1/MockApiHub/wrapper"

	"github.com/sirupsen/logrus"
)

// Manager coordinates and controls the mock APIs
type Manager struct {
	apis           map[string]api.IAPI
	config         *config.AppConfig
	server         wrapper.IServerOps
	hubAPIHandlers map[string]map[string]func(http.ResponseWriter, *http.Request)
	log            *logrus.Entry
	file           wrapper.IFileOps
	configManager  config.IManager
}

// NewManager returns an instance of the Manager type.
func NewManager(appConfig *config.AppConfig) (*Manager, error) {
	mgr := &Manager{}
	mgr.log = log.NewLogger(&appConfig.Log, "manager").WithFields(logrus.Fields{
		log.PortField:     appConfig.HTTP.Port,
		log.UseTLSField:   appConfig.HTTP.UseTLS,
		log.CertFileField: appConfig.HTTP.CertFile,
		log.KeyFileField:  appConfig.HTTP.KeyFile,
	})
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())
	contextLogger.Info("creating new manager")

	server, err := createManagerServer(appConfig.HTTP.Port, mgr)
	if err != nil {
		contextLogger.WithError(err).Error("error creating manager")
		return nil, err
	}

	mgr.config = appConfig
	mgr.server = wrapper.NewServerOps(server)
	mgr.apis = make(map[string]api.IAPI)
	mgr.file = &wrapper.FileOps{}
	mgr.configManager = config.NewConfigManager()
	contextLogger.Info("successfully created new manager")
	return mgr, nil
}

// StartMockAPIHub registers the mock apis and serves them.
func (mgr *Manager) StartMockAPIHub() error {
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())
	contextLogger.Debug("starting mock API hub")

	if err := mgr.loadMockAPIs(); err != nil {
		contextLogger.WithError(err).Error("error loading mock APIs")
		return err
	}
	mgr.startMockAPIs()
	contextLogger.Debug("successfully started mock APIs; will next start the mock API hub")

	mgr.registerHubAPIHandlers()
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

	mgr.shutDownMockAPIs()

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

func (mgr *Manager) shutDownMockAPIs() {
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

	if _, err := mgr.file.Stat(certFile); err != nil {
		contextLogger.WithError(err).Error("error starting hub server using TLS -- cert file does not exist")
		return err
	}

	if _, err := mgr.file.Stat(keyFile); err != nil {
		contextLogger.WithError(err).Error("error starting hub server using TLS -- key file does not exist")
		return err
	}

	if err := mgr.server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
		contextLogger.WithError(err).Error("error starting hub server using TLS")
		return err
	}

	return nil
}

func (mgr *Manager) startMockAPIs() {
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())

	for dir, api := range mgr.apis {
		contextLoggerAPI := contextLogger.WithFields(logrus.Fields{
			log.BaseURLField: api.GetBaseURL(),
			log.PortField:    api.GetPort(),
		})
		contextLoggerAPI.Debug("starting mock API")
		if err := api.Start(dir, mgr.config.HTTP.CertFile, mgr.config.HTTP.KeyFile); err != nil {
			contextLoggerAPI.WithError(err).Error("error starting mock API -- moving on to next mock API")
		}
	}
}

func (mgr *Manager) loadMockAPIs() error {
	contextLogger := mgr.log.WithFields(logrus.Fields{
		log.FuncField: ref.GetFuncName(),
		"apiDir":      constants.APIDir,
	})
	contextLogger.Debug("loading mock APIs")

	files, err := mgr.file.ReadDir(constants.APIDir)
	if err != nil {
		contextLogger.WithError(err).Error("error reading API directory")
		return err
	}

	for _, file := range files {
		contextLoggerFile := contextLogger.WithField("file", file.Name())

		apiConfig, err := mgr.configManager.GetAPIConfig(file)
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
