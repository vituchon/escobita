package model

import (
	"reflect"
	"testing"
)

var cards []Card = []Card{Card{1, GOLD, 1}, Card{2, SWORD, 1}, Card{3, CLUB, 1}, Card{4, CUP, 1}}

func BenchmarkDetermineIsGoldenSevenIsUsed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DetermineIsGoldenSevenIsUsed(cards)
	}
}

func TestCalculatePossibleTakeActions(t *testing.T) {
	testRuns := []struct {
		boardCards          []Card
		handCards           []Card
		expectedSuggestions []PlayerTakeAction
	}{
		{
			boardCards: []Card{
				Card{Id: 1, Suit: GOLD, Rank: 1},
				Card{Id: 2, Suit: GOLD, Rank: 2},
				Card{Id: 3, Suit: CLUB, Rank: 2},
				Card{Id: 4, Suit: SWORD, Rank: 5},
			},
			handCards: []Card{Card{Id: 5, Suit: CLUB, Rank: 12}},
			expectedSuggestions: []PlayerTakeAction{
				PlayerTakeAction{
					basePlayerAction: basePlayerAction{
						Player: BotPlayer,
					},
					BoardCards: []Card{
						Card{Id: 1, Suit: GOLD, Rank: 1},
						Card{Id: 2, Suit: GOLD, Rank: 2},
						Card{Id: 3, Suit: CLUB, Rank: 2},
					},
					HandCard:    Card{Id: 5, Suit: CLUB, Rank: 12},
					Is_Escobita: false,
				},
				PlayerTakeAction{
					basePlayerAction: basePlayerAction{
						Player: BotPlayer,
					},
					BoardCards: []Card{
						Card{Id: 4, Suit: SWORD, Rank: 5},
					},
					HandCard:    Card{Id: 5, Suit: CLUB, Rank: 12},
					Is_Escobita: false,
				},
			},
		},
	}

	for _, testRun := range testRuns {
		generatedSuggestions := CalculatePossibleTakeActions(testRun.boardCards, testRun.handCards)
		t.Log(generatedSuggestions)
		areEquals := reflect.DeepEqual(testRun.expectedSuggestions, generatedSuggestions)
		t.Log(areEquals)
		if !areEquals {
			t.Errorf("Generated was:\n%+v\nExpected is:\n%+v\n", generatedSuggestions, testRun.expectedSuggestions)
		}
	}
}
