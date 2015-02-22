package runner

import (
	"os"
	"testing"
)

func TestEnvironment(t *testing.T) {
	if os.Getenv("IRCD_PORT_6667_TCP_ADDR") == "" {
		t.Fatal("Container not linked.")
	}
}
