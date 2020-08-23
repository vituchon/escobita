
namespace Api {
  export interface Game {
    id: number;
    match: any;
    players: any;
  }

  export interface Player {
    id: number;
    name: string;
  }
}
