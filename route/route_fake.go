package route

import (
	"github.com/stretchr/testify/mock"
)

// FakeTree is a mockable route tree
type FakeTree struct {
	mock.Mock
}

// AddRoute is a mockable route.AddRoute()
func (tree *FakeTree) AddRoute(url string) (string, error) {
	args := tree.Called(url)
	return args.String(0), args.Error(1)
}

// GetRoute is a mockable route.GetRoute()
func (tree *FakeTree) GetRoute(url string) (string, map[string]string, error) {
	args := tree.Called(url)
	return args.String(0), args.Get(1).(map[string]string), args.Error(2)
}
