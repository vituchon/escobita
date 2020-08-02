package model

import (
	"math/rand"
)

var EscobitaRanks []Rank = aggregateRanks(Ranks[:7], Ranks[9:])

// creates the match and prepare it for play
// do note that the initial cards are laydown at moment 0 and not at round one!
func CreateAndServe(players []Player) Match {
	var deck Deck = NewDeck(Suits, EscobitaRanks)
	match := newMatch(players, deck)
	match.serve()
	return match
}

func newMatch(players []Player, deck Deck) Match {
	match := Match{
		Players:          players,
		ActionsByPlayer:  newActionsByPlayer(players),
		MatchCards:       newMatchCards(players, deck),
		RoundNumber:      0,
		FirstPlayerIndex: 0,
	}
	return match
}

func (match *Match) serve() {
	shuffle(match.MatchCards.Left)
	match.MatchCards.Board = copyDeck(match.MatchCards.Left[:4])
	match.MatchCards.Left = match.MatchCards.Left[4:]
	match.FirstPlayerIndex = rand.Intn(len(match.Players))
}

func newActionsByPlayer(players []Player) ActionsByPlayer {
	actionsByPlayer := make(ActionsByPlayer)
	for _, player := range players {
		actionsByPlayer[player] = make([]PlayerAction, 0, 10)
	}
	return actionsByPlayer
}

func newMatchCards(players []Player, deck Deck) MatchCards {
	matchCards := MatchCards{
		Board:     nil,
		Left:      deck,
		PerPlayer: make(map[Player]PlayerMatchCards),
	}
	for _, player := range players {
		matchCards.PerPlayer[player] = PlayerMatchCards{
			Taken: nil,
			Hand:  nil,
		}
	}
	return matchCards
}

// Deal cards to each player for starting a new round
func (match *Match) NextRound() Round {
	for _, player := range match.Players {
		matchPlayerCards := match.MatchCards.PerPlayer[player]
		matchPlayerCards.Hand = copyDeck(match.MatchCards.Left[:3])
		match.MatchCards.PerPlayer[player] = matchPlayerCards
		/*fmt.Printf("\nmatchPlayerCards.Hand%+v\n", matchPlayerCards.Hand)
		fmt.Printf("\nmatch.MatchCards.Left%+v\n", match.MatchCards.Left)*/
		match.MatchCards.Left = match.MatchCards.Left[3:]
	}
	match.RoundNumber++
	return Round{
		Match:              match,
		CurrentPlayerIndex: match.FirstPlayerIndex,
		ConsumedTurns:      0,
		Number:             match.RoundNumber,
	}
}

func (match Match) HasMoreRounds() bool {
	cardsLeft := len(match.MatchCards.Left)
	playersCount := len(match.Players)
	return (cardsLeft/playersCount >= 3)
}

func CanTakeCards(handCard Card, boardCards []Card) bool {
	return sumValues(append(boardCards, handCard)) == 15
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

func (match *Match) Take(player Player, action PlayerTakeAction) PlayerAction {
	/*match.MatchCards.Board = */ match.MatchCards.Board.Without(action.BoardCards...)
	matchPlayerCards := match.MatchCards.PerPlayer[player]
	/*matchPlayerCards.Hand = */ matchPlayerCards.Hand.Without(action.HandCard)
	matchPlayerCards.Taken = append(matchPlayerCards.Taken, action.HandCard)
	matchPlayerCards.Taken = append(matchPlayerCards.Taken, action.BoardCards...)
	match.MatchCards.PerPlayer[player] = matchPlayerCards
	isEscobita := (len(match.MatchCards.Board) == 0)
	action.isEscobita = isEscobita
	match.ActionsByPlayer[player] = append(match.ActionsByPlayer[player], action)
	return action
}

func (match *Match) Drop(player Player, action PlayerDropAction) PlayerAction {
	match.MatchCards.Board = append(match.MatchCards.Board, action.HandCard)
	matchPlayerCards := match.MatchCards.PerPlayer[player]
	/*matchPlayerCards.Hand = */ matchPlayerCards.Hand.Without(action.HandCard)
	match.MatchCards.PerPlayer[player] = matchPlayerCards
	match.ActionsByPlayer[player] = append(match.ActionsByPlayer[player], action)
	return action
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
	Number             int
}

func (r Round) HasNextTurn() bool {
	return r.doHasNextTurnMethod2()
}

// this is slower than above but will fit for every quantity of players
func (r Round) doHasNextTurnMethod2() bool {
	for _, player := range r.Match.Players {
		if len(r.Match.MatchCards.PerPlayer[player].Hand) > 0 {
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
	isEscobita bool
}

func (a PlayerTakeAction) IsEscobita() bool {
	return a.isEscobita
}

type PlayerDropAction struct {
	HandCard Card
}

func (a PlayerDropAction) IsEscobita() bool {
	return false
}

type PlayerAction interface {
	IsEscobita() bool
}
