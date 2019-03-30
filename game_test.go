package poker_test

import (
	"fmt"
	"testing"
	"time"

	poker "github.com/bjrnt/tdd-app"
)

func TestTexasHoldem_Start(t *testing.T) {
	t.Run("it schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(5)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		assertSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("it schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(7)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		assertSchedulingCases(t, cases, blindAlerter)
	})
}

func TestTexasHoldem_Finish(t *testing.T) {
	store := &poker.StubPlayerStore{}
	game := poker.NewTexasHoldem(dummySpyAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	poker.AssertPlayerWin(t, store, winner)
}

func assertSchedulingCases(t *testing.T, want []scheduledAlert, blindAlerter *SpyBlindAlerter) {
	for i, want := range want {
		t.Run(fmt.Sprintf(fmt.Sprint(want)), func(t *testing.T) {
			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}

			got := blindAlerter.alerts[i]
			assertScheduledAlert(t, got, want)
		})
	}
}
