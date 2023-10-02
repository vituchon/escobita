package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Player struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

var playerFieldSeparator = "|"
var InvalidPlayerNameErr = errors.New("Player name is invalid as it containst a reserved char (" + playerFieldSeparator + ")")

func ValidateName(name string) error {
	if strings.Contains(name, playerFieldSeparator) {
		return InvalidPlayerNameErr
	}
	return nil
}

func (player Player) String() string {
	return strconv.Itoa(player.Id) + playerFieldSeparator + player.Name
}

// Dev notes: Implementing encoding.TextMarshaler for be complaint with Marshal/unMarshal on using a Player as a map's key
func (player Player) MarshalText() (text []byte, err error) {
	return []byte(player.String()), nil
}

func (player *Player) UnmarshalText(text []byte) error {
	parts := strings.Split(string(text), playerFieldSeparator)
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		err := fmt.Sprintf("Unexpected error on UnmarshalText, error was: '%v'", err)
		return errors.New(err)
	}
	player.Id = id
	player.Name = string(parts[1])
	return nil
}

func (player Player) MarshalJSON() ([]byte, error) {
	return []byte(`{"name":"` + player.Name + `", "id":` + strconv.Itoa(player.Id) + `}`), nil
}

func (player *Player) UnmarshalJSON(bytes []byte) error {
	if strings.Contains(string(bytes), playerFieldSeparator) {
		return player.UnmarshalText(bytes[1 : len(bytes)-1]) //  removing leading (beginning) and trailing (end) quotes that are present in a json string
	}
	var stuff map[string]interface{}
	err := json.Unmarshal(bytes, &stuff)
	if err != nil {
		return err
	}
	player.Name = stuff["name"].(string)
	player.Id = int(stuff["id"].(float64))
	return nil
}

// Dev notes: An idea for naming
type Party []Player
