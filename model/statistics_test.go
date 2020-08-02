package model

import (
	"testing"
)

/*
func TestSeventiesCalculation(t *testing.T) {
	testRuns := []struct {
		title         string
		cards         []Card
		expectedValue int
	}{
		{
			title:         "1 + 1 + 1 + 1 = 8",
			cards:         []Card{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, SWORD, 1}, Card{MOCK_ID, CLUB, 1}, Card{MOCK_ID, CUP, 1}},
			expectedValue: 8,
		},
	}
	for _, testRun := range testRuns {
		t.Logf("Running unit test: %s ", testRun.title)
		computedValue := calculateSeventiesScore(testRun.cards...)
		if computedValue != testRun.expectedValue {
			t.Errorf("aggregated card values differs! Expected is %d and computed value is %d", testRun.expectedValue, computedValue)
		}
	}
}*/

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

	// hardcoding a serve
	match.FirstPlayerIndex = 0                                    // beto
	match.MatchCards.Board = copyDeck(oneRoundTwoPlayersDeck[:4]) // 14 de puntaje en la mesa
	match.MatchCards.Left = oneRoundTwoPlayersDeck[4:]
	// end

	round := match.NextRound() // round 1

	beto := round.NextTurn()
	takeAction := PlayerTakeAction{
		HandCard:   match.MatchCards.PerPlayer[beto].Hand[0], // 1 de puntaje
		BoardCards: match.MatchCards.Board,                   // 14 de puntaje
	}
	escobitaAction := match.Take(beto, takeAction) // 1 hace una escovid
	if !escobitaAction.IsEscobita() {
		t.Fatalf("player take action shall be escobita and actual result is not")
	}

	pepe := round.NextTurn()
	dropAction := PlayerDropAction{
		HandCard: match.MatchCards.PerPlayer[pepe].Hand[0], // 8 de puntaje
	}
	match.Drop(pepe, dropAction)

	beto = round.NextTurn()
	dropAction = PlayerDropAction{
		HandCard: match.MatchCards.PerPlayer[beto].Hand[0], // 3 de puntaje
	}
	match.Drop(beto, dropAction)

	pepe = round.NextTurn()
	takeAction = PlayerTakeAction{
		HandCard:   match.MatchCards.PerPlayer[pepe].Hand[0], // 4 de puntaje
		BoardCards: match.MatchCards.Board,                   // 11 de puntaje
	}
	escobitaAction = match.Take(pepe, takeAction)
	if !escobitaAction.IsEscobita() {
		t.Fatalf("player take action shall be escobita and actual result is not")
	}

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

func dropCard(t *testing.T, player Player, match Match, card Card) {

}
