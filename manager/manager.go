package manager

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"MockApiHub/api"
	"MockApiHub/config"
	"MockApiHub/log"
	"MockApiHub/reflect"
	"MockApiHub/str"

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
}

const (
	apiDir    = "./api/apis"
	apiDirExt = "Api"
)

// NewManager returns an instance of the Manager type
func NewManager(config *config.AppConfig) (*Manager, error) {
	mgr := &Manager{}
	logger := log.NewLogger(&config.Log)
	mgr.log = logger.WithField("pkg", "manager")

	server, err := createManagerServer(config.HTTP.Port, mgr)
	if err != nil {
		mgr.log.Panic(err)
		return nil, err
	}

	mgr.config = config
	mgr.server = server
	mgr.apis = make(map[string]*api.API)

	return mgr, nil
}

func (mgr *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)
	path := str.CleanURL(r.URL.String())
	contextLogger := mgr.log.WithFields(logrus.Fields{
		"method": method,
		"path":   path,
		"func":   reflect.GetFuncName(),
	})

	if len(method) == 0 || len(path) == 0 {
		contextLogger.Warn("endpoint not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("endpoint not found"))
		return
	}

	if handler, exists := mgr.hubAPIHandlers[method][path]; exists {
		contextLogger.Debug("endpoint exists")
		handler(w, r)
		return
	}
	contextLogger.Warn("endpoint not found")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("endpoint not found"))
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

// StartMockAPIHub registers the mock apis and serves them
func (mgr *Manager) StartMockAPIHub() error {
	err := mgr.loadMockAPIs()
	if err != nil {
		return err
	}

	mgr.registerHubAPIHandlers()
	mgr.registerMockAPIs()
	mgr.startHubServer()

	return nil
}

func (mgr *Manager) startHubServer() error {
	if mgr.config.HTTP.UseTLS {
		if err := mgr.startHubServerUsingTLS(); err != nil {
			panic(err)
		}
	} else {
		mgr.server.ListenAndServe()
	}

	return nil
}

func (mgr *Manager) startHubServerUsingTLS() error {
	certFile := mgr.config.HTTP.CertFile
	keyFile := mgr.config.HTTP.KeyFile

	if _, err := os.Stat(certFile); err != nil {
		return fmt.Errorf(fmt.Sprintf("%s cert file does not exist", certFile))
	}

	if _, err := os.Stat(keyFile); err != nil {
		return fmt.Errorf(fmt.Sprintf("%s key file does not exist", keyFile))
	}

	return mgr.server.ListenAndServeTLS(certFile, keyFile)
}

func (mgr *Manager) registerMockAPIs() {
	for dir, api := range mgr.apis {
		err := api.Register(dir, mgr.config.HTTP.CertFile, mgr.config.HTTP.KeyFile)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (mgr *Manager) loadMockAPIs() error {
	files, err := ioutil.ReadDir(apiDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		apiConfig, err := getAPIConfig(file)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if mgr.apiByPortExists(apiConfig.HTTP.Port) {
			fmt.Println(fmt.Sprintf("Trying to register %s api on port %d, but there is already an "+
				"api registered on that port. Skipping.", file.Name(), apiConfig.HTTP.Port))

			continue
		}

		api, err := api.NewAPI(apiConfig)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if api != nil {
			mgr.apis[file.Name()] = api
		}
	}
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

func getAPIConfig(file os.FileInfo) (*config.APIConfig, error) {
	if !file.IsDir() || !isAPI(file.Name()) {
		return nil, nil
	}

	dir := file.Name()
	fmt.Println("Found the following mock api: ", dir)
	apiConfig, err := getAPIConfigFromDir(dir)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return apiConfig, nil
}

func getAPIConfigFromDir(dir string) (*config.APIConfig, error) {
	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/%s", apiDir, dir))
	for _, file := range files {
		if isAPIConfig(file.Name()) {
			apiConfig, err := decodeAPIConfig(dir, file.Name())
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			return apiConfig, nil
		}
	}
	return nil, nil
}

func decodeAPIConfig(dir string, fileName string) (*config.APIConfig, error) {
	path := fmt.Sprintf("%s/%s/%s", apiDir, dir, fileName)
	var config config.APIConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		fmt.Println(err)
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
