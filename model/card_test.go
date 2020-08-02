package model

import (
	"reflect"
	"testing"
)

func TestDeckWithout(t *testing.T) {
	var deck Deck = Deck{Card{1, GOLD, 1}, Card{2, SWORD, 1}, Card{3, CLUB, 1}, Card{4, CUP, 1}}
	testRuns := []struct {
		title         string
		cards         Deck
		cardsToRemove Deck
		cardsExpected Deck
	}{
		{
			title:         "1",
			cards:         copyDeck(deck),
			cardsToRemove: deck[0:1],
			cardsExpected: deck[1:],
		},
		{
			title:         "2",
			cards:         copyDeck(deck),
			cardsToRemove: aggregateDecks(deck[0:1], Deck{deck[3]}),
			cardsExpected: deck[1:3],
		},
		{
			title:         "3",
			cards:         copyDeck(deck),
			cardsToRemove: aggregateDecks(deck[3:4], Deck{deck[1]}),
			cardsExpected: aggregateDecks(deck[0:1], Deck{deck[2]}),
		},
	}
	for _, testRun := range testRuns {
		t.Logf("\n=====Running unit test: %s=====\n", testRun.title)
		t.Logf("\nCards=%v\nToRemove=%v\n", testRun.cards, testRun.cardsToRemove)
		testRun.cards.Without(testRun.cardsToRemove...)
		cardsRemaining := testRun.cards
		t.Logf("\nRemaining=%v\n", cardsRemaining)
		if !areDeckEquals(cardsRemaining, testRun.cardsExpected) {
			t.Errorf("Cards remaining and cards expected differ!\nRemaining=%v\nExpected=%v", cardsRemaining, testRun.cardsExpected)
		}
	}
}

func areDeckEquals(generated Deck, expected Deck) bool {
	if len(generated) == 0 && len(expected) == 0 {
		return true
	}
	return reflect.DeepEqual(generated, expected)
}
