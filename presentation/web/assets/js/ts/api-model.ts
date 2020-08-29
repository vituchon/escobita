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
    actionsByPlayer: ActionsByPlayerName;
    playerAction: any[];
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
    matchs?: Match[];
    players?: Player[];
    scoreByPlayer?: ScoreByPlayerName;
    currentMatch?: Match;
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
    byPlayer: MatchCardsByPlayerName;
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
    isEscobita: boolean;
  }

  export interface PlayerDropAction extends BasePlayerAction  {
    player: Player;
    handCard: Card;
  }

  export type PlayerAction = PlayerTakeAction | PlayerDropAction


  export interface Message {
    id?: number;
    playerId: number;
    gameId: number;
    text: string;
    created?: number;
  }


}
