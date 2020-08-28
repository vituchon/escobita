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
    actionsByPlayer: ActionsByPlayer;
    playerAction: any[];
    matchCards: MatchCards;
    firstPlayerIndex: number;
    roundNumber: number;
  }

  export interface Game {
    id?: number;
    name: string;
    matchs?: Match[];
    players?: Player[];
    scorePerPlayer?: ScorePerPlayer;
    currentMatch?: Match;
  }

  export interface ScorePerPlayer extends _.Dictionary<number> { // TODO : rename to ScoreByPlayer
    [name:string]: number;
  }

  export interface ActionsByPlayer extends _.Dictionary<PlayerAction> {
    [name:string]: PlayerAction;
  }

  interface BasePlayerAction {
    player: Player;
  }

  export interface MatchCards {
    board?: Card[];
    left: Card[];
    perPlayer: MatchCardsByPlayer; // TODO : rename to byPlayer
  }

  export interface MatchCardsByPlayer extends _.Dictionary<PlayerMatchCards> {
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
