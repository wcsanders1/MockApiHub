//Package wrapper wraps functionality concerning server and file interaction so that the functionality can be mocked.
package wrapper

import (
	"context"
	"net/http"
)

type (
	// IServerOps contains basic server operations
	IServerOps interface {
		Shutdown(context.Context) error
		ListenAndServe() error
		ListenAndServeTLS(string, string) error
	}

	// ServerOps provides a real implementation of IServerOps
	ServerOps struct {
		server *http.Server
	}
)

// NewServerOps returns a pointer to a new ServerOps
func NewServerOps(server *http.Server) *ServerOps {
	return &ServerOps{server}
}

// Shutdown shuts down the server
func (ops *ServerOps) Shutdown(ctx context.Context) error {
	return ops.server.Shutdown(ctx)
}

// ListenAndServe starts the server
func (ops *ServerOps) ListenAndServe() error {
	return ops.server.ListenAndServe()
}

// ListenAndServeTLS starts the server using TLS
func (ops *ServerOps) ListenAndServeTLS(certFile, keyFile string) error {
	return ops.server.ListenAndServeTLS(certFile, keyFile)
}
