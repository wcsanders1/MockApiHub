package wrapper

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type (
	// IFileOps contains basic file operations, such as opening a file and reading it
	IFileOps interface {
		Open(string) (*os.File, error)
		ReadAll(*os.File) ([]byte, error)
		ReadDir(dir string) ([]os.FileInfo, error)
		DecodeFile(file string, v interface{}) (toml.MetaData, error)
	}

	// FileOps offers a real implementation of IFileOpc
	FileOps struct{}
)

// Open opens a file by its name and returns it
func (ops *FileOps) Open(file string) (*os.File, error) {
	return os.Open(file)
}

// ReadAll reads a file to its end
func (ops *FileOps) ReadAll(file *os.File) ([]byte, error) {
	return ioutil.ReadAll(file)
}

// ReadDir reads a directory and returns an array of FileInfo
func (ops *FileOps) ReadDir(dir string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dir)
}

// DecodeFile decodes a toml file
func (ops *FileOps) DecodeFile(file string, v interface{}) (toml.MetaData, error) {
	return toml.DecodeFile(file, v)
}
