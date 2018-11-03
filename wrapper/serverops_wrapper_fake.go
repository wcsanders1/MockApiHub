package wrapper

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type (
	call interface{}

	// FakeServerOps is a fake implementation of IServerOps
	FakeServerOps struct {
		mock.Mock
		Finished chan call
	}
)

// NewFakeServerOps returns a reference to a FakeServerOps
func NewFakeServerOps() *FakeServerOps {
	return &FakeServerOps{
		Finished: make(chan call),
	}
}

// Shutdown is a fake implementation of IServerOps.Shutdown()
func (ops *FakeServerOps) Shutdown(ctx context.Context) error {
	args := ops.Called(ctx)
	return args.Error(0)
}

// ListenAndServe is a fake implementation is IServerOps.ListenAndServe()
func (ops *FakeServerOps) ListenAndServe() error {
	args := ops.Called()
	close(ops.Finished)
	return args.Error(0)
}

// ListenAndServeTLS is a fake implementation is IServerOps.ListenAndServeTLS()
func (ops *FakeServerOps) ListenAndServeTLS(certFile, keyFile string) error {
	args := ops.Called(certFile, keyFile)
	close(ops.Finished)
	return args.Error(0)
}

// WaitForListenAndServe ensures that the ListenAndServe function is called before a test makes assertions about that call
func (ops *FakeServerOps) WaitForListenAndServe() {
	<-ops.Finished
}
