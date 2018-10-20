package wrapper

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/mock"
)

// FakeFileOps is a fake implementation of IFileOps
type FakeFileOps struct {
	mock.Mock
}

// ReadAll is a fake implementation of IFileOps.ReadAll()
func (ops *FakeFileOps) ReadAll(file *os.File) ([]byte, error) {
	args := ops.Called(file)
	return args.Get(0).([]byte), args.Error(1)
}

// Open is a fake implementation of IFileOps.Open()
func (ops *FakeFileOps) Open(file string) (*os.File, error) {
	args := ops.Called(file)
	return args.Get(0).(*os.File), args.Error(1)
}

// ReadDir is a fake implementation of IFileOps.ReadDir()
func (ops *FakeFileOps) ReadDir(dir string) ([]os.FileInfo, error) {
	args := ops.Called(dir)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

// DecodeFile is a fake implementation of IFileOps.DecodeFile()
func (ops *FakeFileOps) DecodeFile(file string, v interface{}) (toml.MetaData, error) {
	args := ops.Called(file, v)
	return args.Get(0).(toml.MetaData), args.Error(1)
}

// Stat is a fake implementation of IFileOps.Stat()
func (ops *FakeFileOps) Stat(file string) (os.FileInfo, error) {
	args := ops.Called(file)
	return args.Get(0).(os.FileInfo), args.Error(1)
}
