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

	return removeDuplicates(takeActions)
}

func removeDuplicates(takeActions []PlayerTakeAction) []PlayerTakeAction {
	// eliminating duplicates
	compareCards := func(left, right Card) int {
		return left.Id - right.Id
	}
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
	return withoutDups
}

func CalculateActionSymbolicScore(action PlayerAction, match Match) int {
	takeAction, isTakeAction := action.(PlayerTakeAction)
	if isTakeAction {
		return CalculateTakeActionSymbolicScore(takeAction, match)
	}
	return 0
}

func CalculateTakeActionSymbolicScore(action PlayerTakeAction, match Match) int {
	var employedCards []Card
	employedCards = append(employedCards, action.BoardCards...)
	employedCards = append(employedCards, action.HandCard)

	seventiesSymbolicScore := CountSevenRankCards(employedCards) * 3

	goldenSuitCardsSymbolicScore := CountGoldenSuitCards(employedCards) * 2
	goldSevenSymbolicScore := boolToInt[DetermineIsGoldenSevenIsUsed(employedCards)] * 10

	isEscobita := len(match.Cards.Board) == len(action.BoardCards)
	escobitaSymbolicScore := boolToInt[isEscobita] * 10

	return len(employedCards) + escobitaSymbolicScore + seventiesSymbolicScore + goldSevenSymbolicScore + goldenSuitCardsSymbolicScore
}

type TakeActionsAnalysisResult struct {
	PossibleActions  []SuggestedTakeAction
	RecomendedAction SuggestedTakeAction
}

func AnalizeActions(actions []PlayerTakeAction, match Match) TakeActionsAnalysisResult {
	var recommendedAction SuggestedTakeAction
	var suggestedTakeActions []SuggestedTakeAction
	var maxSymbolicScore = 0

	for _, action := range actions {
		symbolicScore := CalculateActionSymbolicScore(action, match)
		suggestedAction := SuggestedTakeAction{
			PlayerTakeAction: action,
			SymbolicScore:    symbolicScore,
		}
		suggestedTakeActions = append(suggestedTakeActions, suggestedAction)
		if symbolicScore > maxSymbolicScore {
			maxSymbolicScore = symbolicScore
			recommendedAction = suggestedAction
		}
	}
	return TakeActionsAnalysisResult{
		PossibleActions:  suggestedTakeActions,
		RecomendedAction: recommendedAction,
	}
}

func GetMostImportantCards(cards []Card, uptoCount int) []Card {
	//const set: Set<Card> = new Set() // dev notes: using set for taking leverage that it not possible to add duplicates. For example: after adding golden 7 it is impossible to re-add the golden 7 as golden card as it will be already added

	/*idx := slices.IndexFunc(cards, func(card Card) bool { return card.IsGoldenSeven() })
	sevenCards := util.Filter(cards, func(card Card) bool { return card.isSevenRank() })
	goldenCards := util.Filter(cards, func(card Card) bool { return card.isGoldenSuit() })

	if idx != -1 {

	}*/
	return nil
	/*if (Util.isDefined(goldenSevenCard)) {
	    addToSetUptoLimitSize(set, [goldenSevenCard], uptoCount)
	  }

	  var rest = uptoCount - set.size
	  const atMostTwoSevenCards = sevenCards.slice(0,Math.min(2,rest)) // considering at most 2 cards (if golded seven is included the only one seven will be added as the golded seven DO count as a seven too!)
	  addToSetUptoLimitSize(set, atMostTwoSevenCards, uptoCount)

	  addToSetUptoLimitSize(set, goldenCards, uptoCount)
	  addToSetUptoLimitSize(set, cards, uptoCount)
	  return Array.from(set)*/
}
