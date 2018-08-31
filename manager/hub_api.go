package manager

import (
	"net/http"
)

const (
	refreshAPISPath = "refreshAllAPIs"
)

// RefreshAPI clears all mock apis and reloads them
func (mgr *Manager) RefreshAPI (w http.ResponseWriter, r *http.Request) {
	
	w.Write([]byte("refreshing (but not really)"))
}

// RegisterHubAPIHandlers registers the hub API handlers
func (mgr *Manager) RegisterHubAPIHandlers() error {

	mgr.hubAPIHandlers = make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	mgr.hubAPIHandlers[http.MethodPost] = make(map[string]func(http.ResponseWriter, *http.Request))

	mgr.hubAPIHandlers[http.MethodPost][refreshAPISPath] = mgr.RefreshAPI

	return nil
}