package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Match struct {
	Players          []Player        `json:"players"`
	ActionsByPlayer  ActionsByPlayer `json:"actionsByPlayer"`
	ActionsLog       PlayerActions   `json:"playerActions"`
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
		ActionsLog:       make(PlayerActions, 0, totalTurns),
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

type ActionsByPlayer map[Player]PlayerActions

type PlayerActions []PlayerAction

// Dev notes: Marshalling of interfaces' array or map whose values are interfaces requires "special treatment" request.Body. See https://stackoverflow.com/q/52783848/903998 and https://stackoverflow.com/a/42765078
func (actions *PlayerActions) UnmarshalJSON(b []byte) error {
	var rawActions []map[string]interface{}
	if err := json.Unmarshal(b, &rawActions); err != nil {
		return err
	}

	var parsedActions []PlayerAction
	for _, rawAction := range rawActions {
		_, hasBoardCardsField := rawAction["boardCards"] // if it contains board cards then is take action, else is a drop action. Recall that there aren't more actions.
		if hasBoardCardsField {
			var takeAction PlayerTakeAction = parsePlayerTakeAction(rawAction)
			parsedActions = append(parsedActions, takeAction)
		} else {
			var dropAction PlayerDropAction = parsePlayerDropAction(rawAction)
			parsedActions = append(parsedActions, dropAction)
		}
	}

	*actions = parsedActions
	return nil
}

func parsePlayerDropAction(m map[string]interface{}) PlayerDropAction {
	playerMap := (m["player"]).(map[string]interface{})
	handCardMap := (m["handCard"]).(map[string]interface{})
	return PlayerDropAction{
		basePlayerAction: basePlayerAction{
			Player: parsePlayer(playerMap),
		},
		HandCard: parseSingleCard(handCardMap),
	}
}

func parsePlayerTakeAction(m map[string]interface{}) PlayerTakeAction {
	playerMap := (m["player"]).(map[string]interface{})
	handCardMap := (m["handCard"]).(map[string]interface{})
	rawBoardCards := (m["boardCards"]).([]interface{})

	return PlayerTakeAction{
		basePlayerAction: basePlayerAction{
			Player: parsePlayer(playerMap),
		},
		Is_Escobita: m["isEscobita"].(bool),
		HandCard:    parseSingleCard(handCardMap),
		BoardCards:  parseMultipleCards(rawBoardCards),
	}
}

func parseSingleCard(m map[string]interface{}) Card {
	return Card{
		Id:   int(m["id"].(float64)),
		Rank: Rank(m["rank"].(float64)),
		Suit: Suit(m["suit"].(float64)),
	}
}

func parseMultipleCards(rawCards []interface{}) []Card {
	var boardCards []Card = make([]Card, len(rawCards), len(rawCards))
	for i, rawCard := range rawCards {
		m := rawCard.(map[string]interface{})
		boardCard := parseSingleCard(m)
		boardCards[i] = boardCard
	}
	return boardCards
}

func parsePlayer(m map[string]interface{}) Player {
	return Player{
		Name: m["name"].(string),
	}
}

func newActionsByPlayer(players []Player) ActionsByPlayer {
	actionsByPlayer := make(ActionsByPlayer)
	for _, player := range players {
		actionsByPlayer[player] = make(PlayerActions, 0, 10)
	}
	return actionsByPlayer
}

type MatchCards struct {
	Board    Deck                        `json:"board"` // the cards on the table that anyone can reclaim
	Left     Deck                        `json:"left"`  // the remaining cards to play in the rest of the match
	ByPlayer map[Player]PlayerMatchCards `json:"byPlayer"`
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
	Taken Deck `json:"taken" // the cards on the player has claimed `
	Hand  Deck `json:"hand"` // the cards on the player has to play
}
