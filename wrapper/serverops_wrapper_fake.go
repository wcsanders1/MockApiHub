package wrapper

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// FakeServerOps is a fake implementation of IServerOps
type FakeServerOps struct {
	mock.Mock
}

// Shutdown is a fake implementation of IServerOps.Shutdown()
func (ops *FakeServerOps) Shutdown(ctx context.Context) error {
	args := ops.Called(ctx)
	return args.Error(0)
}

// ListenAndServe is a fake implementation is IServerOps.ListenAndServe()
func (ops *FakeServerOps) ListenAndServe() error {
	args := ops.Called()
	return args.Error(0)
}

// ListenAndServeTLS is a fake implementation is IServerOps.ListenAndServeTLS()
func (ops *FakeServerOps) ListenAndServeTLS(certFile, keyFile string) error {
	args := ops.Called(certFile, keyFile)
	return args.Error(0)
}
