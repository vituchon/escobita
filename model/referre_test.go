package model

import (
	"testing"
)

// TODO: Add at least a test to verify for solo player game, if the player perform only ONE take action... and the end he will got all the 40 cards taken.. so i can check that!

func TestMatchDropFlow(t *testing.T) {
	var party []Player = []Player{
		Player{Id: 1, Name: "P1"},
		Player{Id: 2, Name: "P2"},
		Player{Id: 3, Name: "P3"},
		Player{Id: 4, Name: "P4"},
		Player{Id: 5, Name: "P5"},
		Player{Id: 6, Name: "P6"},
	}
	testRuns := []struct {
		title          string
		players        []Player
		expectedRounds int
	}{
		{
			title:          "match flow for 1 players",
			players:        party[:1],
			expectedRounds: 12,
		},
		{
			title:          "match flow for 2 players",
			players:        party[:2],
			expectedRounds: 6,
		},
		{
			title:          "match flow for 3 players",
			players:        party[:3],
			expectedRounds: 4,
		},
		{
			title:          "match flow for 4 players",
			players:        party[:4],
			expectedRounds: 3,
		},
		/*{
			title:          "match flow for 5 players",
			players:        party[:5],
			expectedRounds: 2, // it doesn't fit perfect as 36(cards) / (5(player) * 3(cards)) = 2,4 rounds
		},*/
		{
			title:          "match flow for 6 players",
			players:        party[:6],
			expectedRounds: 2,
		},
	}

	for _, testRun := range testRuns {
		t.Logf("==== Running unit test: %s ====", testRun.title)
		match := CreateMatch(testRun.players)
		match.Prepare()
		actualRounds := 0
		//t.Log(actualRounds, match)
		for match.HasMoreRounds() && actualRounds < 100 { // the match ends by dealing cards until there a no more cards left
			actualRounds++
			actualTurns := 0
			round := match.NextRound()
			turnsCountByPlayer := make(map[Player]int)
			for round.HasNextTurn() { // the round ends when each player consumes his three turns per round
				player := round.NextTurn()
				//t.Log("turno: ", actualRounds, ", jugador: ", player, "cartas:", match.MatchCards.PerPlayer[player])
				//t.Log(match)
				t.Log(turnsCountByPlayer)
				count, _ := turnsCountByPlayer[player]
				turnsCountByPlayer[player] = count + 1
				dropAction := NewPlayerDropAction(player, match.Cards.ByPlayer[player].Hand[0])

				//t.Logf("dropAction: %+v por parte de %v\n ", dropAction, player)
				_, err := match.Drop(dropAction) // each players just drops a player
				if err != nil {
					t.Error("Unexpected error performing a drop action", err)
				}
				actualTurns++
			}
			for player, turnsCount := range turnsCountByPlayer {
				if turnsCount != 3 {
					t.Errorf("Player %v has %d turns and expected turns are 3", player, turnsCount)
				}
			}

			if round.Number != actualRounds {
				t.Errorf("Round number is %d turns and round.Number is %d", actualRounds, round.Number)
			}
			expectedTurns := len(match.Players) * 3
			if actualTurns != expectedTurns {
				t.Errorf("Round should have %d turns and they were %d", expectedTurns, actualTurns)
			}
		}
		if actualRounds != testRun.expectedRounds {
			t.Errorf("Match should end before %d rounds and rounds were %d", testRun.expectedRounds, actualRounds)
		}

		const actionsLogCount int = 36
		if len(match.ActionsLog) != actionsLogCount {
			t.Errorf("Match quantity of actions must be 36 and computed is %d", len(match.ActionsLog))
		}

		match.Ends()
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

func TestDropActionCanBePerformedOnlyByTheTurnPlayer(t *testing.T) {
	var player1 Player = Player{Id: 1, Name: "P1"}
	var player2 Player = Player{Id: 2, Name: "P2"}
	var party []Player = []Player{
		player1,
		player2,
	}
	testRuns := []struct {
		title   string
		players []Player
	}{
		{
			title:   "Not turn's player tries to perform a drop action",
			players: party,
		},
	}

	for _, testRun := range testRuns {
		t.Logf("==== Running unit test: %s ====", testRun.title)
		match := CreateMatch(testRun.players)
		match.Prepare()
		round := match.NextRound()
		player := round.NextTurn()
		var notTurnPlayer *Player
		if player.Id == player1.Id {
			notTurnPlayer = &player2
		} else {
			notTurnPlayer = &player1
		}
		t.Log("Current turn's player is", player, "and player", notTurnPlayer, "will try to perform a drop action.")
		_, err := dropCard(t, *notTurnPlayer, &match, Card{})
		if err == nil {
			t.Error("An error was expected")
		} else {
			t.Log("Expected error happens", err)
		}
	}
}

// Helper method: (Evaluate if necessary) Starts a match, beginning the first round and then the first turn within that first round
func (match *Match) Begins() (actingPlayer Player) {
	match.Prepare()
	match.NextRound()                    // advances into first round
	return match.CurrentRound.NextTurn() // advances into first turn
}
