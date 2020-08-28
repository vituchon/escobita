
namespace Api {
  export interface Game {
    id?: number;
    match: any;
    players: Player[];
  }

  export interface Player {
    id?: number;
    name: string;
  }

  export interface Message {
    id?: number;
    playerId: number;
    gameId: number;
    text: string;
    created?: number;
  }
}
