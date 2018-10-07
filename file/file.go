package file

import (
	"io/ioutil"
	"os"
)

type (
	// IBasicOps contains basic file operations, such as opening a file and reading it
	IBasicOps interface {
		Open(string) (*os.File, error)
		ReadAll(*os.File) ([]byte, error)
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
