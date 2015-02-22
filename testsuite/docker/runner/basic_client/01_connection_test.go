package main

import (
	"os"
	"testing"

	"github.com/thoj/go-ircevent"
)

var (
	conn *irc.Connection
)

func TestMain(m *testing.M) {
	conn = irc.IRC("elemental_test", "user")
	defer conn.Quit()

	os.Exit(m.Run())
}

func TestBasicConnection(t *testing.T) {
	err := conn.Connect(os.Getenv("IRCD_PORT_6667_TCP_ADDR") + ":6667")
	if err != nil {
		t.Fatal(err)
	}
}
