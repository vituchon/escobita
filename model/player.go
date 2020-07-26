package model

type Player struct {
	Name string
}

type Party []Player // just an idea

func (player Player) String() string {
	return player.Name
}

type PlayerTakeAction struct {
	BoardCards []Card
	HandCard   Card
	isEscobita bool
}

var boolToInt map[bool]int = map[bool]int{
	true:  1,
	false: 0,
}

func (a PlayerTakeAction) Score() int {
	return boolToInt[a.isEscobita]
}

type PlayerDropAction struct {
	HandCard Card
}

func (a PlayerDropAction) Score() int {
	return 0
}

type PlayerAction interface {
	Score() int
}
