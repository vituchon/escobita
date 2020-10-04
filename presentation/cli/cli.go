// The presentation layer contains all resources concerned with creating an application interface
// Contains code designed to be used for command line interface in a interactive way (no offline mode provided!)
package cli

import (
	"bufio"
	"errors"
	"strconv"

	//"errors"
	"fmt"
	"local/escobita/model"
	"os"
)

func StartApp() {

	var players []model.Player = []model.Player{
		model.Player{Name: "Beto"},
		model.Player{Name: "Pepe"},
	}

	quitGame := false
	match := model.CreateAndBegins(players)
	for match.HasMoreRounds() && !quitGame {
		round := match.NextRound()
		for round.HasNextTurn() && !quitGame {
			player := round.NextTurn()
			playerCards := match.Cards.ByPlayer[player]
			fmt.Println("Restantes: ", match.Cards.Left)

			playerMustAct := true
			for playerMustAct {
				fmt.Println("Juega ", player)
				fmt.Println("En mesa: ", match.Cards.Board)
				fmt.Println("En mano: ", playerCards.Hand)
				fmt.Println("Tomadas: ", playerCards.Taken)

				fmt.Print("Que deseas hacer, jugar (t), soltar carta (d), salir (q): ")
				cmd := doReadInput()
				if cmd == TakeCommand {
					takeAction := readTakeActionFromStdin(player, match)
					fmt.Println(takeAction.HandCard)
					fmt.Println(takeAction.BoardCards)
					isValidClaim := model.CanTakeCards(takeAction.HandCard, takeAction.BoardCards)
					if isValidClaim {
						fmt.Println("La jugada es valida!")
						playerMustAct = false
						match.Take(takeAction)
					} else {
						fmt.Println("La jugada NO es valida")
					}
				} else if cmd == DropCommand {
					dropAction := readDropActionFromStdin(player, match)
					match.Drop(dropAction)
					playerMustAct = false
				} else if cmd == QuitCommand {
					playerMustAct = false
					quitGame = true
				}
				staticticsByPlayer := match.CalculateStaticticsByPlayer()
				//fmt.Printf("===ESTADISTICAS DEL PARTIDO===\n%+v\n", staticticsByPlayer)
				fmt.Printf("===Resutaldos===\n%+v\n", staticticsByPlayer.BuildScoreSummaryByPlayer())
			}

		}
		//fmt.Println(match)
	}
} api

func readTakeActionFromStdin(player model.Player, match model.Match) model.PlayerTakeAction {
	fmt.Println("==" + player.Name + " selecciona combinación entre las de mesa y una de mano==")
	playerCards := match.Cards.ByPlayer[player]
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
	return model.NewPlayerTakeAction(player, handCard, boardCards)
}

func readDropActionFromStdin(player model.Player, match model.Match) model.PlayerDropAction {
	fmt.Print("La de mano, id de carta: ")
	playerCards := match.Cards.ByPlayer[player]
	cardId := ReadSingleIntInput()
	handCard, err := playerCards.Hand.GetSingle(cardId)
	for err != nil {
		fmt.Println("Error: ", err)
		fmt.Print("La de mano, id de carta: ")
		cardId := ReadSingleIntInput()
		handCard, err = playerCards.Hand.GetSingle(cardId)
	}
	return model.NewPlayerDropAction(player, handCard)
}

func ReadMultipleIntsInput() (ints []int) {
	value, command := ReadInputAndParseAnInt()
	for command != nil && command != FinishCommand {
		if command == QuitCommand {
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
	if command == QuitCommand {
		os.Exit(0)
	}

	return value
}

func ReadInputAndParseAnInt() (int, Command) {
	value, err, command := doReadIntInput()
	isEnteringAnotherCommand := (command == QuitCommand || command == FinishCommand)
	for err != nil && !isEnteringAnotherCommand {
		value, err, command = doReadIntInput()
		isEnteringAnotherCommand = (command == QuitCommand || command == FinishCommand)
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
	QuitCommandType CommandType = iota
	FinishCommandType
	DropCommandType
	TakeCommandType
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

var QuitCommand = BaseCommand{QuitCommandType, "q"}
var DropCommand = BaseCommand{DropCommandType, "d"}
var TakeCommand = BaseCommand{TakeCommandType, "t"}
var FinishCommand = BaseCommand{FinishCommandType, "f"}

var scanner = bufio.NewScanner(os.Stdin)

func doReadInput() Command {
	scanner.Scan()
	str := scanner.Text()
	if str == QuitCommand.Value {
		return QuitCommand
	}
	if str == DropCommand.Value {
		return DropCommand
	}
	if str == TakeCommand.Value {
		return TakeCommand
	}
	if str == FinishCommand.Value {
		return FinishCommand
	}
	return BaseCommand{InputValueCommandType, str}
}
