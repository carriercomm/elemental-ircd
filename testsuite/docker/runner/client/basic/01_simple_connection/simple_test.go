/*
This test package tests basic CRUD-style operations in IRC networks.
This is purely focused on joining and messaging channels.
*/
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

// TestBasicConnection ensures that the test client can make a basic connection
// to the ircd over port 6667.
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

// TestNotAbleToSaySomething ensures that #foo (which either will not exist
// or must have mode +n set) cannot be externally messaged.
func TestNotAbleToSaySomething(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	conn.AddCallback("404", func(e *irc.Event) {
		wg.Done()
	})

	conn.Privmsg("#foo", "bar")

	wg.Wait()
}

// TestJoinChannel ensures that the test client can join a new or existing channel
// without error due to being banned.
func TestJoinChannel(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	conn.AddCallback("474", func(e *irc.Event) {
		t.Fatalf("Banned from #foo? %s", e.Message())
	})

	conn.AddCallback("353", func(e *irc.Event) {
		wg.Done()
	})

	conn.Join("#foo")

	wg.Wait()
}

// TestAbleToSaySomething ensures that the test client can say something to #foo without
// being denied speaking.
func TestAbleToSaySomething(t *testing.T) {
	conn.AddCallback("404", func(e *irc.Event) {
		t.Fatal(e.Message())
	})

	conn.Privmsg("#foo", "bar")

	time.Sleep(2 * time.Second)
}
