package model

import (
	"errors"
	"strconv"
	"strings"
)

// Spanish card
type Card struct {
	Id   int
	Suit Suit
	Rank Rank
}

func (s Card) String() string {
	return "[" + strconv.Itoa(s.Id) + ", " + strconv.Itoa(determineValue(s)) + "] " + s.Suit.String() + "," + strconv.Itoa(s.Rank)
}

// The suit that a card belongs to
type Suit int

const (
	SWORD Suit = iota
	CLUB
	CUP
	GOLD
)

var Suits []Suit = []Suit{SWORD, CLUB, CUP, GOLD}

var suitCodenames = []string{
	"sword",
	"club",
	"cup",
	"gold",
}

func (s Suit) String() string {
	index := int(s)
	if index < 0 || index > len(suitCodenames) {
		return "??"
	}
	return suitCodenames[index]
}

// The rank or number of the card
type Rank = int

var Ranks []Rank = []Rank{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

// It is collection of cards
type Deck []Card

var NoCardFoundErr = errors.New("There is no card with the given id")
var NoCard = Card{}

func (d Deck) GetSingle(cardId int) (Card, error) {
	for _, card := range d {
		if card.Id == cardId {
			return card, nil
		}
	}
	return NoCard, NoCardFoundErr
}

func (d Deck) GetMultiple(cardIds ...int) ([]Card, error) {
	cards := []Card{}
	for _, cardId := range cardIds {
		card, err := d.GetSingle(cardId)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (d Deck) Without(cards ...Card) Deck {
	var filtered Deck = []Card{}
	for _, deckCard := range d {
		discard := false
		for _, card := range cards {
			if deckCard.Id == card.Id {
				discard = true
				break
			}
		}
		if !discard {
			filtered = append(filtered, deckCard)
		}
	}
	return filtered
}

func (d Deck) String() string {
	cardStrings := make([]string, 0, len(d))
	for _, card := range d {
		cardStrings = append(cardStrings, card.String())
	}
	return strings.Join(cardStrings, "|")
}

func NewDeck(suits []Suit, ranks []Rank) Deck {
	deck := make([]Card, 0, len(suits)*len(ranks))
	i := 1
	for _, suit := range suits {
		for _, rank := range ranks {
			card := Card{
				Suit: suit,
				Rank: rank,
				Id:   i,
			}
			deck = append(deck, card)
			i++
		}
	}
	return deck
}
