package main

import (
	"bufio"
	"errors"
	"strconv"

	//"errors"
	"fmt"
	"local/escobita/model"
	"os"
)

const DefaultConfigFilePath = "default-config.json"

func main() {

	var players []model.Player = []model.Player{
		model.Player{Name: "Beto"},
		model.Player{Name: "Pepe"},
	}

	quitGame := false
	match := model.CreateAndServe(players)
	for match.MatchCanHaveMoreRounds() && !quitGame {
		round := match.NextRound()
		for round.HasNextTurn() && !quitGame {
			player := round.NextTurn()
			playerCards := match.Cards.PerPlayer[player]
			fmt.Println("Restantes: ", match.Cards.Left)

			playerMustAct := true
			for playerMustAct {
				fmt.Println("Juega ", player)
				fmt.Println("En mesa: ", match.Cards.Board)
				fmt.Println("En mano: ", playerCards.Hand)
				fmt.Println("Tomadas: ", playerCards.Taken)

				fmt.Print("Que deseas hacer, jugar (c), soltar carta (d), salir (q): ")
				cmd := doReadInput()
				if cmd == ClaimCommand {
					takeAction := readTakeActionFromStdin(player, match)
					fmt.Println(takeAction.HandCard)
					fmt.Println(takeAction.BoardCards)
					isValidClaim := model.CanTakeCards(takeAction.HandCard, takeAction.BoardCards)
					if isValidClaim {
						fmt.Println("La jugada es valida!")
						playerMustAct = false
						match.Take(player, takeAction)
					} else {
						fmt.Println("La jugada NO es valida")
					}
				} else if cmd == DropCommand {
					dropAction := readDropActionFromStdin(player, match)
					match.Drop(player, dropAction)
					playerMustAct = false
				} else if cmd == ExitCommand {
					playerMustAct = false
					quitGame = true
				}
			}

		}
		//fmt.Println(match)
	}
}

func readTakeActionFromStdin(player model.Player, match model.Match) model.PlayerTakeAction {
	fmt.Println("==" + player.Name + " selecciona combinación entre las de mesa y una de mano==")
	playerCards := match.Cards.PerPlayer[player]
	fmt.Print("La de mano, id de carta: ")
	cardId := ReadSingleIntInput()
	handCard, err := playerCards.Hand.GetSingle(cardId)
	for err != nil {
		fmt.Println("Error: ", err)
		fmt.Print("La de mano, id de carta: ")
		cardId := ReadSingleIntInput()
		handCard, err = playerCards.Hand.GetSingle(cardId)
	}

	fmt.Print("La de la mesa, ids de cartas (f para terminar ingreso): ")
	cardsIds := ReadMultipleIntsInput()
	boardCards, err := match.Cards.Board.GetMultiple(cardsIds...)
	for err != nil {
		fmt.Println("Error: ", err)
		fmt.Print("La de la mesa, ids de cartas (f para terminar ingreso): ")
		cardsIds := ReadMultipleIntsInput()
		boardCards, err = match.Cards.Board.GetMultiple(cardsIds...)
	}
	return model.PlayerTakeAction{
		BoardCards: boardCards,
		HandCard:   handCard,
	}
}

func readDropActionFromStdin(player model.Player, match model.Match) model.PlayerDropAction {
	fmt.Print("La de mano, id de carta: ")
	playerCards := match.Cards.PerPlayer[player]
	cardId := ReadSingleIntInput()
	handCard, err := playerCards.Hand.GetSingle(cardId)
	for err != nil {
		fmt.Println("Error: ", err)
		fmt.Print("La de mano, id de carta: ")
		cardId := ReadSingleIntInput()
		handCard, err = playerCards.Hand.GetSingle(cardId)
	}
	return model.PlayerDropAction{
		HandCard: handCard,
	}
}

func ReadMultipleIntsInput() (ints []int) {
	value, command := ReadInputAndParseAnInt()
	for command != nil && command != FinishCommand {
		if command == ExitCommand {
			os.Exit(0)
		}
		ints = append(ints, value)
		value, command = ReadInputAndParseAnInt()
	}
	return ints
}

func ReadSingleIntInput() int {
	value, command := ReadInputAndParseAnInt()
	for command != nil && command == FinishCommand {
		value, command = ReadInputAndParseAnInt()
	}
	if command == ExitCommand {
		os.Exit(0)
	}

	return value
}

func ReadInputAndParseAnInt() (int, Command) {
	value, err, command := doReadIntInput()
	isEnteringAnotherCommand := (command == ExitCommand || command == FinishCommand)
	for err != nil && !isEnteringAnotherCommand {
		value, err, command = doReadIntInput()
		isEnteringAnotherCommand = (command == ExitCommand || command == FinishCommand)
	}
	return value, command

}

var NotANumberErr = errors.New("Input value is not a number")

func doReadIntInput() (int, error, Command) {
	command := doReadInput()
	if command.GetType() == InputValueCommandType {
		value, atoiErr := strconv.Atoi(command.GetValue())
		return value, atoiErr, command
	} else {
		return 0, NotANumberErr, command
	}

}

type Command interface {
	GetType() CommandType
	GetValue() string
}

type CommandType int

const (
	ExitCommandType CommandType = iota
	FinishCommandType
	DropCommandType
	ClaimCommandType
	InputValueCommandType
)

type BaseCommand struct {
	Type  CommandType
	Value string
}

func (bc BaseCommand) GetType() CommandType {
	return bc.Type
}
func (bc BaseCommand) GetValue() string {
	return bc.Value
}

var ExitCommand = BaseCommand{ExitCommandType, "q"}
var DropCommand = BaseCommand{DropCommandType, "d"}
var ClaimCommand = BaseCommand{ClaimCommandType, "c"}
var FinishCommand = BaseCommand{FinishCommandType, "f"}

var scanner = bufio.NewScanner(os.Stdin)

func doReadInput() Command {
	scanner.Scan()
	str := scanner.Text()
	if str == "q" {
		return ExitCommand
	}
	if str == "d" {
		return DropCommand
	}
	if str == "c" {
		return ClaimCommand
	}
	if str == "f" {
		return FinishCommand
	}
	return BaseCommand{InputValueCommandType, str}
}

/*
var scanner = bufio.NewScanner(os.Stdin)

func ReadLine() string {
	scanner.Scan()
	return scanner.Text()
}

func main() {

	line := ""
	for line != "f" {
		line = ReadLine()
		fmt.Printf("Entered %s\n", line)
	}

}
*/
