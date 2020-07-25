package model

import (
	"math/rand"
)

var EscobitaRanks []Rank = append(Ranks[:7], Ranks[9:]...)

// creates the match and prepare it for play (do note that the initial cards are laydown in this step, and not in the first round)
func CreateAndServe(players []Player) Match {
	match := newMatch(players)
	shuffle(match.Cards.Left)
	match.Cards.Board = match.Cards.Left[:4]
	match.Cards.Left = match.Cards.Left[4:]
	match.FirstPlayerIndex = rand.Intn(len(match.Players))
	return match
}

func newMatch(players []Player) Match {
	var deck Deck = NewDeck(Suits, EscobitaRanks)
	match := Match{
		Players:        players,
		ScorePerPlayer: make(map[Player]int),
		Cards:          newMatchCards(players, deck),
	}
	for _, player := range players {
		match.ScorePerPlayer[player] = 0
	}
	return match
}

func newMatchCards(players []Player, deck Deck) MatchCards {
	matchCards := MatchCards{
		Board:     make([]Card, 0, 4),
		Left:      deck,
		PerPlayer: make(map[Player]PlayerMatchCards),
	}
	for _, player := range players {
		matchCards.PerPlayer[player] = PlayerMatchCards{
			Taken: make([]Card, 0, 10),
			Hand:  make([]Card, 0, 3),
		}
	}
	return matchCards
}

// Deal cards to each player for starting a new round
func (match *Match) NextRound() Round {
	for _, player := range match.Players {
		hand := match.Cards.Left[:3]
		matchPlayerCards := match.Cards.PerPlayer[player]
		matchPlayerCards.Hand = hand
		match.Cards.PerPlayer[player] = matchPlayerCards
		match.Cards.Left = match.Cards.Left[3:]
	}
	return Round{
		Match:              match,
		CurrentPlayerIndex: match.FirstPlayerIndex,
		ConsumedTurns:      0,
	}
}

func (match Match) MatchCanHaveMoreRounds() bool {
	cardsLeft := len(match.Cards.Left)
	playersCount := len(match.Players)
	return (cardsLeft/playersCount >= 3)
}

func CanTakeCards(handCard Card, boardCards []Card) bool {
	return sumValues(append(boardCards, handCard)) == 15
}

func (match *Match) Take(player Player, action PlayerTakeAction) {
	match.Cards.Board = match.Cards.Board.Without(action.BoardCards...)
	matchPlayerCards := match.Cards.PerPlayer[player]
	matchPlayerCards.Hand = matchPlayerCards.Hand.Without(action.HandCard)
	matchPlayerCards.Taken = append(matchPlayerCards.Taken, action.HandCard)
	matchPlayerCards.Taken = append(matchPlayerCards.Taken, action.BoardCards...)
	match.Cards.PerPlayer[player] = matchPlayerCards
}

func (match *Match) Drop(player Player, card Card) {
	match.Cards.Board = append(match.Cards.Board, card)
	matchPlayerCards := match.Cards.PerPlayer[player]
	matchPlayerCards.Hand = matchPlayerCards.Hand.Without(card)
	match.Cards.PerPlayer[player] = matchPlayerCards
}

func sumValues(cards []Card) int {
	total := 0
	for _, card := range cards {
		total += determineValue(card)
	}
	return total
}

func determineValue(card Card) int {
	if card.Rank < 8 {
		return card.Rank
	} else {
		return card.Rank - 2
	}
}

func shuffle(deck Deck) {
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

type Round struct {
	Match              *Match
	CurrentPlayerIndex int
	ConsumedTurns      int
}

func (r Round) HasNextTurn() bool {
	return r.doHasNextTurnMethod2()
}

// this is slower than above but will fit for every quantity of players
func (r Round) doHasNextTurnMethod2() bool {
	for _, player := range r.Match.Players {
		if len(r.Match.Cards.PerPlayer[player].Hand) > 0 {
			return true
		}
	}
	return false
}

// this is faster but won't work for matchs where "36 % len(r.Match.Players) > 0"
// so to use both an state pattern or somelike that (set on initialization time) would be required,a nice to do thing
func (r Round) doHasNextTurnMethod1() bool {
	return r.ConsumedTurns < len(r.Match.Players)*3
}

func (r *Round) NextTurn() Player {
	party := r.Match.Players
	nextPlayer := party[r.CurrentPlayerIndex%len(party)]
	r.CurrentPlayerIndex++
	r.ConsumedTurns++
	return nextPlayer
}

type PlayerTakeAction struct {
	BoardCards []Card
	HandCard   Card
}

type PlayerDropAction struct {
	HandCard Card
}
