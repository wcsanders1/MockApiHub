package fake

import (
	"os"

	"github.com/stretchr/testify/mock"
)

// BasicOps offers a fake, mockable implementation of IBasicOps
type BasicOps struct {
	mock.Mock
}

// ReadAll is a fake implementation of BasicOps.ReadAll
func (ops *BasicOps) ReadAll(file *os.File) ([]byte, error) {
	args := ops.Called(file)
	return args.Get(0).([]byte), args.Error(1)
}

// Open is a fake implementation of BasicOps.Open
func (ops *BasicOps) Open(file string) (*os.File, error) {
	args := ops.Called(file)
	return args.Get(0).(*os.File), args.Error(1)
}

// ReadDir is a fake implementation of BasicOps.ReadDir
func (ops *BasicOps) ReadDir(dir string) ([]os.FileInfo, error) {
	args := ops.Called(dir)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}
