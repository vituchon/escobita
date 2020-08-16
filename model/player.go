package model

type Player struct {
	Name string `json:"name"`
}

type Party []Player // just an idea

func (player Player) String() string {
	return player.Name
}

// implementing encoding.TextMarshaler for be complaint with Marshal/unMarshal when key of map is a struct
func (player Player) MarshalText() (text []byte, err error) {
	return []byte(player.Name), nil
}

func (player *Player) UnmarshalText(text []byte) error {
	player.Name = string(text)
	return nil
}
