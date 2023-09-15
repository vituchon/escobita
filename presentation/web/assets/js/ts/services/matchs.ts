
namespace Matchs {

  export namespace Rules {
    function sumValues(cards :Api.Card[]){
      const total = _.reduce(cards,(acc,card) => {
        return acc + determineValue(card)
      },0)
      return total
    }

    function determineValue(card: Api.Card) {
      if (card.rank < 8) {
        return card.rank
      } else {
        return card.rank - 2
      }
    }
    // Determines if a valid take action can be performed
    // TODO: rename to isValidTakeAction
    export function canTakeCards(handCard: Api.Card, boardCards: Api.Card[]) {
      return sumValues(boardCards.concat(handCard)) == 15
    }

    export interface PositionByPlayerName extends _.Dictionary<number> {
      [name:string]: number;
    }

    interface PlayerScore {
      name: string,
      score: number,
    }
    export function calculatePositionByPlayerName(stats: Api.ScoreSummaryByPlayerName): PositionByPlayerName {
      const asArray: PlayerScore[] = _.map(stats,(summary,playerName) => {
        return {
          name: playerName,
          score: summary.score,
        }
      })
      const sorted = _.sortBy(asArray,(elem) => -elem.score) // sort desc so the higher score goes in "first"
      const asMap = _.reduce(sorted,(acc, playerScore, index) => {
        acc[playerScore.name] = index
        return acc
      },<PositionByPlayerName>{})
      return asMap
    }

  }

  export function isTakeAction(action: Api.PlayerTakeAction | Api.PlayerDropAction): action is Api.PlayerTakeAction {
    return !_.isUndefined((action as Api.PlayerTakeAction ).boardCards)
  }

  export function createTakeAction(player: Api.Player,boardCards: Api.Card[], handCard: Api.Card): Api.PlayerTakeAction {
    return {
      player: player,
      boardCards: boardCards,
      handCard: handCard,
    }
  }

  export function createDropAction(player: Api.Player, handCard: Api.Card): Api.PlayerDropAction {
    return {
      player: player,
      handCard: handCard,
    }
  }

  export namespace Engine {

    export const BotPlayer: Api.Player = {
      id: 0,
      name: "Botty 🤖",
    }

    export function calculatePossibleTakeActions(boardCards: Api.Card[], handCards: Api.Card[], player: Api.Player = BotPlayer): Api.PlayerTakeAction[] {
      if (_.size(boardCards) == 0) {
        return []
      }
      const boardCombinations = Arrays.combine(boardCards)
      const takeActions: Api.PlayerTakeAction[]  = [];
      _.forEach(handCards,(handCard) => {
        _.forEach(boardCombinations,(boardCombination) => {
          const cardsCount = _.size(boardCombination)
          for (var i = 1; i <= cardsCount; i++) {
            const boardSubcombination = boardCombination.slice(0,i)
            const isFeasible = Rules.canTakeCards(handCard,boardSubcombination)
            if (isFeasible) {
              const takeAction = Matchs.createTakeAction(player,boardSubcombination, handCard)
              takeActions.push(takeAction)
            }
          }
        })
      })

      var withoutDups: Api.PlayerTakeAction[] = []
      takeActions.forEach(takeAction => {
        const isNotContained = !_.find(withoutDups,(tk) => {
          const hasSameBoardCards = Arrays.hasSameValues(takeAction.boardCards,tk.boardCards, (c1,c2) => c1.id === c2.id)
          const hasSameHandCard = takeAction.handCard.id === tk.handCard.id
          return hasSameBoardCards && hasSameHandCard
        })
        if (isNotContained) {
          withoutDups.push(takeAction)
        }
      });

      return withoutDups
    }
  }

}
