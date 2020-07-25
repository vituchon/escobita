package model

import (
	"testing"
)

func TestMatchFlow(t *testing.T) {
	var party []Player = []Player{
		Player{Name: "P1"},
		Player{Name: "P2"},
		Player{Name: "P3"},
		Player{Name: "P4"},
		Player{Name: "P5"},
		Player{Name: "P6"},
	}
	testRuns := []struct {
		title          string
		players        []Player
		expectedRounds int
	}{
		{
			title:          "match flow for 2 players",
			players:        party[:2],
			expectedRounds: 6,
		},
		/*{
			title:          "match flow for 3 players",
			players:        party[:3],
			expectedRounds: 4,
		},
		{
			title:          "match flow for 4 players",
			players:        party[:4],
			expectedRounds: 3,
		},
		{
			title:          "match flow for 5 players",
			players:        party[:5],
			expectedRounds: 2, // it doesn't fit perfect as 36(cards) / (5(player) * 3(cards)) = 2,4 rounds
		},
		{
			title:          "match flow for 6 players",
			players:        party[:6],
			expectedRounds: 2,
		},*/
	}

	for _, testRun := range testRuns {
		t.Logf("Running unit test: %s", testRun.title)
		match := CreateAndServe(testRun.players)
		actualRounds := 0
		//t.Log(actualRounds, match)
		for match.MatchCanHaveMoreRounds() && actualRounds < 10 { // the match ends by dealing cards until there a no more cards left
			actualRounds++
			actualTurns := 0
			round := match.NextRound()
			addressedPlayers := make(map[Player]bool)
			for round.HasNextTurn() && actualTurns < 20 { // the round ends when each player consumes his turn
				player := round.NextTurn()
				//t.Log(player)
				_, exists := addressedPlayers[player]
				if exists {
					t.Error("Players should be addressed once per round")
				}
				match.Drop(player, match.Cards.PerPlayer[player].Hand[0]) // each players just drops a player
				actualTurns++
			}
			expectedTurns := len(match.Players) * 3
			if actualTurns != expectedTurns {
				t.Errorf("Round should have %d turns and they were %d", expectedTurns, actualTurns)
			}
		}
		if actualRounds != testRun.expectedRounds {
			t.Errorf("Match should end before %d rounds and rounds were %d", testRun.expectedRounds, actualRounds)
		}
	}
}

const MOCK_ID = 0

func TestCardsCombinationValues(t *testing.T) {

	testRuns := []struct {
		title         string
		cards         []Card
		expectedValue int
	}{
		{
			title:         "1 + 1 + 1 + 1 = 4",
			cards:         []Card{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, SWORD, 1}, Card{MOCK_ID, CLUB, 1}, Card{MOCK_ID, CUP, 1}},
			expectedValue: 1 + 1 + 1 + 1,
		},
		{
			title:         "2 + 2 + 1 + 1 = 6",
			cards:         []Card{Card{MOCK_ID, GOLD, 2}, Card{MOCK_ID, SWORD, 2}, Card{MOCK_ID, CLUB, 1}, Card{MOCK_ID, CUP, 1}},
			expectedValue: 2 + 2 + 1 + 1,
		},
		{
			title:         "3 + 4 + 5 + 6 + 7 = 25",
			cards:         []Card{Card{MOCK_ID, GOLD, 3}, Card{MOCK_ID, SWORD, 4}, Card{MOCK_ID, CLUB, 5}, Card{MOCK_ID, CUP, 6}, Card{MOCK_ID, CUP, 7}},
			expectedValue: 3 + 4 + 5 + 6 + 7,
		},
		{
			title:         "8 + 8 + 1 = 17",
			cards:         []Card{Card{MOCK_ID, GOLD, 10}, Card{MOCK_ID, SWORD, 10}, Card{MOCK_ID, CLUB, 1}},
			expectedValue: 8 + 8 + 1,
		},
		{
			title:         "8 + 9 + 10 = 27",
			cards:         []Card{Card{MOCK_ID, GOLD, 10}, Card{MOCK_ID, SWORD, 11}, Card{MOCK_ID, CLUB, 12}},
			expectedValue: 8 + 9 + 10,
		},
		{
			title:         "7 + 7 + 1 = 15",
			cards:         []Card{Card{MOCK_ID, CUP, 7}, Card{MOCK_ID, SWORD, 7}, Card{MOCK_ID, SWORD, 1}},
			expectedValue: 15,
		},
		{
			title:         "9 + 10 + 1 = 20",
			cards:         []Card{Card{MOCK_ID, CLUB, 11}, Card{MOCK_ID, CUP, 12}, Card{MOCK_ID, CLUB, 1}},
			expectedValue: 20,
		},
	}
	for _, testRun := range testRuns {
		t.Logf("Running unit test: %s ", testRun.title)
		computedValue := sumValues(testRun.cards)
		if computedValue != testRun.expectedValue {
			t.Errorf("aggregated card values differs! Expected is %d and computed value is %d", testRun.expectedValue, computedValue)
		}
	}
}
