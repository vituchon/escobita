/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />

namespace Cards {

  export function isGoldenSeven(card: Api.Card) {
    return isSevenRank(card) && isGoldenSuit(card)
  }

  export function isGoldenSuit(card: Api.Card) {
    return card.suit == Cards.Suits.gold
  }

  export function isSevenRank(card: Api.Card) {
    return card.rank == 7
  }

  export namespace Suits {
    // dev notes: the values must match some at /local/escobita/model/card.go#Line:25
    export const sword = 0
    export const club = 1
    export const cup = 2
    export const gold = 3
    export const all = [sword, club, cup, gold]

    export const labels: _.Dictionary<string> = {
      [sword]: "espada",
      [club]: "basto",
      [cup]: "copa",
      [gold]: "oro",
    }

    export function translate(suit: number) {
      return labels[suit]
    }
  }

  export namespace Sprites {
    const width = 208;
    const height = 319;
    const sourceImgPath = 'presentation/web/assets/images/cards.png'

    export const offsetYBySuit: _.Dictionary<number> = {
      [Suits.gold]: 0,
      [Suits.cup]: 1 * height,
      [Suits.sword]: 2 * height,
      [Suits.club]: 3 * height,
    }

    export function generateImageTag (suit: number, rank: number) {
      const offsetX = width * rank;
      const offsetY = offsetYBySuit[suit];

      // following tips at https://www.htmlgoodies.com/beyond/css/working-with-css-image-sprites.html
      const style = `background-size: 40% 40%; zoom: 0.4; width: ${width}px; height: ${height}px; background: url(${sourceImgPath}) -${offsetX - width}px -${offsetY}px no-repeat;`
      return `<div title="${rank}_${Suits.translate(suit)}" style="${style}" />`
    }

  }
}

namespace Rounds {

  export interface Round extends Api.Round {
  }

  export function isPlayerTurn(round: Round, player: Players.Player): boolean {
    return !_.isEmpty(round.currentTurnPlayer) && round.currentTurnPlayer.name == player.name; //  TODO: now game players have id, USE ID to avoid a clash!
  }
}
