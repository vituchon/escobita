package model

import (
	"fmt"
	"github.com/vituchon/escobita/util"
	"strings"
)

var ComputerPlayer Player = Player{
	Id:   0,
	Name: "Computer ",
}

func CalculateAction(match Match) PlayerAction {
	cards := GetMostImportantCards(match.Cards.Board, 8)
	playerCards := match.Cards.ByPlayer[ComputerPlayer]
	possibleTakeActions := CalculatePossibleTakeActions(cards, playerCards.Hand)
	analizedActions := AnalizeActions(possibleTakeActions, match)
	fmt.Println("analizedActions", analizedActions)
	if len(analizedActions.PossibleActions) > 0 {
		fmt.Println("Take action", analizedActions.RecomendedAction.PlayerTakeAction)
		return analizedActions.RecomendedAction.PlayerTakeAction
	} else {
		// TODO : There is a lot of room to improve this answer! func CalculatePossibleDropActions comming soon!
		dropAction := NewPlayerDropAction(ComputerPlayer, playerCards.Hand[0])
		fmt.Println("Drop action", dropAction)
		return dropAction
	}
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
							Player: ComputerPlayer,
						},
						BoardCards: boardSubcombination,
						HandCard:   handCard,
					}
					takeActions = append(takeActions, takeAction)
				}
			}
		}
	}

	//fmt.Println("before",len(takeActions))
	after := removeDuplicates(takeActions)
	//fmt.Println("after",len(after))
	return  after
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
			//fmt.Println("BoardCards A:",takeAction.BoardCards,"\nBoardCards B:",anotherTakeAction.BoardCards, "\n", "hasSameBoardCards" ,hasSameBoardCards, "hasSameHandCard",hasSameHandCard)
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

func (analysis TakeActionsAnalysisResult) String() string {
	var sb strings.Builder
	for i, possibleAction := range analysis.PossibleActions {
		sb.WriteString(fmt.Sprintf("%d. Sugerencia:%+v\n", i+1, possibleAction))
	}
	sb.WriteString(fmt.Sprintf("Recomendada:%+v\n", analysis.RecomendedAction))
	return sb.String()
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
	goldenSevenCard := util.Find(cards, func(card Card) bool { return card.IsGoldenSeven() })
	sevenCards := util.Filter(cards, func(card Card) bool { return card.IsSevenRank() })
	goldenCards := util.Filter(cards, func(card Card) bool { return card.IsGoldenSuit() })

	var cardsSet map[Card]bool = make(map[Card]bool)
	if goldenSevenCard != nil {
		addToSetUptoLimitSize(cardsSet, []Card{*goldenSevenCard}, uptoCount)
	}
	//fmt.Println("1 cardset:", cardsSet)
	rest := uptoCount - len(cardsSet)
	maxSevenCards := min(2, len(sevenCards))
	atMostTwoSevenCards := sevenCards[0:min(maxSevenCards, rest)] // considering at most 2 cards (if golded seven is included the only one seven will be added as the golded seven DO count as a seven too!)
	addToSetUptoLimitSize(cardsSet, atMostTwoSevenCards, uptoCount)
    //fmt.Println("2 cardset:", cardsSet)
	addToSetUptoLimitSize(cardsSet, goldenCards, uptoCount)
	//fmt.Println("3 cardset:", cardsSet)
	addToSetUptoLimitSize(cardsSet, cards, uptoCount)
	//fmt.Println("4 cardset:", cardsSet)

	mostImportantCards := make([]Card,0,  len(cardsSet))
	for card := range cardsSet {
		//fmt.Println("card",card)
		mostImportantCards = append(mostImportantCards, card)
	}
	return mostImportantCards
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func addToSetUptoLimitSize[T comparable](set map[T]bool, values []T, limitCount int) {
	if len(set) >= limitCount {
		return
	}

	for _, value := range values {
		set[value] = true
		if len(set) >= limitCount {
			return
		}
	}
}
