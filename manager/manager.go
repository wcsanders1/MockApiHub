package manager

import (
	"MockApiHub/api"
)

// Manager coordinates and controls the apis
type Manager struct{
	apis map[string]*api.API
}

// NewManager returns an instance of the Manager type
func NewManager() *Manager {
	return &Manager{}
}

func (mgr *Manager) StartMockApiHub() {
	mgr.apis = make(map[string]*api.API)
}