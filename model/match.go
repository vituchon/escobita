package model

import (
	"fmt"
	"strings"
)

type Match struct {
	Players          []Player
	FirstPlayerIndex int
	ScorePerPlayer   map[Player]int
	Cards            MatchCards
	RoundNumber      int

	//Status         string
}

type MatchPlayer struct {
	Player  Player
	Actions []PlayerAction
}

/*
type matchStatus struct {
	Served  string
	OnGoing string
	Finish  string
}

var MatchStatus = matchStatus{
	Served:  "served",   // ready to play
	OnGoing: "on-going", // playing
	Finish:  "finish",
}*/

type MatchCards struct {
	Board     Deck // the cards on the table that anyone can reclaim
	Left      Deck // the remaining cards to play in the rest of the match
	PerPlayer map[Player]PlayerMatchCards
}

type PlayerMatchCards struct {
	Taken Deck // the cards on the player has claimed
	Hand  Deck // the cards on the player has to play
}

func (m Match) String() string {
	playersDescription := make([]string, 0, len(m.Players))
	for _, player := range m.Players {
		cardsInHands := Deck(m.Cards.PerPlayer[player].Hand).String()
		cardsTaken := Deck(m.Cards.PerPlayer[player].Taken).String()
		playerDescription := fmt.Sprintf("%s, score: %d\nCards taken:%v\nCards in hand:%v", player.Name, m.ScorePerPlayer[player], cardsTaken, cardsInHands)
		playersDescription = append(playersDescription, playerDescription)
	}
	joinedPlayersDescription := strings.Join(playersDescription, "\n")
	matchBoardCards := Deck(m.Cards.Board).String()
	matchLeftCards := Deck(m.Cards.Left).String()
	return fmt.Sprintf("Match\nLeft cards:%v\nBoard cards: %v\nPlayers:\n%v", matchLeftCards, matchBoardCards, joinedPlayersDescription)
}
