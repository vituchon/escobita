package model

import (
	"fmt"
	"math/rand"
)

var EscobitaRanks []Rank = aggregateRanks(Ranks[:7], Ranks[9:])

// creates the match and prepare it for play
// do note that the initial cards are laydown at moment 0 and not at round one!
func CreateAndBegins(players []Player) Match {
	match := CreateMatch(players)
	match.Begins()
	return match
}

func CreateMatch(players []Player) Match {
	var deck Deck = NewDeck(Suits, EscobitaRanks)
	return newMatch(players, deck)
}

// performs the initial lay down of cards and select the initial player so the match is ready to begin
func (match *Match) Begins() {
	shuffle(match.Cards.Left)
	match.Cards.Board = copyDeck(match.Cards.Left[:4])
	match.Cards.Left = match.Cards.Left[4:]
	match.FirstPlayerIndex = rand.Intn(len(match.Players))
}

// finalizes the match grating the left cards to the last taker
func (match *Match) Ends() {
	if len(match.Cards.Board) > 0 {
		player := match.getLastCardTaker()
		if player != nil {
			matchPlayerCards := match.Cards.PerPlayer[*player]
			// dev notes: a design decision here may be to track this "movement" as an other type of Action, like CleanBoardAction
			matchPlayerCards.Taken = append(matchPlayerCards.Taken, match.Cards.Board...)
			match.Cards.PerPlayer[*player] = matchPlayerCards
			match.Cards.Board = match.Cards.Board[:0] // practical way to empty an slice
			fmt.Printf("The last card taker is %v", *player)
		} else {
			fmt.Println("Nobody takes cards")
		}
	}
}

func (match Match) getLastCardTaker() *Player {
	for i := len(match.ActionsLog) - 1; i > 0; i-- {
		action := match.ActionsLog[i]
		_, isTakeAction := action.(PlayerTakeAction)
		if isTakeAction {
			player := action.GetPlayer()
			return &player
		}
	}
	return nil
}

// Performs cards take from board using a hand card.
// It is assumed that the combination of cards is valid (sums 15)
func (match *Match) Take(action PlayerTakeAction) PlayerAction {
	player := action.Player
	match.Cards.Board.Without(action.BoardCards...)
	matchPlayerCards := match.Cards.PerPlayer[player]
	matchPlayerCards.Hand.Without(action.HandCard)
	matchPlayerCards.Taken = append(matchPlayerCards.Taken, action.HandCard)
	matchPlayerCards.Taken = append(matchPlayerCards.Taken, action.BoardCards...)
	match.Cards.PerPlayer[player] = matchPlayerCards
	isEscobita := (len(match.Cards.Board) == 0)
	action.isEscobita = isEscobita
	match.ActionsByPlayer[player] = append(match.ActionsByPlayer[player], action)
	match.ActionsLog = append(match.ActionsLog, action)
	return action
}

// Performs a card drop
func (match *Match) Drop(action PlayerDropAction) PlayerAction {
	player := action.Player
	match.Cards.Board = append(match.Cards.Board, action.HandCard)
	matchPlayerCards := match.Cards.PerPlayer[player]
	matchPlayerCards.Hand.Without(action.HandCard)
	match.Cards.PerPlayer[player] = matchPlayerCards
	match.ActionsByPlayer[player] = append(match.ActionsByPlayer[player], action)
	match.ActionsLog = append(match.ActionsLog, action)
	return action
}

// Deal cards to each player for starting a new round
func (match *Match) NextRound() Round {
	for _, player := range match.Players {
		matchPlayerCards := match.Cards.PerPlayer[player]
		matchPlayerCards.Hand = copyDeck(match.Cards.Left[:3])
		match.Cards.PerPlayer[player] = matchPlayerCards
		/*fmt.Printf("\nmatchPlayerCards.Hand%+v\n", matchPlayerCards.Hand)
		fmt.Printf("\nmatch.MatchCards.Left%+v\n", match.MatchCards.Left)*/
		match.Cards.Left = match.Cards.Left[3:]
	}
	match.RoundNumber++
	round := Round{
		Match:              match,
		currentPlayerIndex: match.FirstPlayerIndex,
		ConsumedTurns:      0,
		Number:             match.RoundNumber,
	}
	match.CurrentRound = &round
	return round
}

func (match Match) HasMoreRounds() bool {
	cardsLeft := len(match.Cards.Left)
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

func shuffle(deck Deck) {
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

type Round struct {
	Match              *Match  `json:"match"`
	CurrentTurnPlayer  *Player `json:"currentTurnPlayer"`
	currentPlayerIndex int     `json:"-"`
	ConsumedTurns      int     `json:"consumedTurns"`
	Number             int     `json:"number"`
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
	nextPlayer := party[r.currentPlayerIndex%len(party)]
	r.CurrentTurnPlayer = &nextPlayer
	r.currentPlayerIndex++
	r.ConsumedTurns++
	return nextPlayer
}

type basePlayerAction struct {
	Player Player `json:"player"` // the performer
}

func (bpa basePlayerAction) GetPlayer() Player {
	return bpa.Player
}

type PlayerTakeAction struct {
	basePlayerAction
	BoardCards []Card `json:"boardCards"`
	HandCard   Card   `json:"handCard"`
	isEscobita bool   `json:"isEscobita"`
}

func NewPlayerTakeAction(player Player, handCard Card, boardCards []Card) PlayerTakeAction {
	return PlayerTakeAction{
		basePlayerAction: basePlayerAction{
			Player: player,
		},
		HandCard:   handCard,
		BoardCards: boardCards,
	}
}

func (a PlayerTakeAction) IsEscobita() bool {
	return a.isEscobita
}

type PlayerDropAction struct {
	basePlayerAction
	HandCard Card `json:"handCard"`
}

func NewPlayerDropAction(player Player, handCard Card) PlayerDropAction {
	return PlayerDropAction{
		basePlayerAction: basePlayerAction{
			Player: player,
		},
		HandCard: handCard,
	}
}

func (a PlayerDropAction) IsEscobita() bool {
	return false
}

type PlayerAction interface {
	IsEscobita() bool
	GetPlayer() Player
}
