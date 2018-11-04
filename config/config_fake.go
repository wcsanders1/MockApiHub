package config

import (
	"os"

	"github.com/stretchr/testify/mock"
)

// FakeManager offers a fake, mockable implementation of IManager.
type FakeManager struct {
	mock.Mock
}

// GetAPIConfig is a fake implementation of GetAPIConfig().
func (mgr *FakeManager) GetAPIConfig(file os.FileInfo) (*APIConfig, error) {
	args := mgr.Called(file)
	return args.Get(0).(*APIConfig), args.Error(1)
}
