package file

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type (
	// IBasicOps contains basic file operations, such as opening a file and reading it
	IBasicOps interface {
		Open(string) (*os.File, error)
		ReadAll(*os.File) ([]byte, error)
		ReadDir(dir string) ([]os.FileInfo, error)
		DecodeFile(file string, v interface{}) (toml.MetaData, error)
	}

	// BasicOps offers a real implementation of IBasicOpc
	BasicOps struct{}
)

// Open opens a file by its name and returns it
func (ops *BasicOps) Open(file string) (*os.File, error) {
	return os.Open(file)
}

// ReadAll reads a file to its end
func (ops *BasicOps) ReadAll(file *os.File) ([]byte, error) {
	return ioutil.ReadAll(file)
}

// ReadDir reads a directory and returns an array of FileInfo
func (ops *BasicOps) ReadDir(dir string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dir)
}

// DecodeFile decodes a toml file
func (ops *BasicOps) DecodeFile(file string, v interface{}) (toml.MetaData, error) {
	return toml.DecodeFile(file, v)
}
