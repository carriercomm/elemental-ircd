/*
This test package tests opering up.
*/
package operup

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/thoj/go-ircevent"
)

var (
	conn *irc.Connection
)

// TestMain is the "main" wrapper for testing.
func TestMain(m *testing.M) {
	conn = irc.IRC("elemental_oper", "user")
	err := conn.Connect(os.Getenv("IRCD_PORT_6667_TCP_ADDR") + ":6667")
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	conn.Quit()

	os.Exit(code)
}

// TestFailOperUp tests a failure of opering up.
func TestFailOperUp(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	conn.AddCallback("381", func(e *irc.Event) {
		t.Fatal(e.Message())
	})

	conn.AddCallback("491", func(e *irc.Event) {
		wg.Done()
	})

	conn.SendRaw("OPER god ninjas")

	wg.Wait()
}

// TestOperUp tests being able to oper up.
func TestOperUp(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	conn.AddCallback("491", func(e *irc.Event) {
		t.Fatal(e.Message())
	})

	conn.AddCallback("381", func(e *irc.Event) {
		wg.Done()
	})

	conn.SendRaw("OPER god test")

	wg.Wait()
}

// TestKillClient tests killing a client.
func TestKillClient(t *testing.T) {
	conn.AddCallback("481", func(e *irc.Event) {
		t.Fatal(e.Message())
	})

	smurfconn := irc.IRC("elemental_test_smurf", "smurf")
	smurfconn.Connect(os.Getenv("IRCD_PORT_6667_TCP_ADDR") + ":6667")
	time.Sleep(500 * time.Millisecond)

	conn.SendRaw("KILL elemental_test_smurf :die you test smurf")

	time.Sleep(2 * time.Second)
}

// TestCantKillServer tests not being able to kill a server.
func TestCantKillServer(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	conn.AddCallback("483", func(e *irc.Event) {
		wg.Done()
	})

	conn.SendRaw("KILL services.int")

	wg.Wait()
}
