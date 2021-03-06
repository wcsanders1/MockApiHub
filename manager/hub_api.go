package manager

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wcsanders1/MockApiHub/api"
	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/log"
	"github.com/wcsanders1/MockApiHub/ref"
)

type apiDisplay struct {
	BaseURL   string
	Endpoints map[string]config.Endpoint
	Port      int
}

const (
	refreshAPIsPath = "refresh-all-mock-apis"
	showAllAPIsPath = "show-all-registered-mock-apis"
)

func (mgr *Manager) refreshMockAPIs(w http.ResponseWriter, r *http.Request) {
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())
	contextLogger.Debug("refreshing all mock APIs")

	mgr.shutDownMockAPIs()
	mgr.apis = make(map[string]api.IAPI)
	if err := mgr.loadMockAPIs(); err != nil {
		contextLogger.WithError(err).Error("error loading mock APIs")
		return
	}

	mgr.startMockAPIs()
	msg := "successfully refreshed mock apis"
	w.Write([]byte(msg))
	contextLogger.Debug(msg)
}

func (mgr *Manager) showRegisteredMockAPIs(w http.ResponseWriter, r *http.Request) {
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())
	contextLogger.Debug("showing all registered mock APIs")

	apis := make(map[string]apiDisplay)
	for apiName, api := range mgr.apis {
		apis[apiName] = apiDisplay{
			BaseURL:   api.GetBaseURL(),
			Port:      api.GetPort(),
			Endpoints: api.GetEndpoints(),
		}
	}

	apisJSON, err := json.Marshal(apis)
	if err != nil {
		contextLogger.WithError(err).Error("error displaying mock APIs")
		return
	}

	w.Write(apisJSON)
	contextLogger.WithField("registeredAPIs", apis).Debug("successfully showed all registered mock APIs")
}

func (mgr *Manager) registerHubAPIHandlers() {
	contextLogger := mgr.log.WithField(log.FuncField, ref.GetFuncName())
	contextLogger.Debug("registering hub API handlers")

	mgr.hubAPIHandlers = make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	mgr.hubAPIHandlers[http.MethodPost] = make(map[string]func(http.ResponseWriter, *http.Request))
	mgr.hubAPIHandlers[http.MethodGet] = make(map[string]func(http.ResponseWriter, *http.Request))

	mgr.hubAPIHandlers[http.MethodPost][strings.ToLower(refreshAPIsPath)] = mgr.refreshMockAPIs
	mgr.hubAPIHandlers[http.MethodGet][strings.ToLower(showAllAPIsPath)] = mgr.showRegisteredMockAPIs

	contextLogger.Debug("successfully registered hub API handlers")
}
