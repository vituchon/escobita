package model

import (
	"encoding/json"
)

type Player struct {
	Name string `json:"name"`
}

type Party []Player

func (player Player) String() string {
	return player.Name
}

// Dev notes: Implementing encoding.TextMarshaler for be complaint with Marshal/unMarshal on using a Player as a map's key
func (player Player) MarshalText() (text []byte, err error) {
	return []byte(player.Name), nil
}

func (player *Player) UnmarshalText(text []byte) error {
	player.Name = string(text)
	return nil
}

func (p Player) MarshalJSON() ([]byte, error) {
	return []byte(`{"name":"` + p.Name + `"}`), nil
}

func (p *Player) UnmarshalJSON(b []byte) error {
	var stuff map[string]interface{}
	err := json.Unmarshal(b, &stuff)
	// dev notes: on Match#ActionsByPlayer the marshalling generates a key with the name (an string, refer to Player#MarshalText), so when unmarshalling I need to take that string as the name
	if err != nil {
		p.Name = string(b)
	} else {
		p.Name = stuff["name"].(string)
	}
	return nil
}
