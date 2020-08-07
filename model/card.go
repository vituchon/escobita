package model

import (
	"errors"
	"fmt"
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
	//return "(id=" + strconv.Itoa(s.Id) + ",value=" + strconv.Itoa(determineValue(s)) + ") " + s.Suit.String() + "," + strconv.Itoa(s.Rank)
	return s.Suit.String() + "," + strconv.Itoa(s.Rank)
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

func (d *Deck) Without(cards ...Card) {
	f := (*d)[:0]
	for _, dc := range *d {
		include := true
		for _, c := range cards {
			if dc.Id == c.Id {
				include = false
				break
			}
		}
		if include {
			f = append(f, dc)
		}
	}
	(*d) = f
}

func (d Deck) SplitByRank() map[Suit]Deck {
	var splitted map[Suit]Deck = make(map[Suit]Deck)
	for _, card := range d {
		splitted[card.Suit] = append(splitted[card.Suit], card)
	}
	return splitted
}

// Gets the card that is most closer (from the "left": lower or equals) if exists, if not a nil is returned
func (d Deck) getLeftCloserToRank(rank int) *Card {
	higherRank := 0
	var higherCard *Card
	for _, card := range d {
		fmt.Printf("Comparando carta %v contra %d\n ", card, higherRank)
		if card.Rank <= rank && card.Rank > higherRank {
			cardCopy := card
			higherCard = &cardCopy
			higherRank = cardCopy.Rank
		}
	}
	return higherCard
}

func (d Deck) String() string {
	cardStrings := make([]string, 0, len(d))
	for _, card := range d {
		cardStrings = append(cardStrings, card.String())
	}
	return strings.Join(cardStrings, ";")
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

func copyDeck(original Deck) (replica Deck) {
	replica = make(Deck, len(original), len(original))
	copy(replica, original)
	return
}

func aggregateDecks(d1, d2 Deck) Deck {
	lenD1 := len(d1)
	deck := make(Deck, lenD1, lenD1+len(d2))
	_ = copy(deck, d1)
	deck = append(deck, d2...)
	return deck
}

func aggregateRanks(r1, r2 []Rank) []Rank {
	lenR1 := len(r1)
	ranks := make([]Rank, lenR1, lenR1+len(r2))
	_ = copy(ranks, r1)
	ranks = append(ranks, r2...)
	return ranks
}
