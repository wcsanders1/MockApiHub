package str

import (
	"fmt"
)

// GetPort returns the port in the format that http server expects
func GetPort(port int) string {
	return fmt.Sprintf(":%d", port)
}