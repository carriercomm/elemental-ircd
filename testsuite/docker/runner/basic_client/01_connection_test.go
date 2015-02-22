package main

import (
	"os"
	"sync"
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
	wg := sync.WaitGroup{}
	wg.Add(1)

	conn.AddCallback("ERROR", func(e *irc.Event) {
		t.Fatalf("%s", e.Message())
	})

	conn.AddCallback("001", func(e *irc.Event) {
		t.Log(e.Message())
		wg.Done()
	})

	err := conn.Connect(os.Getenv("IRCD_PORT_6667_TCP_ADDR") + ":6667")
	if err != nil {
		t.Fatal(err)
	}

	go conn.Loop()
	wg.Wait()
}
