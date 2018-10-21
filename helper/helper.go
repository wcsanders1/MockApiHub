package helper

import (
	"os"

	"github.com/wcsanders1/MockApiHub/fake"
)

// GetFakeFileInfoAndCollection returns fake os.FileInfo and a collection of os.FileInfo
// with that os.FileInfo in it, for testing.
func GetFakeFileInfoAndCollection(dir, file string) (os.FileInfo, []os.FileInfo) {
	fileInfo := new(fake.FileInfo)
	fileInfo.On("Name").Return(dir)
	fileInfo.On("IsDir").Return(true)

	fileInfoCollection := []os.FileInfo{}
	fileInfoInner := new(fake.FileInfo)
	fileInfoInner.On("Name").Return(file)
	fileInfoCollection = append(fileInfoCollection, fileInfoInner)

	return fileInfo, fileInfoCollection
}
