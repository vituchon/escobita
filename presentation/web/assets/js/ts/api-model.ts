
namespace Api {
  export interface Game {
    id: number;
    match: any;
    players: Player[];
  }

  export interface Player {
    id: number;
    name: string;
  }
}
