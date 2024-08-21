package model

import (
	"errors"
	"fmt"
	"testing"
)

func TestScoreCalculationOnBeginMatch(t *testing.T) {
	var beto = Player{Name: "beto"}
	var pepe = Player{Name: "pepe"}
	var players []Player = []Player{
		beto,
		pepe,
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
		Card{8, SWORD, 12}, // 4 d (pepe) - 6 t (pepe)
		Card{9, GOLD, 4},   // 6 t (pepe)
		Card{10, CUP, 12},  // 2 t (pepe)
	}
	match := newMatch(players, oneRoundTwoPlayersDeck)
	staticticsByPlayer := match.CalculateStaticticsByPlayer()
	scoreSummaryByPlayer := staticticsByPlayer.BuildScoreSummaryByPlayer()

	if scoreSummaryByPlayer[beto].Score != 0 {
		t.Errorf("beto hasn't make any action and computed score is %d, is must be 0", scoreSummaryByPlayer[beto].Score)
	}
	if scoreSummaryByPlayer[pepe].Score != 0 {
		t.Errorf("pepe hasn't make any action and computed score is %d, is must be 0", scoreSummaryByPlayer[pepe].Score)
	}
}

func TestSeventiesCalculation(t *testing.T) {
	testRuns := []struct {
		title                       string
		cards                       []Card
		expectedScore               int
		expectedSevenRankCardsCount int
	}{
		{
			title:                       "1 (GOLD) + 1 (SWORD) + 1 (CLUB) + 1 (CUP) = 1 + 1 + 1 + 1",
			cards:                       []Card{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, SWORD, 1}, Card{MOCK_ID, CLUB, 1}, Card{MOCK_ID, CUP, 1}},
			expectedScore:               4,
			expectedSevenRankCardsCount: 0,
		},
		{
			title:                       "3 (GOLD) + 1 (CUP) = 4 + 1",
			cards:                       Deck{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, GOLD, 2}, Card{MOCK_ID, GOLD, 3}, Card{MOCK_ID, CUP, 1}},
			expectedScore:               5,
			expectedSevenRankCardsCount: 0,
		},
		{
			title:                       "4 (GOLD) + 5 (CUP) = 8 + 16",
			cards:                       Deck{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, GOLD, 4}, Card{MOCK_ID, CUP, 5}, Card{MOCK_ID, CUP, 1}},
			expectedScore:               24,
			expectedSevenRankCardsCount: 0,
		},
		{
			title:                       "6 (GOLD) + 1 (CUP) + 7 (CUP) = 32 + 2 + 64",
			cards:                       Deck{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, GOLD, 6}, Card{MOCK_ID, GOLD, 3}, Card{MOCK_ID, CUP, 2}, Card{MOCK_ID, SWORD, 6}, Card{MOCK_ID, SWORD, 7}},
			expectedScore:               98,
			expectedSevenRankCardsCount: 1,
		},
		// same as above plus not valid cards for seventies (> 7)
		{
			title:                       "1 (GOLD) + 1 (SWORD) + 1 (CLUB) + 1 (CUP) = 1 + 1 + 1 + 1",
			cards:                       []Card{Card{MOCK_ID, GOLD, 10}, Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, SWORD, 1}, Card{MOCK_ID, CLUB, 1}, Card{MOCK_ID, CUP, 1}},
			expectedScore:               4,
			expectedSevenRankCardsCount: 0,
		},
		{
			title:                       "3 (GOLD) + 1 (CUP) = 4 + 1",
			cards:                       Deck{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, GOLD, 2}, Card{MOCK_ID, GOLD, 3}, Card{MOCK_ID, CUP, 11}, Card{MOCK_ID, CUP, 1}},
			expectedScore:               5,
			expectedSevenRankCardsCount: 0,
		},
		{
			title:                       "4 (GOLD) + 5 (CUP) = 8 + 16",
			cards:                       Deck{Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, GOLD, 12}, Card{MOCK_ID, GOLD, 4}, Card{MOCK_ID, CUP, 5}, Card{MOCK_ID, CUP, 10}, Card{MOCK_ID, CUP, 1}},
			expectedScore:               24,
			expectedSevenRankCardsCount: 0,
		},
		{
			title:                       "6 (GOLD) + 1 (CUP) + 7 (CUP) = 32 + 2 + 64",
			cards:                       Deck{Card{MOCK_ID, CLUB, 10}, Card{MOCK_ID, GOLD, 1}, Card{MOCK_ID, GOLD, 6}, Card{MOCK_ID, GOLD, 3}, Card{MOCK_ID, CUP, 2}, Card{MOCK_ID, SWORD, 6}, Card{MOCK_ID, GOLD, 11}, Card{MOCK_ID, SWORD, 7}},
			expectedScore:               98,
			expectedSevenRankCardsCount: 1,
		},
		{
			title:                       "7 (GOLD) + 7 (SWORD) + 7 (CUP) + 7 (CLUB) = 64 * 4 = 256",
			cards:                       Deck{Card{MOCK_ID, GOLD, 7}, Card{MOCK_ID, SWORD, 7}, Card{MOCK_ID, CUP, 7}, Card{MOCK_ID, CLUB, 7}},
			expectedScore:               256,
			expectedSevenRankCardsCount: 4,
		},
	}
	for _, testRun := range testRuns {
		t.Logf("Running unit test: %s ", testRun.title)
		computedScore := calculateSeventiesScore(testRun.cards)
		if computedScore != testRun.expectedScore {
			t.Errorf("aggregated card values differs! Expected is %d and computed value is %d", testRun.expectedScore, computedScore)
		}
		computedCountSevenRankCards := CountSevenRankCards(testRun.cards)
		if computedCountSevenRankCards != testRun.expectedSevenRankCardsCount {
			t.Errorf("CountSevenRankCards differs! Expected is %d and computed value is %d", testRun.expectedSevenRankCardsCount, computedCountSevenRankCards)
		}
	}
}

