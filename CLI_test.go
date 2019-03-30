package poker_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	poker "github.com/bjrnt/tdd-app"
)

type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{duration, amount})
}

type GameSpy struct {
	StartedWith  int
	FinishedWith string
	StartCalled  bool
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}

var dummySpyAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

func TestCLI(t *testing.T) {
	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("3", "Chris wins")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Chris")
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &GameSpy{}

		in := userSends("8", "Cleo wins")
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Cleo")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		game := &GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("pies")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})
}

func userSends(messages ...string) io.Reader {
	rdr := strings.NewReader(strings.Join(messages, "\n"))
	return rdr
}

func assertFinishCalledWith(t *testing.T, game *GameSpy, want string) {
	t.Helper()
	got := game.FinishedWith
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertGameStartedWith(t *testing.T, game *GameSpy, numberOfPlayers int) {
	t.Helper()
	got := game.StartedWith
	want := numberOfPlayers
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertGameNotStarted(t *testing.T, game *GameSpy) {
	t.Helper()
	if game.StartCalled {
		t.Errorf("game should not have been started")
	}
}

func assertMessagesSentToUser(t *testing.T, out *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := out.String()
	if got != want {
		t.Errorf("got '%s' sent to out but expected %+v", got, messages)
	}
}

func assertScheduledAlert(t *testing.T, got scheduledAlert, want scheduledAlert) {
	t.Helper()

	if got.amount != want.amount {
		t.Errorf("got amount %d want %d", got.amount, want.amount)
	}

	if got.at != want.at {
		t.Errorf("got %v want %v", got.at, want.at)
	}
}
