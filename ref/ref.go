//Package ref provides functionality allowing a function to get information about itself.
package ref

import (
	"runtime"
)

// GetFuncName returns the name of the calling function.
func GetFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