func TestEscobita(t *testing.T) {
	var players []Player = []Player{
		Player{Name: "Beto"},
		Player{Name: "Pepe"},
	}

	var oneRoundTwoPlayersDeck Deck = Deck{
		Card{1, GOLD, 5},  // 1 t (beto)
		Card{2, CUP, 5},   // 1 t (beto)
		Card{3, CLUB, 3},  // 1 t (beto)
		Card{4, SWORD, 1}, // 1 t (beto)
		// las de beto
		Card{5, GOLD, 3}, // 3 d (beto)
		Card{6, CUP, 1},  // 1 t (beto) - 4 t (pepe)
		Card{7, CLUB, 7}, // 5 d (beto)
		// las de pepe
		Card{8, SWORD, 10}, // 2 d (pepe) - 4 t (pepe)
		Card{9, GOLD, 4},   // 4 t (pepe)
		Card{10, CUP, 4},   // 6 d (pepe)
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
	takeCard(t, beto, &match, match.Cards.ByPlayer[beto].Hand[1], match.Cards.Board, true)

	pepe := round.NextTurn()
	// drop: 10
	dropCard(t, pepe, &match, match.Cards.ByPlayer[pepe].Hand[0])

	beto = round.NextTurn()
	// drop: 3
	dropCard(t, beto, &match, match.Cards.ByPlayer[beto].Hand[0])

	pepe = round.NextTurn()
	// take: 4 hand + (8+3) board = 15, escobita
	takeCard(t, pepe, &match, match.Cards.ByPlayer[pepe].Hand[0], match.Cards.Board, true)

	beto = round.NextTurn()
	// drop: 7
	dropCard(t, beto, &match, match.Cards.ByPlayer[beto].Hand[0])

	pepe = round.NextTurn()
	// drop: 4
	dropCard(t, pepe, &match, match.Cards.ByPlayer[pepe].Hand[0])

	if match.HasMoreRounds() {
		t.Errorf("Match is one round only")
	}
	match.Ends()
	staticticsByPlayer := match.CalculateStaticticsByPlayer()
	//scoreSummaryByPlayer := staticticsByPlayer.BuildScoreBoard()
	//t.Logf("scoreSummaryByPlayer %+v\n", scoreSummaryByPlayer)

	if staticticsByPlayer[beto].EscobitasCount != 1 {
		t.Errorf("Beto shall have 1 escobita")
	}
	if staticticsByPlayer[pepe].EscobitasCount != 1 {
		t.Errorf("Pepe shall have 1 escobita")
	}

	if staticticsByPlayer[beto].CardsTakenCount != 5 {
		t.Errorf("Beto shall have 5 cards taken")
	}
	if staticticsByPlayer[pepe].CardsTakenCount != 5 {
		t.Errorf("Pepe shall have 5 cards taken")
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
		Card{8, SWORD, 12}, // 4 d (pepe) - 6 t (pepe)
		Card{9, GOLD, 4},   // 6 t (pepe)
		Card{10, CUP, 12},  // 2 t (pepe)
	}

	match := newMatch(players, oneRoundTwoPlayersDeck)

	// TODO : try to use match.Serve()
	// begin: hardcoding a serve
	match.FirstPlayerIndex = 0                               // beto
	match.Cards.Board = copyDeck(oneRoundTwoPlayersDeck[:4]) // 12 de puntaje en la mesa
	match.Cards.Left = copyDeck(oneRoundTwoPlayersDeck[4:])
	// end: hardcoding a serve

	round := match.NextRound() // round 1

	beto := round.NextTurn()
	// take: 7 hand + (3+2+3) board = 15 (GOLD 7)
	boardCards, err := match.Cards.Board.GetMultiple(1, 2, 3)
	takeCard(t, beto, &match, match.Cards.ByPlayer[beto].Hand[0], boardCards, false)
	//t.Logf("%+v", match.Cards.ByPlayer[beto])

	pepe := round.NextTurn()
	// take: 10 hand + (5) board = 15, escobita
	boardCards, err = match.Cards.Board.GetMultiple(4)
	takeCard(t, pepe, &match, match.Cards.ByPlayer[pepe].Hand[2], boardCards, true)

	beto = round.NextTurn()
	// drop: 1
	dropCard(t, beto, &match, match.Cards.ByPlayer[beto].Hand[1])

	pepe = round.NextTurn()
	// drop: 10
	dropCard(t, pepe, &match, match.Cards.ByPlayer[pepe].Hand[0])

	beto = round.NextTurn()
	// drop: 3
	dropCard(t, beto, &match, match.Cards.ByPlayer[beto].Hand[0])

	pepe = round.NextTurn()
	// take: 4 hand + (11) board = 15,
	boardCards, err = match.Cards.Board.GetMultiple(7, 8)
	if err != nil {
		t.Fatalf("unexpected error: '%v'", err)
	}
	takeCard(t, pepe, &match, match.Cards.ByPlayer[pepe].Hand[0], boardCards, false)

	if match.HasMoreRounds() {
		t.Errorf("Match is one round only")
	}
	match.Ends()
	staticticsByPlayer := match.CalculateStaticticsByPlayer()
	scoreSummaryByPlayer := staticticsByPlayer.BuildScoreSummaryByPlayer()
	t.Logf("scoreSummaryByPlayer %+v\n", scoreSummaryByPlayer)

	if !staticticsByPlayer[beto].HasGoldSeven {
		t.Errorf("beto shall have gold seven!")
	}
	if staticticsByPlayer[beto].EscobitasCount != 0 {
		t.Errorf("beto shall have no escobita")
	}
	if staticticsByPlayer[pepe].EscobitasCount != 1 {
		t.Errorf("pepe shall have 1 escobita")
	}
	if *staticticsByPlayer.calculateMostGoldCardsPlayer() != beto {
		t.Errorf("beto shall be the player with most gold cards")
	}

	if scoreSummaryByPlayer[beto].Score != 3 {
		t.Errorf("beto shall have score 3 and pepe computed is %d", scoreSummaryByPlayer[beto].Score)
	}
	if scoreSummaryByPlayer[pepe].Score != 2 {
		t.Errorf("pepe shall have score 2 and pepe computed is %d", scoreSummaryByPlayer[pepe].Score)
	}
}

func TestLastTaker(t *testing.T) {
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
		Card{5, GOLD, 7}, // 3 d (beto)
		Card{6, CUP, 3},  // 5 d (beto)
		Card{7, CLUB, 7}, // 1 t (beto)
		// las de pepe
		Card{8, SWORD, 12}, // 4 d (pepe)
		Card{9, GOLD, 4},   // 6 d (pepe)
		Card{10, CUP, 12},  // 2 t (pepe)
	}

	match := newMatch(players, oneRoundTwoPlayersDeck)

	// TODO : try to use match.Serve()
	// begin: hardcoding a serve
	match.FirstPlayerIndex = 0                               // beto
	match.Cards.Board = copyDeck(oneRoundTwoPlayersDeck[:4]) // 12 de puntaje en la mesa
	match.Cards.Left = copyDeck(oneRoundTwoPlayersDeck[4:])
	// end: hardcoding a serve

	round := match.NextRound() // round 1

	beto := round.NextTurn()
	// take: 7 hand + (3+2+3) board = 15
	boardCards, _ := match.Cards.Board.GetMultiple(1, 2, 3)
	takeCard(t, beto, &match, match.Cards.ByPlayer[beto].Hand[2], boardCards, false)
	t.Logf("%+v", match.Cards.ByPlayer[beto])

	pepe := round.NextTurn()
	// take: 10 hand + (5) board = 15, escobita
	boardCards, _ = match.Cards.Board.GetMultiple(4)
	takeCard(t, pepe, &match, match.Cards.ByPlayer[pepe].Hand[2], boardCards, true)
	t.Logf("%+v", match.Cards.ByPlayer[pepe])

	beto = round.NextTurn()
	// drop: 1
	dropCard(t, beto, &match, match.Cards.ByPlayer[beto].Hand[0])

	pepe = round.NextTurn()
	// drop: 10
	dropCard(t, pepe, &match, match.Cards.ByPlayer[pepe].Hand[0])

	beto = round.NextTurn()
	// drop: 3
	dropCard(t, beto, &match, match.Cards.ByPlayer[beto].Hand[0])

	pepe = round.NextTurn()
	dropCard(t, pepe, &match, match.Cards.ByPlayer[pepe].Hand[0])

	if match.HasMoreRounds() {
		t.Errorf("Match is one round only")
	}
	pepeCardsCountBeforeEnds := len(match.Cards.ByPlayer[pepe].Taken)
	match.Ends()
	// pepe takes 4 cards from the table including the seven gold!
	pepeCardsCountAfterEnds := len(match.Cards.ByPlayer[pepe].Taken)
	pepeCardsTakenAtEnds := pepeCardsCountAfterEnds - pepeCardsCountBeforeEnds
	if pepeCardsTakenAtEnds != 4 {
		t.Errorf("Pepe shall get the remaining 4 cards in the table and he gets %d cards ", pepeCardsTakenAtEnds)
	}
	t.Logf("%+v", match.ActionsByPlayer[pepe])

	staticticsByPlayer := match.CalculateStaticticsByPlayer()
	scoreSummaryByPlayer := staticticsByPlayer.BuildScoreSummaryByPlayer()
	//t.Logf("scoreSummaryByPlayer %+v\n", scoreSummaryByPlayer)

	/*if !staticticsByPlayer[beto].HasGoldSeven {
		t.Errorf("beto shall have gold seven!")
	}
	if staticticsByPlayer[beto].EscobitasCount != 0 {
		t.Errorf("beto shall have no escobita")
	}
	if staticticsByPlayer[pepe].EscobitasCount != 1 {
		t.Errorf("pepe shall have 1 escobita")
	}
	if *staticticsByPlayer.calculateMostGoldCardsPlayer() != beto {
		t.Errorf("beto shall be the player with most gold cards")
	}*/

	if scoreSummaryByPlayer[beto].Score != 0 {
		t.Errorf("beto shall have score 0 and pepe computed is %d", scoreSummaryByPlayer[beto].Score)
	}
	if scoreSummaryByPlayer[pepe].Score != 5 {
		t.Errorf("pepe shall have score 5 and pepe computed is %d", scoreSummaryByPlayer[pepe].Score)
	}
}

func takeCard(t *testing.T, player Player, match *Match, handCard Card, boardCards []Card, mustBeEscobita bool) (PlayerAction, error) {
	if CanTakeCards(handCard, boardCards) {
		takeAction := NewPlayerTakeAction(player, handCard, boardCards)
		action, err := match.Take(takeAction)
		if err != nil {
			return nil, err
		}
		if mustBeEscobita && !action.IsEscobita() {
			t.Fatalf("player take action shall be escobita and actual result is not")
		}
		return action, nil
	} else {
		errMsg := fmt.Sprintf("Invalid take action handCard=%v, boardCards=%v\n", handCard, boardCards)
		return nil, errors.New(errMsg)
	}
}

func dropCard(t *testing.T, player Player, match *Match, card Card) (PlayerAction, error) {
	dropAction := NewPlayerDropAction(player, card)
	return match.Drop(dropAction)
}
