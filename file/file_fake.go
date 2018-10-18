package file

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/mock"
)

// FakeBasicOps is a fake implementation of IBasicOps
type FakeBasicOps struct {
	mock.Mock
}

// ReadAll is a fake implementation of IBasicOps.ReadAll()
func (ops *FakeBasicOps) ReadAll(file *os.File) ([]byte, error) {
	args := ops.Called(file)
	return args.Get(0).([]byte), args.Error(1)
}

// Open is a fake implementation of IBasicOps.Open()
func (ops *FakeBasicOps) Open(file string) (*os.File, error) {
	args := ops.Called(file)
	return args.Get(0).(*os.File), args.Error(1)
}

// ReadDir is a fake implementation of IBasicOps.ReadDir()
func (ops *FakeBasicOps) ReadDir(dir string) ([]os.FileInfo, error) {
	args := ops.Called(dir)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

// DecodeFile is a fake implementation of IBasicOps.DecodeFile()
func (ops *FakeBasicOps) DecodeFile(file string, v interface{}) (toml.MetaData, error) {
	args := ops.Called(file, v)
	return args.Get(0).(toml.MetaData), args.Error(1)
}
