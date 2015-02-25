/*
Tests for basic channel moderation.
*/
package chanmod

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/thoj/go-ircevent"
)

var (
	conn       *irc.Connection
	peonconn   *irc.Connection
	voicedconn *irc.Connection
)

// TestMain is the "main" wrapper for testing.
func TestMain(m *testing.M) {
	conn = irc.IRC("elemental_chanop", "user")
	err := conn.Connect(os.Getenv("IRCD_PORT_6667_TCP_ADDR") + ":6667")
	if err != nil {
		log.Fatal(err)
	}
	go conn.Loop()

	peonconn = irc.IRC("elemental_peon", "user")
	err = peonconn.Connect(os.Getenv("IRCD_PORT_6667_TCP_ADDR") + ":6667")
	if err != nil {
		log.Fatal(err)
	}
	go peonconn.Loop()

	voicedconn = irc.IRC("elemental_voiced", "user")
	err = voicedconn.Connect(os.Getenv("IRCD_PORT_6667_TCP_ADDR") + ":6667")
	if err != nil {
		log.Fatal(err)
	}
	go voicedconn.Loop()

	code := m.Run()

	conn.Quit()
	peonconn.Quit()

	os.Exit(code)
}

// TestJoinChannelAndBeAnOp
func TestBasicConnection(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	conn.AddCallback("353", func(e *irc.Event) {
		// :rarity.shadownet.int 353 a = #yolo :@a
		// ^                     ^   ^
		// |                     |   \- Arguments
		// \- source             |
		//                       \- verb (Code here)
		if e.Arguments[3] != "@"+conn.GetNick() {
			t.Fatal("Expected to be an op in #yolo")
		}

		wg.Done()
	})

	conn.Join("#yolo")

	wg.Wait()
}

func TestMakePeonJoin(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	peonconn.AddCallback("353", func(e *irc.Event) {
		// :rarity.shadownet.int 353 a = #yolo :@elemental_chanop elemental_peon
		// ^                     ^   ^
		// |                     |   \- Arguments
		// \- source             |
		//                       \- verb (Code here)
		if e.Arguments[3] != peonconn.GetNick()+" @"+conn.GetNick() {
			t.Fatal("Expected to be a peon in #yolo")
		}

		wg.Done()
	})

	peonconn.Join("#yolo")

	wg.Wait()
}

func TestPeonCannotKick(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	peonconn.AddCallback("482", func(e *irc.Event) {
		wg.Done()
	})

	peonconn.SendRaw("KICK #yolo elemental_chanop")

	wg.Wait()

	wg.Add(9999999)
}

func TestPeonCannotSetChannelProperties(t *testing.T) {
	modes := []string{
		"n", "t", "s", "p", "m", "i", "r", "c", "d", "g", "z", "u",
		"L", "P", "F", "Q", "C", "D", "T", "E", "J", "K",
	}

	wg := sync.WaitGroup{}
	wg.Add(len(modes))

	peonconn.AddCallback("482", func(e *irc.Event) {
		wg.Done()
	})

	peonconn.AddCallback("481", func(e *irc.Event) {
		wg.Done()
	})

	conn.AddCallback("MODE", func(e *irc.Event) {
		t.Fatal(e.Message())
	})

	for _, mode := range modes {
		peonconn.Mode("#yolo", "+"+mode)
	}

	wg.Wait()
	wg.Add(999999)
}

func TestPeonCannotSetChannelBans(t *testing.T) {
	modes := []string{"q", "e", "I", "b"}

	wg := sync.WaitGroup{}
	wg.Add(len(modes))

	peonconn.AddCallback("482", func(e *irc.Event) {
		wg.Done()
	})

	for _, mode := range modes {
		peonconn.Mode("#yolo", "+"+mode+" swag")
	}

	wg.Wait()
	wg.Add(999999)
}

func TestPeonCannotSetKeyMode(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	peonconn.AddCallback("482", func(e *irc.Event) {
		wg.Done()
	})

	peonconn.Mode("#yolo", "+k ninjas")

	wg.Wait()
	wg.Add(999999)
}

func TestPeonCannotSpeakOverStuff(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	peonconn.AddCallback("404", func(e *irc.Event) {
		wg.Done()
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		t.Fatal(e.Message())
	})

	conn.Mode("#yolo", "+m")
	peonconn.Privmsg("#yolo", "nuuuuuuuuuuuuuuuuu")
	wg.Wait()
	wg.Add(1)

	conn.Mode("#yolo", "+b "+peonconn.GetNick())
	peonconn.Privmsg("#yolo", "nuuuuuuuuuuuuuuuuu")
	wg.Wait()
	wg.Add(1)

	conn.Mode("#yolo", "-b+q "+peonconn.GetNick()+" "+peonconn.GetNick())
	peonconn.Privmsg("#yolo", "nuuuuuuuuuuuuuuuuu")
	wg.Wait()

	conn.Mode("#yolo", "-q "+peonconn.GetNick())

	wg.Add(999999)
}

func TestVoicedUserInitialSetup(t *testing.T) {
	voicedconn.Join("#yolo")
	time.Sleep(250 * time.Millisecond)
	conn.Mode("#yolo", "+v "+voicedconn.GetNick())
}

func TestVoicedUserCannotKick(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	voicedconn.AddCallback("482", func(e *irc.Event) {
		wg.Done()
	})

	voicedconn.SendRaw("KICK #yolo elemental_chanop")

	wg.Wait()

	wg.Add(9999999)
}

func TestVoicedUserCannotSetChannelProperties(t *testing.T) {
	modes := []string{
		"n", "t", "s", "p", "m", "i", "r", "c", "d", "g", "z", "u",
		"L", "P", "F", "Q", "C", "D", "T", "E", "J", "K",
	}

	wg := sync.WaitGroup{}
	wg.Add(len(modes))

	voicedconn.AddCallback("482", func(e *irc.Event) {
		wg.Done()
	})

	voicedconn.AddCallback("481", func(e *irc.Event) {
		wg.Done()
	})

	conn.AddCallback("MODE", func(e *irc.Event) {
		t.Fatal(e.Message())
	})

	for _, mode := range modes {
		voicedconn.Mode("#yolo", "+"+mode)
	}

	wg.Wait()
	wg.Add(999999)
}

func TestVoicedUserCannotSetChannelBans(t *testing.T) {
	modes := []string{"q", "e", "I", "b"}

	wg := sync.WaitGroup{}
	wg.Add(len(modes))

	voicedconn.AddCallback("482", func(e *irc.Event) {
		wg.Done()
	})

	for _, mode := range modes {
		voicedconn.Mode("#yolo", "+"+mode+" swag")
	}

	wg.Wait()
	wg.Add(999999)
}

func TestVoicedUserCanSpeakOverStuff(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		wg.Done()
	})

	voicedconn.AddCallback("404", func(e *irc.Event) {
		t.Fatal(e.Message())
	})

	conn.Mode("#yolo", "+m")
	voicedconn.Privmsg("#yolo", "nuuuuuuuuuuuuuuuuu")
	wg.Wait()
	wg.Add(1)

	conn.Mode("#yolo", "+b "+voicedconn.GetNick())
	voicedconn.Privmsg("#yolo", "nuuuuuuuuuuuuuuuuu")
	wg.Wait()
	wg.Add(1)

	conn.Mode("#yolo", "-b+q "+voicedconn.GetNick()+" "+peonconn.GetNick())
	voicedconn.Privmsg("#yolo", "nuuuuuuuuuuuuuuuuu")
	wg.Wait()

	conn.Mode("#yolo", "-q "+voicedconn.GetNick())

	wg.Add(999999)
}
