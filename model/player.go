package model

type Player struct {
	Name string
}

type Party []Player // just an idea

func (player Player) String() string {
	return player.Name
}
