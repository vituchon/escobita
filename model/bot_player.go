package model

import (
	"github.com/vituchon/escobita/util"
)

var BotPlayer Player = Player{
	Id:   0,
	Name: "Botty ",
}

type SuggestedTakeAction struct {
	PlayerTakeAction
	SymbolicScore int
}

func CalculatePossibleTakeActions(boardCards []Card, handCards []Card) []PlayerTakeAction {
	if len(boardCards) == 0 {
		return nil
	}

	var takeActions []PlayerTakeAction = []PlayerTakeAction{}
	boardCombinations := util.GeneratePermutations(boardCards)
	for _, boardCombination := range boardCombinations {
		for _, handCard := range handCards {
			//fmt.Println("tengo esto", boardCombination)
			for i := 1; i < len(boardCombination); i++ {
				//fmt.Println("Recordanto de 0 a", i, "resulta", boardCombination[0:i])
				boardSubcombination := copyDeck(boardCombination[0:i])
				isFeasible := CanTakeCards(handCard, boardSubcombination)
				if isFeasible {
					takeAction := PlayerTakeAction{
						basePlayerAction: basePlayerAction{
							Player: BotPlayer,
						},
						BoardCards: boardSubcombination,
						HandCard:   handCard,
					}
					takeActions = append(takeActions, takeAction)
				}
			}
		}
	}

	compareCards := func(left, right Card) int {
		return left.Id - right.Id
	}
	// eliminating duplicates
	withoutDups := []PlayerTakeAction{}
	for _, takeAction := range takeActions {
		var isContained bool = false
		for _, anotherTakeAction := range withoutDups {
			hasSameBoardCards := util.HasSameValuesDisregardingOrder(takeAction.BoardCards, anotherTakeAction.BoardCards, compareCards)
			hasSameHandCard := takeAction.HandCard.Id == anotherTakeAction.HandCard.Id
			if hasSameBoardCards && hasSameHandCard {
				isContained = true
				break
			}
		}
		if !isContained {
			withoutDups = append(withoutDups, takeAction)
		}
	}

	// perform analyze action over each withoutDup for adding a symbolic score
	// TODO: implement calculate symbolic score func

	return withoutDups
}
