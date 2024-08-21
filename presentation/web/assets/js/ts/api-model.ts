/// <reference path='./third_party_definitions/_definitions.ts' />

namespace Api {

  export interface Card {
    id: number;
    suit: number;
    rank: number;
  }

  export interface Player {
    id: number;
    name: string;
  }

  export interface MapByPlayer<T> {
    set(player: Player, data: T): void
    get(player: Player): T
  }

  export interface Match {
    players?: Player[];
    actionsByPlayer: MapByPlayer<Api.PlayerAction>;
    playerActions: PlayerAction[];
    matchCards: MatchCards;
    firstPlayerIndex: number;
    roundNumber: number;
    currentRound?: Round;
  }

  export interface Round {
    currentTurnPlayer: Player;
    consumedTurns: number;
    number: number;
  }

  export interface Game {
    id?: number;
    name: string;
    owner?: Player;
    matchs?: Match[]; // previous played matchs
    players?: Player[];
    currentMatch?: Match;
  }

  interface BasePlayerAction {
    player: Player;
  }

  export interface MatchCards {
    board?: Card[];
    left: Card[];
    byPlayer: MapByPlayer<PlayerMatchCards>;
  }

  export interface MatchCardsByPlayerUniqueKey extends _.Dictionary<PlayerMatchCards> {
    [uniqueKey:string]: PlayerMatchCards;
  }

  export interface PlayerMatchCards {
    taken: Card[];
    hand: Card[];
  }

  export interface PlayerTakeAction extends BasePlayerAction {
    boardCards: Card[];
    handCard: Card;
    isEscobita?: boolean; // resolved on server (by courtesy, cuz the client has the tools to do for itself as well)
  }

  export interface PlayerDropAction extends BasePlayerAction  {
    player: Player;
    handCard: Card;
  }

  export type PlayerAction = PlayerTakeAction | PlayerDropAction

  export interface PlayerStatictics {
    cardsTakenCount: number;
    escobitasCount: number;
    seventiesScore: number;
    hasGoldSeven: boolean;
    goldCardsCount: number;
}

  export interface PlayerScoreSummary {
      score: number;
      statictics: PlayerStatictics;
  }

  export interface ScoreSummaryByPlayerUniqueKey extends _.Dictionary<PlayerScoreSummary> {
    [uniqueKey:string]: PlayerScoreSummary;
  }

  export interface Message {
    id?: number;
    playerId: number;
    gameId: number;
    text: string;
    created?: number;
  }


}
