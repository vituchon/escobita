package model

import (
	"testing"
)

func TestSeventiesCalculation(t *testing.T) {
	testRuns := []struct {
		title         string
		cards         []Card
		expectedValue int
	}{
		{
			title:         "1 (GOLD) + 1 (SWORD) + 1 (CLUB) + 1 (CUP) = 1 + 1 + 1 + 1",
			cards:         []Card{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, SWORD, 1}, Card{MOCK_ID, CLUB, 1}, Card{MOCK_ID, CUP, 1}},
			expectedValue: 4,
		},
		{
			title:         "3 (GOLD) + 1 (CUP) = 4 + 1",
			cards:         Deck{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, GOLD, 2}, Card{MOCK_ID, GOLD, 3}, Card{MOCK_ID, CUP, 1}},
			expectedValue: 5,
		},
		{
			title:         "4 (GOLD) + 5 (CUP) = 8 + 16",
			cards:         Deck{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, GOLD, 4}, Card{MOCK_ID, CUP, 5}, Card{MOCK_ID, CUP, 1}},
			expectedValue: 24,
		},
		{
			title:         "6 (GOLD) + 1 (CUP) + 7 (CUP) = 32 + 2 + 64",
			cards:         Deck{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, GOLD, 6}, Card{MOCK_ID, GOLD, 3}, Card{MOCK_ID, CUP, 2}, Card{MOCK_ID, SWORD, 6}, Card{MOCK_ID, SWORD, 7}},
			expectedValue: 98,
		},
		// same as above plus not valid cards for seventies (> 7)
		{
			title:         "1 (GOLD) + 1 (SWORD) + 1 (CLUB) + 1 (CUP) = 1 + 1 + 1 + 1",
			cards:         []Card{Card{MOCK_ID, GOLD, 10}, Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, SWORD, 1}, Card{MOCK_ID, CLUB, 1}, Card{MOCK_ID, CUP, 1}},
			expectedValue: 4,
		},
		{
			title:         "3 (GOLD) + 1 (CUP) = 4 + 1",
			cards:         Deck{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, GOLD, 2}, Card{MOCK_ID, GOLD, 3}, Card{MOCK_ID, CUP, 11}, Card{MOCK_ID, CUP, 1}},
			expectedValue: 5,
		},
		{
			title:         "4 (GOLD) + 5 (CUP) = 8 + 16",
			cards:         Deck{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, GOLD, 12}, Card{MOCK_ID, GOLD, 4}, Card{MOCK_ID, CUP, 5}, Card{MOCK_ID, CUP, 10}, Card{MOCK_ID, CUP, 1}},
			expectedValue: 24,
		},
		{
			title:         "6 (GOLD) + 1 (CUP) + 7 (CUP) = 32 + 2 + 64",
			cards:         Deck{Card{MOCK_ID, CLUB, 10}, Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, GOLD, 6}, Card{MOCK_ID, GOLD, 3}, Card{MOCK_ID, CUP, 2}, Card{MOCK_ID, SWORD, 6}, Card{MOCK_ID, GOLD, 11}, Card{MOCK_ID, SWORD, 7}},
			expectedValue: 98,
		},
	}
	for _, testRun := range testRuns {
		t.Logf("Running unit test: %s ", testRun.title)
		computedValue := calculateSeventiesScore(testRun.cards)
		if computedValue != testRun.expectedValue {
			t.Errorf("aggregated card values differs! Expected is %d and computed value is %d", testRun.expectedValue, computedValue)
		}
	}
}

/*
func TestEscobita(t *testing.T) {
	var players []Player = []Player{
		Player{Name: "Beto"},
		Player{Name: "Pepe"},
	}

	var oneRoundTwoPlayersDeck Deck = Deck{
		Card{1, GOLD, 5},
		Card{2, CUP, 5},
		Card{3, CLUB, 3},
		Card{4, SWORD, 1},
		// las de beto
		Card{5, GOLD, 3},
		Card{6, CUP, 1},
		Card{7, CLUB, 8},
		// las de pepe
		Card{8, SWORD, 10},
		Card{9, GOLD, 4},
		Card{10, CUP, 4},
	}

	match := newMatch(players, oneRoundTwoPlayersDeck)

	// begin: hardcoding a serve
	match.FirstPlayerIndex = 0                                    // beto
	match.MatchCards.Board = copyDeck(oneRoundTwoPlayersDeck[:4]) // 14 de puntaje en la mesa
	match.MatchCards.Left = oneRoundTwoPlayersDeck[4:]
	// end: hardcoding a serve

	round := match.NextRound() // round 1

	beto := round.NextTurn()
	// take: 1 hand + (5+5+3+1) board = 15, escobita
	takeCard(t, beto, &match, match.MatchCards.PerPlayer[beto].Hand[0], match.MatchCards.Board, true)

	pepe := round.NextTurn()
	// drop: 8
	dropCard(t, pepe, &match, match.MatchCards.PerPlayer[pepe].Hand[0])

	beto = round.NextTurn()
	// drop: 3
	dropCard(t, beto, &match, match.MatchCards.PerPlayer[beto].Hand[0])

	pepe = round.NextTurn()
	// take: 4 hand + (8+3) board = 15, escobita
	takeCard(t, pepe, &match, match.MatchCards.PerPlayer[pepe].Hand[0], match.MatchCards.Board, true)

	staticticsByPlayer := match.CalculateStaticticsByPlayer()
	if staticticsByPlayer[beto].EscobitasCount != 1 {
		t.Fatalf("Beto shall have 1 escobita")
	}
	if staticticsByPlayer[pepe].EscobitasCount != 1 {
		t.Fatalf("Pepe shall have 1 escobita")
	}

	if match.HasMoreRounds() {
		t.Fatalf("Match is one round only")
	}

}

func takeCard(t *testing.T, player Player, match *Match, handCard Card, boardCards []Card, mustBeEscobita bool) PlayerAction {
	takeAction := PlayerTakeAction{
		HandCard:   handCard,
		BoardCards: boardCards,
	}
	action := match.Take(player, takeAction)
	if mustBeEscobita && !action.IsEscobita() {
		t.Fatalf("player take action shall be escobita and actual result is not")
	}
	return action
}

func dropCard(t *testing.T, player Player, match *Match, card Card) PlayerAction {
	dropAction := PlayerDropAction{
		HandCard: card,
	}
	return match.Drop(player, dropAction)
}*/
