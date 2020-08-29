package model

import (
	"fmt"
	"strings"
)

type Match struct {
	Players          []Player        `json:"players"`
	ActionsByPlayer  ActionsByPlayer `json:"actionsByPlayerName"`
	ActionsLog       []PlayerAction  `json:"playerAction"`
	Cards            MatchCards      `json:"matchCards"`
	FirstPlayerIndex int             `json:"firstPlayerIndex"`
	RoundNumber      int             `json:"roundNumber"`
	// dev notes: este campo se agrego acá, pero tranquilamente podria haber sido agregado a nivel WebGame dado que es allí donde aparece la
	// necesidad de acceder a la ronda de juego en progreso, lo cual tambien se podria haber programado a nivel WebGame guardando en dicha struct el round
	// que devuelve el método NextRound
	CurrentRound *Round `json:"currentRound,omitempty"` // TODO : add some test to ensure this is tracking OK
}

func newMatch(players []Player, deck Deck) Match {
	totalTurns := len(deck) - 4 // it should be 36, but in the tests is lower because i employ an small deck
	match := Match{
		Players:          players,
		ActionsByPlayer:  newActionsByPlayer(players),
		ActionsLog:       make([]PlayerAction, 0, totalTurns),
		Cards:            newMatchCards(players, deck),
		RoundNumber:      0,
		FirstPlayerIndex: 0,
		CurrentRound:     nil,
	}
	return match
}

func (m Match) String() string {
	playersDescription := make([]string, 0, len(m.Players))
	for _, player := range m.Players {
		cardsInHands := Deck(m.Cards.ByPlayer[player].Hand).String()
		cardsTaken := Deck(m.Cards.ByPlayer[player].Taken).String()
		playerDescription := fmt.Sprintf("%s\nCards taken:%v\nCards in hand:%v", player.Name, cardsTaken, cardsInHands)
		playersDescription = append(playersDescription, playerDescription)
	}
	joinedPlayersDescription := strings.Join(playersDescription, "\n")
	matchBoardCards := Deck(m.Cards.Board).String()
	matchLeftCards := Deck(m.Cards.Left).String()
	return fmt.Sprintf("Match, first player is %v and current round is %v,\nLeft cards:%v\nBoard cards: %v\nPlayers:\n%v", m.Players[m.FirstPlayerIndex], m.RoundNumber, matchLeftCards, matchBoardCards, joinedPlayersDescription)
}

type ActionsByPlayer map[Player][]PlayerAction

func newActionsByPlayer(players []Player) ActionsByPlayer {
	actionsByPlayer := make(ActionsByPlayer)
	for _, player := range players {
		actionsByPlayer[player] = make([]PlayerAction, 0, 10)
	}
	return actionsByPlayer
}

type MatchCards struct {
	Board    Deck                        `json:"board"` // the cards on the table that anyone can reclaim
	Left     Deck                        `json:"left"`  // the remaining cards to play in the rest of the match
	ByPlayer map[Player]PlayerMatchCards `json:"byPlayerName"`
}

func newMatchCards(players []Player, deck Deck) MatchCards {
	matchCards := MatchCards{
		Board:    nil,
		Left:     deck,
		ByPlayer: make(map[Player]PlayerMatchCards),
	}
	for _, player := range players {
		matchCards.ByPlayer[player] = PlayerMatchCards{
			Taken: nil,
			Hand:  nil,
		}
	}
	return matchCards
}

type PlayerMatchCards struct {
	Taken Deck // the cards on the player has claimed `json:"taken"`
	Hand  Deck // the cards on the player has to play `json:"hand"`
}
