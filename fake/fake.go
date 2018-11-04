package fake

import (
	"net/http"
	"os"
	"time"

	"github.com/stretchr/testify/mock"
)

type (
	// FileInfo offers a fake, mockable implementation of os.FileInfo.
	FileInfo struct {
		mock.Mock
	}

	// ResponseWriter offers a fake, mockable implementation of http.ResponseWriter.
	ResponseWriter struct {
		mock.Mock
	}
)

// Name is a fake implementation of os.FileInfo.Name().
func (fi *FileInfo) Name() string {
	args := fi.Called()
	return args.String(0)
}

// Size is a fake implementation of os.FileInfo.Size().
func (fi *FileInfo) Size() int64 {
	args := fi.Called()
	return args.Get(0).(int64)
}

// Mode is a fake implementation of os.FileInfo.Mode().
func (fi *FileInfo) Mode() os.FileMode {
	args := fi.Called()
	return args.Get(0).(os.FileMode)
}

// ModTime is a fake implementation of os.FileInfo.ModTime().
func (fi *FileInfo) ModTime() time.Time {
	args := fi.Called()
	return args.Get(0).(time.Time)
}

// IsDir is a fake implementation of os.FileInfo.IsDir().
func (fi *FileInfo) IsDir() bool {
	args := fi.Called()
	return args.Bool(0)
}

// Sys is a fake implementation of os.FileInfo.Sys().
func (fi *FileInfo) Sys() interface{} {
	args := fi.Called()
	return args.Get(0).(interface{})
}

// Header is a fake implementation of http.Header().
func (r *ResponseWriter) Header() http.Header {
	args := r.Called()
	return args.Get(0).(http.Header)
}

// Write is a fake implementation of http.Write().
func (r *ResponseWriter) Write(data []byte) (int, error) {
	args := r.Called(data)
	return args.Int(0), args.Error(1)
}

// WriteHeader is a fake implementation of http.WriteHeader().
func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.Called(statusCode)
	return
}
