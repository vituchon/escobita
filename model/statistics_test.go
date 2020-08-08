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
	match.FirstPlayerIndex = 0                               // beto
	match.Cards.Board = copyDeck(oneRoundTwoPlayersDeck[:4]) // 14 de puntaje en la mesa
	match.Cards.Left = oneRoundTwoPlayersDeck[4:]
	// end: hardcoding a serve

	round := match.NextRound() // round 1

	beto := round.NextTurn()
	// take: 1 hand + (5+5+3+1) board = 15, escobita
	takeCard(t, beto, &match, match.Cards.PerPlayer[beto].Hand[0], match.Cards.Board, true)

	pepe := round.NextTurn()
	// drop: 8
	dropCard(t, pepe, &match, match.Cards.PerPlayer[pepe].Hand[0])

	beto = round.NextTurn()
	// drop: 3
	dropCard(t, beto, &match, match.Cards.PerPlayer[beto].Hand[0])

	pepe = round.NextTurn()
	// take: 4 hand + (8+3) board = 15, escobita
	takeCard(t, pepe, &match, match.Cards.PerPlayer[pepe].Hand[0], match.Cards.Board, true)

	staticticsByPlayer := match.CalculateStaticticsByPlayer()
	if staticticsByPlayer[beto].EscobitasCount != 1 {
		t.Errorf("Beto shall have 1 escobita")
	}
	if staticticsByPlayer[pepe].EscobitasCount != 1 {
		t.Errorf("Pepe shall have 1 escobita")
	}

	if match.HasMoreRounds() {
		t.Errorf("Match is one round only")
	}

}

func TestGoldSeven(t *testing.T) {
	var players []Player = []Player{
		Player{Name: "Beto"},
		Player{Name: "Pepe"},
	}

	var oneRoundTwoPlayersDeck Deck = Deck{
		Card{1, GOLD, 3},  // 1 t (beto)
		Card{2, CUP, 2},   // 1 t (beto)
		Card{3, CLUB, 3},  // 1 t (beto)
		Card{4, SWORD, 5}, // 2 t (pepe)
		// las de beto
		Card{5, GOLD, 7}, // 1 t (beto)
		Card{6, CUP, 3},  // 5 d (beto)
		Card{7, CLUB, 1}, // 3 d (beto) - 6 t (pepe)
		// las de pepe
		Card{8, SWORD, 10}, // 4 d (pepe) - 6 t (pepe)
		Card{9, GOLD, 4},   // 6 t (pepe)
		Card{10, CUP, 10},  // 2 t (pepe)
	}

	match := newMatch(players, oneRoundTwoPlayersDeck)

	// begin: hardcoding a serve
	match.FirstPlayerIndex = 0                               // beto
	match.Cards.Board = copyDeck(oneRoundTwoPlayersDeck[:4]) // 12 de puntaje en la mesa
	match.Cards.Left = copyDeck(oneRoundTwoPlayersDeck[4:])
	// end: hardcoding a serve

	round := match.NextRound() // round 1

	beto := round.NextTurn()
	// take: 7 hand + (3+2+3) board = 15 (GOLD 7)
	boardCards, err := match.Cards.Board.GetMultiple(1, 2, 3)
	takeCard(t, beto, &match, match.Cards.PerPlayer[beto].Hand[0], boardCards, false)
	t.Logf("%+v", match.Cards.PerPlayer[beto])

	pepe := round.NextTurn()
	// take: 10 hand + (5) board = 15, escobita
	boardCards, err = match.Cards.Board.GetMultiple(4)
	takeCard(t, pepe, &match, match.Cards.PerPlayer[pepe].Hand[2], boardCards, true)

	beto = round.NextTurn()
	// drop: 1
	dropCard(t, beto, &match, match.Cards.PerPlayer[beto].Hand[1])

	pepe = round.NextTurn()
	// drop: 10
	dropCard(t, pepe, &match, match.Cards.PerPlayer[pepe].Hand[0])

	beto = round.NextTurn()
	// drop: 3
	dropCard(t, beto, &match, match.Cards.PerPlayer[beto].Hand[0])

	pepe = round.NextTurn()
	// take: 4 hand + (11) board = 15,
	boardCards, err = match.Cards.Board.GetMultiple(7, 8)
	if err != nil {
		t.Fatalf("unexpected error: '%v'", err)
	}
	takeCard(t, pepe, &match, match.Cards.PerPlayer[pepe].Hand[0], boardCards, false)

	staticticsByPlayer := match.CalculateStaticticsByPlayer()
	if !staticticsByPlayer[beto].HasGoldSeven {
		t.Errorf("Beto shall have gold seven!")
	}
	if staticticsByPlayer[beto].EscobitasCount != 0 {
		t.Errorf("Pepe shall have no escobita")
	}
	if staticticsByPlayer[pepe].EscobitasCount != 1 {
		t.Errorf("Pepe shall have 1 escobita")
	}
	if staticticsByPlayer.calculateMostGoldCardsPlayer() != beto {
		t.Errorf("beto shall be the player with most gold cards")
	}

	scoreSummaryByPlayer := staticticsByPlayer.BuildScoreBoard()

	if scoreSummaryByPlayer[beto].Score != 3 {
		t.Errorf("pepe shall have score 3 and pepe computed is %d", scoreSummaryByPlayer[beto].Score)
	}
	if scoreSummaryByPlayer[pepe].Score != 2 {
		t.Errorf("pepe shall have score 2 and pepe computed is %d", scoreSummaryByPlayer[pepe].Score)
	}

	//t.Logf("match %+v\n", match)
	//t.Logf("scoreSummaryByPlayer %+v\n", scoreSummaryByPlayer)

	if match.HasMoreRounds() {
		t.Errorf("Match is one round only")
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
}
