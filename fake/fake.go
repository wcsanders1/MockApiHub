package fake

import (
	"os"
	"time"

	"github.com/stretchr/testify/mock"
)

type (
	// BasicOps offers a fake, mockable implementation of file.IBasicOps
	BasicOps struct {
		mock.Mock
	}

	// FileInfo offers a fake, mockable implementation of os.FileInfo
	FileInfo struct {
		mock.Mock
	}
)

// ReadAll is a fake implementation of IBasicOps.ReadAll()
func (ops *BasicOps) ReadAll(file *os.File) ([]byte, error) {
	args := ops.Called(file)
	return args.Get(0).([]byte), args.Error(1)
}

// Open is a fake implementation of IBasicOps.Open()
func (ops *BasicOps) Open(file string) (*os.File, error) {
	args := ops.Called(file)
	return args.Get(0).(*os.File), args.Error(1)
}

// ReadDir is a fake implementation of IBasicOps.ReadDir()
func (ops *BasicOps) ReadDir(dir string) ([]os.FileInfo, error) {
	args := ops.Called(dir)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

// Name is a fake implementation of os.FileInfo.Name()
func (fi *FileInfo) Name() string {
	args := fi.Called()
	return args.String(0)
}

// Size is a fake implementation of os.FileInfo.Size()
func (fi *FileInfo) Size() int64 {
	args := fi.Called()
	return args.Get(0).(int64)
}

// Mode is a fake implementation of os.FileInfo.Mode()
func (fi *FileInfo) Mode() os.FileMode {
	args := fi.Called()
	return args.Get(0).(os.FileMode)
}

// ModTime is a fake implementation of os.FileInfo.ModTime()
func (fi *FileInfo) ModTime() time.Time {
	args := fi.Called()
	return args.Get(0).(time.Time)
}

// IsDir is a fake implementation of os.FileInfo.IsDir()
func (fi *FileInfo) IsDir() bool {
	args := fi.Called()
	return args.Bool(0)
}

// Sys is a fake implementation of os.FileInfo.Sys()
func (fi *FileInfo) Sys() interface{} {
	args := fi.Called()
	return args.Get(0).(interface{})
}
