/// <reference path='./third_party_definitions/_definitions.ts' />

namespace Api {

  export interface Card {
    id: number;
    suit: number;
    rank: number;
  }

  export interface Player {
    id?: number;
    name: string;
  }

  export interface Match {
    players?: Player[];
    actionsByPlayerName: ActionsByPlayerName;
    playerActions: PlayerAction[];
    matchCards: MatchCards; // TODO: rename member to cards to avoid stuttering
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
    playerId?: number; // owner
    matchs?: Match[];
    players?: Player[];
    scoreByPlayerName?: ScoreByPlayerName;
    currentMatch?: Match;

    // these below properties exists on the back...
    /*actionsByPlayerName?: ActionsByPlayerName;
    playerActions?: PlayerAction[];*/
  }

  export interface ScoreByPlayerName extends _.Dictionary<number> {
    [name:string]: number;
  }

  export interface ActionsByPlayerName extends _.Dictionary<PlayerAction> {
    [name:string]: PlayerAction;
  }

  interface BasePlayerAction {
    player: Player;
  }

  export interface MatchCards {
    board?: Card[];
    left: Card[];
    byPlayerName: MatchCardsByPlayerName;
  }

  export interface MatchCardsByPlayerName extends _.Dictionary<PlayerMatchCards> {
    [name:string]: PlayerMatchCards;
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

  export interface ScoreSummaryByPlayerName extends _.Dictionary<PlayerScoreSummary> {
    [name:string]: PlayerScoreSummary;
  }

  export interface Message {
    id?: number;
    playerId: number;
    gameId: number;
    text: string;
    created?: number;
  }


}
