package main

import (
	"os"
	"sync"
	"testing"
	"time"

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
		wg.Done()
	})

	err := conn.Connect(os.Getenv("IRCD_PORT_6667_TCP_ADDR") + ":6667")
	if err != nil {
		t.Fatal(err)
	}

	go conn.Loop()
	wg.Wait()
}

func TestNotAbleToSaySomething(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	conn.AddCallback("404", func(e *irc.Event) {
		wg.Done()
	})

	conn.Privmsg("#foo", "bar")

	wg.Wait()
}

func TestJoinChannel(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	conn.AddCallback("ERROR", func(e *irc.Event) {
		t.Fatalf("%s", e.Message())
	})

	conn.AddCallback("353", func(e *irc.Event) {
		wg.Done()
	})

	conn.Join("#foo")

	wg.Wait()
}

func TestAbleToSaySomething(t *testing.T) {
	conn.AddCallback("404", func(e *irc.Event) {
		t.Fatal(e.Message())
	})

	conn.Privmsg("#foo", "bar")

	time.Sleep(2 * time.Second)
}
