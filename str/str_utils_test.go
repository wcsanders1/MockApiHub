package str

import "testing"

func TestGetPort(t *testing.T) {
	port := GetPort(5000)
	if port != ":5000" {
		t.Errorf("Port incorrect, got: %s, expected: %s", port, ":5000")
	}
}